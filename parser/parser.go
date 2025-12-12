package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gotedo/bible-verse-parser/data"
)

var defaultSeparators = []string{"&", ",", ";", "and"}

type BiblePassageParser struct {
	separators []string
	books      map[int]*Book
	bookAbbr   map[string]int
}

func NewBiblePassageParser() *BiblePassageParser {
	p := &BiblePassageParser{separators: defaultSeparators, books: map[int]*Book{}, bookAbbr: map[string]int{}}
	for num, bd := range data.BibleStructure {
		b := NewBook(num, bd.Name, bd.SingularName, bd.Abbreviations, bd.ChapterStructure)
		p.books[num] = b
		p.bookAbbr[StandardiseString(b.Name)] = num
		for _, a := range b.Abbreviations {
			p.bookAbbr[StandardiseString(a)] = num
		}
	}
	return p
}

func (p *BiblePassageParser) Parse(versesString string) ([]*BiblePassage, error) {
	if strings.TrimSpace(versesString) == "" {
		return nil, fmt.Errorf("unable to parse reference")
	}
	substitutions := []struct{ re, rep string }{
		// insert spaces between letters and digits and vice versa to normalize inputs like 'chapter3verse16'
		{`(?i)([A-Za-z])([0-9])`, `$1 $2`},
		// avoid splitting numeric+fragment (e.g. 15a). only split when the following letter is not a/b/c
		{`(?i)([0-9])([d-z])`, `$1 $2`},
		{`(?i)(—|–)`, `-`},
		{`(?i)[^a-z]to[^a-z]`, `-`},
		{`(?i)([^a-z])chapter([^a-z])`, `$1ch$2`},
		{`(?i)([^a-z])c([^a-z])`, `$1ch$2`},
		{`(?i)([^a-z])verses?([^a-z])`, `$1 v $2`},
		{`(?i)(^|;\ *|\-)([\d])([a-zA-Z])`, `$1 $2 $3`},
	}
	for _, s := range substitutions {
		r := regexp.MustCompile(s.re)
		versesString = r.ReplaceAllString(versesString, s.rep)
	}

	sections := SplitOnSeparators(p.separators, versesString)

	passages := []*BiblePassage{}
	lastBook := ""
	var lastChapter *int
	var lastVerse *int

	for _, section := range sections {
		section = strings.TrimSpace(section)
		if section == "" {
			continue
		}

		splitSection := strings.Split(section, "-")
		if len(splitSection) > 2 {
			return nil, fmt.Errorf("Range is too complex")
		}

		fromReference, startVerse, lb, lc, lv, lastFragment, err := p.parseStartReference(splitSection[0], lastBook, lastChapter, lastVerse)
		if err != nil {
			return nil, err
		}

		lastBook = lb
		lastChapter = lc
		lastVerse = lv

		var toReference *BibleReference

		if len(splitSection) == 1 {
			endBookObject, err := p.getBookFromAbbreviation(lastBook)
			if err != nil {
				return nil, err
			}
			// if the parsed start contained an explicit verse (startVerse != nil), then
			// the range is that single verse. Otherwise default to whole chapter/end as before.
			if startVerse != nil {
				toReference = fromReference
			} else {
				endChapterForReference := 0
				if lastChapter != nil {
					endChapterForReference = *lastChapter
				} else {
					endChapterForReference = endBookObject.ChaptersInBook()
				}

				endVerse := 0
				if lastVerse != nil {
					endVerse = *lastVerse
				} else {
					vv, _ := endBookObject.VersesInChapter(endChapterForReference)
					endVerse = vv
				}

				tr, err := NewBibleReference(endBookObject, endChapterForReference, endVerse, lastFragment)
				if err != nil {
					return nil, err
				}
				toReference = tr
			}
		} else {
			matches, err := p.parseReference(splitSection[1])
			if err != nil {
				return nil, err
			}

			endBook := ""
			var endChapter *int
			var endVerse *int
			endFragment := ""

			if matches["book"] != "" {
				endBook = matches["book"]
			}

			endBookObject, err := p.getBookFromAbbreviation(func() string {
				if endBook != "" {
					return endBook
				}
				return lastBook
			}())
			if err != nil {
				return nil, err
			}

			// (explicit verse marker handled in parseStartReference; nothing to do here)

			if matches["chapter_or_verse"] != "" {
				if startVerse != nil && matches["verse"] == "" {
					// this is an end verse
					if matches["chapter_or_verse"] == "end" {
						v, _ := endBookObject.VersesInChapter(*lastChapter)
						ev := v
						endVerse = &ev
					} else {
						s := matches["chapter_or_verse"]
						if len(s) > 0 {
							last := s[len(s)-1]
							if last == 'a' || last == 'b' || last == 'c' || last == 'A' || last == 'B' || last == 'C' {
								endFragment = strings.ToLower(string(last))
							}
						}
						// strip trailing fragment letters before atoi
						numStr := strings.TrimRightFunc(matches["chapter_or_verse"], func(r rune) bool {
							return r == 'a' || r == 'b' || r == 'c' || r == 'A' || r == 'B' || r == 'C'
						})
						vi, _ := strconv.Atoi(numStr)
						endVerse = &vi
					}
				} else {
					if matches["chapter_or_verse"] == "end" {
						ec := endBookObject.ChaptersInBook()
						endChapter = &ec
					} else {
						ciStr := matches["chapter_or_verse"]
						// strip trailing fragment if present
						ciNumStr := strings.TrimRightFunc(ciStr, func(r rune) bool {
							return r == 'a' || r == 'b' || r == 'c' || r == 'A' || r == 'B' || r == 'C'
						})
						ci, _ := strconv.Atoi(ciNumStr)
						endChapter = &ci
					}
				}
			}

			if matches["verse"] != "" {
				// parse verse with optional fragment
				vs := matches["verse"]
				if len(vs) > 0 {
					last := vs[len(vs)-1]
					if last == 'a' || last == 'b' || last == 'c' || last == 'A' || last == 'B' || last == 'C' {
						endFragment = strings.ToLower(string(last))
					}
				}
				numStr := strings.TrimRightFunc(matches["verse"], func(r rune) bool { return r == 'a' || r == 'b' || r == 'c' || r == 'A' || r == 'B' || r == 'C' })
				vi, _ := strconv.Atoi(numStr)
				endVerse = &vi
			}

			endChapterForReference := 0
			if endChapter != nil {
				endChapterForReference = *endChapter
			} else if lastChapter != nil {
				endChapterForReference = *lastChapter
			} else {
				endChapterForReference = endBookObject.ChaptersInBook()
			}

			// set last values
			lastBook = endBookObject.NameFn()
			if endChapter != nil {
				lastChapter = endChapter
			}
			lastVerse = endVerse

			tr, err := NewBibleReference(endBookObject, endChapterForReference, func() int {
				if endVerse != nil {
					return *endVerse
				}
				v, _ := endBookObject.VersesInChapter(endChapterForReference)
				return v
			}(), endFragment)
			if err != nil {
				return nil, err
			}
			toReference = tr
		}

		if fromReference.IntegerNotation() > toReference.IntegerNotation() {
			return nil, fmt.Errorf("references end is before beginning")
		}

		passages = append(passages, NewBiblePassage(fromReference, toReference))
	}

	return passages, nil
}

func (p *BiblePassageParser) parseStartReference(textReference, lastBook string, lastChapter, lastVerse *int) (*BibleReference, *int, string, *int, *int, string, error) {
	matches, err := p.parseReference(textReference)
	if err != nil {
		return nil, nil, "", nil, nil, "", err
	}

	book := ""
	var chapter *int
	var verse *int
	fragment := ""

	if matches["book"] != "" {
		book = matches["book"]
		lastChapter = nil
		lastVerse = nil
	} else {
		book = lastBook
	}

	// get start book object early so we can make decisions for single-chapter books
	startBookObject, err := p.getBookFromAbbreviation(book)
	if err != nil {
		return nil, nil, "", nil, nil, "", err
	}

	// detect whether parseReference marked an explicit verse indicator (e.g. 'v' or 'verse')
	explicitVerse := false
	if strings.HasSuffix(matches["verse"], "__explicitverse") {
		explicitVerse = true
		matches["verse"] = strings.TrimSuffix(matches["verse"], "__explicitverse")
	}

	if matches["chapter_or_verse"] != "" {
		// If the book has only one chapter, prefer to treat the numeric token as a
		// verse only when the original text explicitly indicated a verse (e.g. used
		// 'v' or the word 'verse'). Otherwise treat it as a chapter (to match PHP behaviour).
		if startBookObject.ChaptersInBook() == 1 && explicitVerse {
			if matches["chapter_or_verse"] == "end" {
				ci := -1
				chapter = &ci
			} else {
				ci, frag := parseNumFragment(matches["chapter_or_verse"])
				// if numeric token looks like a verse, use it
				vMax, _ := startBookObject.VersesInChapter(1)
				if ci > 0 && ci <= vMax {
					verse = &ci
					if frag != "" {
						fragment = frag
					}
				} else {
					chapter = &ci
					if frag != "" {
						fragment = frag
					}
				}
			}
			lastVerse = nil
		} else if lastVerse == nil || matches["verse"] != "" {
			// chapter
			if matches["chapter_or_verse"] == "end" {
				ci := -1
				chapter = &ci
			} else {
				ci, frag := parseNumFragment(matches["chapter_or_verse"])
				chapter = &ci
				if frag != "" {
					fragment = frag
				}
			}
			lastVerse = nil
		} else {
			// verse
			vi, frag := parseNumFragment(matches["chapter_or_verse"])
			verse = &vi
			if frag != "" {
				fragment = frag
			}
		}
	}

	if matches["verse"] != "" {
		// If chapter is nil and the book only has one chapter, treat this as a verse, not a chapter
		if chapter == nil && startBookObject.ChaptersInBook() == 1 {
			vi, frag := parseNumFragment(matches["verse"])
			verse = &vi
			if frag != "" {
				fragment = frag
			}
		} else if chapter == nil {
			ci, _ := strconv.Atoi(matches["verse"])
			chapter = &ci
		} else {
			vi, frag := parseNumFragment(matches["verse"])
			verse = &vi
			if frag != "" {
				fragment = frag
			}
		}
	}

	if chapter == nil {
		chapter = lastChapter
	}

	lastBookName := startBookObject.NameFn()

	ch := 1
	if chapter != nil && *chapter > 0 {
		ch = *chapter
	}
	v := 1
	if verse != nil && *verse > 0 {
		v = *verse
	}

	fromRef, err := NewBibleReference(startBookObject, ch, v, fragment)
	if err != nil {
		return nil, nil, "", nil, nil, "", err
	}

	return fromRef, verse, lastBookName, chapter, verse, fragment, nil
}

func (p *BiblePassageParser) parseReference(reference string) (map[string]string, error) {
	reference = strings.ToLower(reference)
	regex := regexp.MustCompile(`^\s*(?P<book>(?:[0-9]+\s+)?[^0-9]+)?(?:(?P<chapter_or_verse>[0-9]+[abc]?)?(?:\s*[\. \:v]+\s*(?P<verse>[0-9]+[abc]?(?:end)?))?)?\s*$`)
	result := regex.FindStringSubmatch(reference)
	if result == nil {
		return nil, fmt.Errorf("unable to parse reference")
	}
	matches := map[string]string{"book": "", "chapter_or_verse": "", "verse": ""}
	for i, name := range regex.SubexpNames() {
		if i != 0 && name != "" {
			matches[name] = strings.TrimSpace(result[i])
		}
	}

	if (matches["book"] == "start" || matches["book"] == "end") && matches["chapter_or_verse"] == "" && matches["verse"] == "" {
		matches["chapter_or_verse"] = matches["book"]
		matches["book"] = ""
	}

	// Remove trailing " ch" or " v" from book if present
	if len(matches["book"]) >= 3 {
		if strings.HasSuffix(matches["book"], " ch") {
			trimmed := strings.TrimSpace(matches["book"][0 : len(matches["book"])-3])
			if regexp.MustCompile(`[A-Za-z]+`).MatchString(trimmed) {
				matches["book"] = trimmed
			}
		}
	}

	if len(matches["book"]) >= 2 {
		if strings.HasSuffix(matches["book"], " v") {
			trimmed := strings.TrimSpace(matches["book"][0 : len(matches["book"])-2])
			if regexp.MustCompile(`[A-Za-z]+`).MatchString(trimmed) {
				matches["book"] = trimmed
				if matches["verse"] == "" {
					matches["verse"] = matches["chapter_or_verse"]
					matches["chapter_or_verse"] = "1"
					// mark that the original text explicitly used 'v' or 'verse'
					matches["verse"] = matches["verse"] + "__explicitverse"
				}
			}
		}
		// also accept the full word 'verse' or 'verses' attached to the book
		if strings.HasSuffix(matches["book"], " verse") || strings.HasSuffix(matches["book"], " verses") {
			// remove the trailing word
			parts := strings.Fields(matches["book"]) // split on whitespace
			if len(parts) > 0 {
				trimmed := strings.Join(parts[0:len(parts)-1], " ")
				if regexp.MustCompile(`[A-Za-z]+`).MatchString(trimmed) {
					matches["book"] = trimmed
					if matches["verse"] == "" {
						matches["verse"] = matches["chapter_or_verse"]
						matches["chapter_or_verse"] = "1"
						// mark explicit use of the word 'verse'
						matches["verse"] = matches["verse"] + "__explicitverse"
					}
				}
			}
		}
	}

	return matches, nil
}

func (p *BiblePassageParser) getBookFromAbbreviation(bookAbbreviation string) (*Book, error) {
	bn, err := p.getBookNumber(bookAbbreviation)
	if err != nil {
		return nil, err
	}
	b, ok := p.books[bn]
	if !ok {
		return nil, fmt.Errorf("invalid book number \"%d\"", bn)
	}
	return b, nil
}

func (p *BiblePassageParser) getBookNumber(bookAbbreviation string) (int, error) {
	s := StandardiseString(bookAbbreviation)
	if v, ok := p.bookAbbr[s]; ok {
		return v, nil
	}
	return 0, fmt.Errorf("invalid book name \"%s\"", bookAbbreviation)
}

// parseNumFragment parses a string like "16b" or "36B" and returns the integer and the
// optional fragment (a/b/c) in lowercase. If no fragment is present, fragment is empty.
func parseNumFragment(s string) (int, string) {
	if s == "" {
		return 0, ""
	}
	last := s[len(s)-1]
	frag := ""
	if last == 'a' || last == 'b' || last == 'c' || last == 'A' || last == 'B' || last == 'C' {
		frag = strings.ToLower(string(last))
		s = s[:len(s)-1]
	}
	if s == "" {
		return 0, frag
	}
	n, _ := strconv.Atoi(s)
	return n, frag
}
