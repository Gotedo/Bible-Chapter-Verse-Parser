package parser

import "fmt"

type Book struct {
	Number           int
	Name             string
	SingularName     string
	Abbreviations    []string
	ChapterStructure map[int]int
}

func NewBook(number int, name, singular string, abbr []string, chapterStructure map[int]int) *Book {
	return &Book{Number: number, Name: name, SingularName: singular, Abbreviations: abbr, ChapterStructure: chapterStructure}
}

func (b *Book) NumberFn() int          { return b.Number }
func (b *Book) NameFn() string         { return b.Name }
func (b *Book) SingularNameFn() string { return b.SingularName }

func (b *Book) ChaptersInBook() int {
	return len(b.ChapterStructure)
}

func (b *Book) VersesInChapter(chapter int) (int, error) {
	v, ok := b.ChapterStructure[chapter]
	if !ok {
		return 0, fmt.Errorf("chapter %d does not exist in %s", chapter, b.Name)
	}
	return v, nil
}
