package parser

import "fmt"

type BibleReference struct {
	Book     *Book
	Chapter  int
	Verse    int
	Fragment string
}

func NewBibleReference(book *Book, chapter, verse int, fragment string) (*BibleReference, error) {
	if verse > 0 {
		vmax, ok := book.ChapterStructure[chapter]
		if !ok {
			return nil, fmt.Errorf("chapter %d does not exist in %s", chapter, book.Name)
		}
		if verse > vmax {
			return nil, fmt.Errorf("verse %d does not exist in chapter %d of book %s", verse, chapter, book.Name)
		}
	}
	if fragment != "" && fragment != "a" && fragment != "b" && fragment != "c" {
		return nil, fmt.Errorf("invalid fragment")
	}
	return &BibleReference{Book: book, Chapter: chapter, Verse: verse, Fragment: fragment}, nil
}

func (r *BibleReference) IntegerNotation() int {
	return (1000000 * r.Book.Number) + (1000 * r.Chapter) + r.Verse
}

func (r *BibleReference) String() string {
	if r.Verse == 0 {
		return fmt.Sprintf("%s %d", r.Book.Name, r.Chapter)
	}
	return fmt.Sprintf("%s %d:%d%s", r.Book.Name, r.Chapter, r.Verse, r.Fragment)
}
