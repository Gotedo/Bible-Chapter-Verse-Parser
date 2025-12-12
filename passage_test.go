package parser

import (
	"testing"
)

func TestPassage_Stringify(t *testing.T) {
	p := NewBiblePassageParser()

	// helper to create a reference using same parser (so Book pointers match)
	makeRef := func(bookName string, ch, v int, frag string) *BibleReference {
		// find book object by name
		var b *Book
		for _, bb := range p.books {
			if bb.Name == bookName {
				b = bb
				break
			}
		}
		if b == nil {
			t.Fatalf("book %s not found", bookName)
		}
		r, err := NewBibleReference(b, ch, v, frag)
		if err != nil {
			t.Fatalf("NewBibleReference error: %v", err)
		}
		return r
	}

	cases := []struct {
		name string
		from *BibleReference
		to   *BibleReference
		want string
	}{
		{"readme example a", makeRef("1 John", 5, 4, ""), makeRef("1 John", 5, 17, ""), "1 John 5:4-17"},
		{"readme example b", makeRef("1 John", 5, 19, ""), makeRef("1 John", 5, 21, ""), "1 John 5:19-21"},
		{"readme example c", makeRef("Esther", 2, 1, ""), makeRef("Esther", 2, 23, ""), "Esther 2"},
		{"fragment", makeRef("Philippians", 2, 14, ""), makeRef("Philippians", 2, 15, "a"), "Philippians 2:14-15a"},
		{"another fragment", makeRef("Mark", 1, 4, "b"), makeRef("Mark", 1, 15, ""), "Mark 1:4b-15"},
		{"entire book", makeRef("John", 1, 1, ""), makeRef("John", 21, 25, ""), "John"},
		{"whole chapter", makeRef("John", 3, 1, ""), makeRef("John", 3, 36, ""), "John 3"},
		{"single verse", makeRef("John", 3, 16, ""), makeRef("John", 3, 16, ""), "John 3:16"},
		{"multiple whole books", makeRef("Genesis", 1, 1, ""), makeRef("Exodus", 40, 38, ""), "Genesis 1:1 - Exodus 40:38"},
		{"passage spanning different chapters", makeRef("Genesis", 1, 1, ""), makeRef("Genesis", 4, 26, ""), "Genesis 1-4"},
		{"passage spanning different chapters with odd verses", makeRef("Genesis", 1, 5, ""), makeRef("Genesis", 4, 10, ""), "Genesis 1:5-4:10"},
		{"passage spanning different book", makeRef("Genesis", 1, 1, ""), makeRef("Exodus", 5, 2, ""), "Genesis 1:1 - Exodus 5:2"},
		{"singular Psalm", makeRef("Psalms", 1, 1, ""), makeRef("Psalms", 1, 6, ""), "Psalm 1"},
		{"verses in a single Psalm", makeRef("Psalms", 1, 2, ""), makeRef("Psalms", 1, 3, ""), "Psalm 1:2-3"},
		{"plural Psalms", makeRef("Psalms", 120, 1, ""), makeRef("Psalms", 134, 3, ""), "Psalms 120-134"},
		{"All of Psalms", makeRef("Psalms", 1, 1, ""), makeRef("Psalms", 150, 6, ""), "Psalms"},
		{"Psalm to Psalm", makeRef("Psalms", 117, 2, ""), makeRef("Psalms", 118, 1, ""), "Psalm 117:2-118:1"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := NewBiblePassage(c.from, c.to)
			got := p.String()
			if got != c.want {
				t.Fatalf("mismatch: got %q want %q", got, c.want)
			}
		})
	}
}
