package parser

import (
	"reflect"
	"testing"
)

func TestParse_Batch1(t *testing.T) {
	p := NewBiblePassageParser()

	cases := []struct {
		name string
		in   string
		want [][]string
	}{
		{"readme example", "1 John 5:4-17, 19-21 & Esther 2", [][]string{{"1 John 5:4", "1 John 5:17"}, {"1 John 5:19", "1 John 5:21"}, {"Esther 2:1", "Esther 2:23"}}},
		{"colon as verse delimiter", "John 3:16", [][]string{{"John 3:16", "John 3:16"}}},
		{"v character as verse delimiter", "John 3v16", [][]string{{"John 3:16", "John 3:16"}}},
		{"colon and v character as verse delimiter", "John 3: v16", [][]string{{"John 3:16", "John 3:16"}}},
		{"space and v character as verse delimiter", "John 3 v16", [][]string{{"John 3:16", "John 3:16"}}},
		{"vv characters as verse delimiter", "John 3vv16", [][]string{{"John 3:16", "John 3:16"}}},
		{"space and vv character as verse delimiter", "John 3 vv16", [][]string{{"John 3:16", "John 3:16"}}},
		{"c and v characters as verse delimiter", "John c3v16", [][]string{{"John 3:16", "John 3:16"}}},
		{"ch and v characters as verse delimiter", "John ch3v16", [][]string{{"John 3:16", "John 3:16"}}},
		{"ch and vv characters as verse delimiter", "John ch3vv16", [][]string{{"John 3:16", "John 3:16"}}},
		{"space, ch and v characters as verse delimiter", "John ch3 v16", [][]string{{"John 3:16", "John 3:16"}}},
		{"chapter and verse characters as verse delimiter", "John chapter3verse16", [][]string{{"John 3:16", "John 3:16"}}},
		{"space, chapter and verse characters as verse delimiter", "John chapter3 verse16", [][]string{{"John 3:16", "John 3:16"}}},
		{"spaces, chapter and verse characters as verse delimiter", "John chapter 3 verse 16", [][]string{{"John 3:16", "John 3:16"}}},
		{"period as verse delimiter", "John 3.16", [][]string{{"John 3:16", "John 3:16"}}},
		{"space as verse delimiter", "John 3 16", [][]string{{"John 3:16", "John 3:16"}}},
		{"entire book", "John", [][]string{{"John 1:1", "John 21:25"}}},
		{"whole chapter", "John 3", [][]string{{"John 3:1", "John 3:36"}}},
		{"two whole chapters", "John 3, 4", [][]string{{"John 3:1", "John 3:36"}, {"John 4:1", "John 4:54"}}},
		{"two verse ranges in same chapters", "John 3:16-18, 19-22", [][]string{{"John 3:16", "John 3:18"}, {"John 3:19", "John 3:22"}}},
		{"verses spanning different chapters", "Gen 1:1-4:26", [][]string{{"Genesis 1:1", "Genesis 4:26"}}},
		{"verses spanning different chapters with numeric book", "1 John 3:1-4:12", [][]string{{"1 John 3:1", "1 John 4:12"}}},
		{"verses spanning different chapters shorthand", "Gen 1-4:26", [][]string{{"Genesis 1:1", "Genesis 4:26"}}},
		{"verse range without hyphen", "1 John 3:1 to 4:12", [][]string{{"1 John 3:1", "1 John 4:12"}}},
		{"verse range without hyphen longhand", " 1 John 3:12 to 1 John 4:21", [][]string{{"1 John 3:12", "1 John 4:21"}}},
		{"two single verses in different chapters", "Gen 1:1; 4:26", [][]string{{"Genesis 1:1", "Genesis 1:1"}, {"Genesis 4:26", "Genesis 4:26"}}},
		{"single verse and whole chapter in different books", "John 3:16 & Isiah 22", [][]string{{"John 3:16", "John 3:16"}, {"Isaiah 22:1", "Isaiah 22:25"}}},
		{"verse range with end keyword", "John 3:16-end", [][]string{{"John 3:16", "John 3:36"}}},
		{"chapter range with end keyword", "John 3-end", [][]string{{"John 3:1", "John 21:25"}}},
		{"abbreviated books with single and verse ranges", "Is 53: 1-6 & 2 Cor 5: 20-21", [][]string{{"Isaiah 53:1", "Isaiah 53:6"}, {"2 Corinthians 5:20", "2 Corinthians 5:21"}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := p.Parse(c.in)
			if err != nil {
				t.Fatalf("parse error for %q: %v", c.in, err)
			}
			gotPairs := make([][]string, len(got))
			for i, pass := range got {
				gotPairs[i] = []string{pass.From.String(), pass.To.String()}
			}
			if !reflect.DeepEqual(gotPairs, c.want) {
				t.Fatalf("mismatch for %q\n got: %#v\nwant: %#v", c.in, gotPairs, c.want)
			}
		})
	}
}
