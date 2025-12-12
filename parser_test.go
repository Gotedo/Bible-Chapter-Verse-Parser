package parser_test

import (
	"testing"

	parser "github.com/gotedo/bible-chapter-verse-parser"
)

func TestParseSimple(t *testing.T) {
	p := parser.NewBiblePassageParser()
	res, err := p.Parse("John 3:16")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 passage, got %d", len(res))
	}
	if res[0].From.String() == "" {
		t.Fatalf("unexpected empty from string")
	}
}
