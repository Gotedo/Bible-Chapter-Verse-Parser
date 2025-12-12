package parser

import "fmt"

type BiblePassage struct {
	From *BibleReference
	To   *BibleReference
}

func NewBiblePassage(from, to *BibleReference) *BiblePassage {
	return &BiblePassage{From: from, To: to}
}

func (p *BiblePassage) String() string {
	from := p.From
	to := p.To
	// Mirror the PHP formatting rules precisely.
	// Check for entire book: from 1:1 to last chapter:last verse in same book
	if to.Book == from.Book && from.Chapter == 1 && from.Verse == 1 {
		if to.Book.ChaptersInBook() == to.Chapter {
			if vmax, _ := to.Book.VersesInChapter(to.Chapter); vmax == to.Verse {
				return from.Book.Name
			}
		}
	}

	trailer := fmt.Sprintf(" %d", from.Chapter)

	// Format "John 3" or "Psalm 3"
	if to.Book == from.Book && to.Chapter == from.Chapter && (from.Verse == 0 || (from.Verse == 1 && func() bool {
		vmax, _ := to.Book.VersesInChapter(to.Chapter)
		return vmax == to.Verse
	}())) {
		return from.Book.SingularName + trailer
	}

	trailer = trailer + fmt.Sprintf(":%d%s", from.Verse, from.Fragment)

	// Format "John 3:16"
	if to.Book == from.Book && to.Chapter == from.Chapter && to.Verse == from.Verse {
		return from.Book.SingularName + trailer
	}

	// Format "John 3:16-17"
	if from.Chapter == to.Chapter {
		return from.Book.SingularName + trailer + "-" + fmt.Sprintf("%d%s", to.Verse, to.Fragment)
	}

	toString := fmt.Sprintf("%d:%d%s", to.Chapter, to.Verse, to.Fragment)

	// Format cross-book: "John 3:16 - Acts 1:1" (note spaces around dash)
	if from.Book != to.Book {
		return from.Book.Name + trailer + " - " + to.Book.Name + " " + toString
	}

	// Psalms plural case: "Psalms 120-134"
	if from.Verse == 1 {
		if vmax, _ := to.Book.VersesInChapter(to.Chapter); vmax == to.Verse {
			return from.Book.Name + " " + fmt.Sprintf("%d-%d", from.Chapter, to.Chapter)
		}
	}

	return from.Book.SingularName + trailer + "-" + toString
}
