package parser

import "testing"

func TestParse_InvalidVerses(t *testing.T) {
	p := NewBiblePassageParser()

	cases := []string{"", "Psalm 34-20"}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			_, err := p.Parse(c)
			if err == nil {
				t.Fatalf("expected error for invalid input %q, got nil", c)
			}
		})
	}
}

func TestParse_InvalidVerseBooks(t *testing.T) {
	p := NewBiblePassageParser()

	cases := []string{"Bob", "1"}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			_, err := p.Parse(c)
			if err == nil {
				t.Fatalf("expected error for invalid book %q, got nil", c)
			}
		})
	}
}
