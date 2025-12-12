# Gotedo Bible Chapter Verse (BCV) Parser

## Table of contents

- Features
- Quick start
- API (parser package)
- Data source
- Tests and development
- Contributing notes
- LICENSE

## Features

- Parse many human-friendly Bible passage formats, including:
  - single verses, ranges, whole chapters, entire books
  - fragments (verse parts) like `15a`, `36B` (case-insensitive)
  - ranges spanning chapters and books
  - shorthand abbreviations and numeric book prefixes (e.g., `1 John`, `2 Cor`)
  - flexible separators: `,`, `;`, `&`, `and`
  - en-dash/em-dash and `to` for ranges
- Produces structured `BibleReference` objects with validation against canonical chapter/verse counts.

## Quick start

Prerequisites

- Go 1.20+ installed

Run tests

```bash
gofmt -w .
go test ./... -v
```

Use the parser in your code

```go
import (
		"fmt"
		"github.com/gotedo/bible-verse-parser/parser"
)

func main() {
		p := parser.NewBiblePassageParser()
		passages, err := p.Parse("John 3:16-18, Psalm 23")
		if err != nil {
				panic(err)
		}
		for _, pass := range passages {
				fmt.Println(pass.From.String(), "-", pass.To.String())
				// or fmt.Println(pass) to print shorthand formatting
		}
}
```

## API (parser package)

Public types and functions (overview)

- parser.NewBiblePassageParser() \*BiblePassageParser

  - Create a parser instance. It initialises books from `data.BibleStructure`.

- (*BiblePassageParser).Parse(versesString string) ([]*BiblePassage, error)

  - Parse an input string and return a slice of `*BiblePassage` or an error.

- type BiblePassage

  - Fields: `From *BibleReference`, `To *BibleReference`.
  - String() returns a PHP-like shorthand representation (e.g., `John 3:16-18`).

- type BibleReference

  - Fields: `Book *Book`, `Chapter int`, `Verse int`, `Fragment string` (optional: `a`, `b`, `c`).
  - Methods: `IntegerNotation() int` (sortable numeric notation), `String() string` (longhand name form).

- type Book
  - Fields: `Number int`, `Name string`, `SingularName string`, `Abbreviations []string`, `ChapterStructure map[int]int`.
  - Methods: `ChaptersInBook() int`, `VersesInChapter(ch int) (int, error)`.

Example usage (parsing and printing):

```go
p := parser.NewBiblePassageParser()
res, err := p.Parse("1 John 5:4-17, 19-21 & Esther 2")
if err != nil { /* handle */ }
for _, r := range res {
		// shorthand
		fmt.Println(r.String())
		// longhand
		fmt.Println(r.From.String(), "to", r.To.String())
}
```

Error cases

- Passing an empty string returns an error (mirrors PHPUnit's invalid tests).
- Invalid book names cause an error.
- Invalid ranges (end before start) cause an error.

## Tests and development

- Unit tests live in the `parser` package; tests were ported from the original PHPUnit suite in batches and cover many parsing edge cases.
- Run tests with `go test ./... -v`.
- The repository includes an `.editorconfig` with Go-friendly style settings. Run `gofmt -w .` before committing.

CI

- A GitHub Actions workflow at `.github/workflows/tests.yml` runs `go vet` and `go test` on pushes and PRs (matrix across Go versions).

## Contributing notes

- The PHP files remain in this repository for reference; if you remove them, ensure the generator still has the canonical data or preserve `data/static_bible_structure.go`.
- Follow the existing test-porting pattern: port tests in small batches (10–20 cases), run `go test` and iterate until green.

Developer checklist

- Run `gofmt -w .`.
- Run `go test ./... -v`.
- Run `go vet ./...`.
- Optionally run `golangci-lint run` (CI tries to run it if present).

## CREDIT

This is a port of the PHP version of https://github.com/TechWilk/bible-verse-parser to Golang.
