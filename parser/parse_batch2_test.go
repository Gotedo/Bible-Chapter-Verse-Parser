package parser

import (
	"reflect"
	"testing"
)

func TestParse_Batch2(t *testing.T) {
	p := NewBiblePassageParser()

	cases := []struct {
		name string
		in   string
		want [][]string
	}{
		{"multiple ranges, one with end keyword", "Deut 6: 4-9, 16-end & Luke 15: 1-10", [][]string{{"Deuteronomy 6:4", "Deuteronomy 6:9"}, {"Deuteronomy 6:16", "Deuteronomy 6:25"}, {"Luke 15:1", "Luke 15:10"}}},
		{"three entire chapters from different books", "1 Peter 2, 5 & Job 34", [][]string{{"1 Peter 2:1", "1 Peter 2:25"}, {"1 Peter 5:1", "1 Peter 5:14"}, {"Job 34:1", "Job 34:37"}}},
		{"multiple ranges from book with prefix number", "1 Peter 2:15-16, 18-20", [][]string{{"1 Peter 2:15", "1 Peter 2:16"}, {"1 Peter 2:18", "1 Peter 2:20"}}},
		{"one entire psalm", "Psalm 34", [][]string{{"Psalms 34:1", "Psalms 34:22"}}},
		{"same abbreviation with dot", "2 Cor. 5: 11-21", [][]string{{"2 Corinthians 5:11", "2 Corinthians 5:21"}}},
		{"roman numeral book name uppercase", "I Samuel 10:22", [][]string{{"1 Samuel 10:22", "1 Samuel 10:22"}}},
		{"roman numeral book name lowercase", "i Samuel 10:22", [][]string{{"1 Samuel 10:22", "1 Samuel 10:22"}}},
		{"roman numerals on multiple book name", "II Kings 1:1 & I Samuel 10:22", [][]string{{"2 Kings 1:1", "2 Kings 1:1"}, {"1 Samuel 10:22", "1 Samuel 10:22"}}},
		{"horribly complex", "Genesis 1:1 - Exodus 5:2 & 6:3-4", [][]string{{"Genesis 1:1", "Exodus 5:2"}, {"Exodus 6:3", "Exodus 6:4"}}},
		{"complex example from issues/18", "Gen 1:1, 3-4; 4:26-5:1; Lev 4:5; 5:2; Phlm 1:2; 1 John 1;2 John 1; 3 John; Pss 1-2", [][]string{{"Genesis 1:1", "Genesis 1:1"}, {"Genesis 1:3", "Genesis 1:4"}, {"Genesis 4:26", "Genesis 5:1"}, {"Leviticus 4:5", "Leviticus 4:5"}, {"Leviticus 5:2", "Leviticus 5:2"}, {"Philemon 1:2", "Philemon 1:2"}, {"1 John 1:1", "1 John 1:10"}, {"2 John 1:1", "2 John 1:13"}, {"3 John 1:1", "3 John 1:15"}, {"Psalms 1:1", "Psalms 2:12"}}},
		{"end fragment", "Philippians 2:14-15a", [][]string{{"Philippians 2:14", "Philippians 2:15a"}}},
		{"start fragment", "John 4:7b-4:8", [][]string{{"John 4:7b", "John 4:8"}}},
		{"start fragment with more letters", "Mark 1v4b-15", [][]string{{"Mark 1:4b", "Mark 1:15"}}},
		{"single fragment", "Acts 2:39a", [][]string{{"Acts 2:39a", "Acts 2:39a"}}},
		{"fragment to fragment", "John 3:16b-17a", [][]string{{"John 3:16b", "John 3:17a"}}},
		{"numbered book without white space between book number and book name, without chapter and verse", "1Kings", [][]string{{"1 Kings 1:1", "1 Kings 22:53"}}},
		{"numbered book without white space between book number and book name, with chapter", "1Kings 1", [][]string{{"1 Kings 1:1", "1 Kings 1:53"}}},
		{"numbered book without white space between book number and book name, with chapter and verse", "1Kings 1:1", [][]string{{"1 Kings 1:1", "1 Kings 1:1"}}},
		{"numbered book without white space between book number and book name, with chapter and verses interval", "1Kings 1:1-2", [][]string{{"1 Kings 1:1", "1 Kings 1:2"}}},
		{"numbered book without white space between book number and book name, with chapter and selected verses", "1Kings 1:1,12", [][]string{{"1 Kings 1:1", "1 Kings 1:1"}, {"1 Kings 1:12", "1 Kings 1:12"}}},
		{"multiple numbered book without white space between book number and book name, without chapter and verse", "1Kings; 2Kings", [][]string{{"1 Kings 1:1", "1 Kings 22:53"}, {"2 Kings 1:1", "2 Kings 25:30"}}},
		{"multiple numbered book without white space between book number and book name, with chapter", "1Kings 1; 2Kings 1", [][]string{{"1 Kings 1:1", "1 Kings 1:53"}, {"2 Kings 1:1", "2 Kings 1:18"}}},
		{"multiple numbered book without white space between book number and book name, with chapter and verse", "1Kings 1:1; 2Kings 1:1", [][]string{{"1 Kings 1:1", "1 Kings 1:1"}, {"2 Kings 1:1", "2 Kings 1:1"}}},
		{"multiple numbered book without white space between book number and book name, with chapter and verses interval", "1Kings 1:1-2; 2Kings 1:1-2", [][]string{{"1 Kings 1:1", "1 Kings 1:2"}, {"2 Kings 1:1", "2 Kings 1:2"}}},
		{"multiple numbered book without white space between book number and book name, with chapter and selected verses", "1Kings 1:1,12; 2Kings 1:1,12", [][]string{{"1 Kings 1:1", "1 Kings 1:1"}, {"1 Kings 1:12", "1 Kings 1:12"}, {"2 Kings 1:1", "2 Kings 1:1"}, {"2 Kings 1:12", "2 Kings 1:12"}}},
		{"multiple numbered book without white space between book number and book name, in range", "1Kings 22:53-2Kings 1:12", [][]string{{"1 Kings 22:53", "2 Kings 1:12"}}},
		{"numbered book with white space between book number and book name, without chapter and verse", "1 Kings", [][]string{{"1 Kings 1:1", "1 Kings 22:53"}}},
		{"numbered book with white space between book number and book name, with chapter", "1 Kings 1", [][]string{{"1 Kings 1:1", "1 Kings 1:53"}}},
		{"numbered book with white space between book number and book name, with chapter and verse", "1 Kings 1:1", [][]string{{"1 Kings 1:1", "1 Kings 1:1"}}},
		{"numbered book with white space between book number and book name, with chapter and verses interval", "1 Kings 1:1-2", [][]string{{"1 Kings 1:1", "1 Kings 1:2"}}},
		{"numbered book with white space between book number and book name, with chapter and selected verses", "1 Kings 1:1,12", [][]string{{"1 Kings 1:1", "1 Kings 1:1"}, {"1 Kings 1:12", "1 Kings 1:12"}}},
		{"multiple numbered book with white space between book number and book name, without chapter and verse", "1 Kings; 2 Kings", [][]string{{"1 Kings 1:1", "1 Kings 22:53"}, {"2 Kings 1:1", "2 Kings 25:30"}}},
		{"multiple numbered book with white space between book number and book name, with chapter", "1 Kings 1; 2 Kings 1", [][]string{{"1 Kings 1:1", "1 Kings 1:53"}, {"2 Kings 1:1", "2 Kings 1:18"}}},
		{"multiple numbered book with white space between book number and book name, with chapter and verse", "1 Kings 1:1; 2 Kings 1:1", [][]string{{"1 Kings 1:1", "1 Kings 1:1"}, {"2 Kings 1:1", "2 Kings 1:1"}}},
		{"multiple numbered book with white space between book number and book name, with chapter and verses interval", "1 Kings 1:1-2; 2 Kings 1:1-2", [][]string{{"1 Kings 1:1", "1 Kings 1:2"}, {"2 Kings 1:1", "2 Kings 1:2"}}},
		{"multiple numbered book with white space between book number and book name, with chapter and selected verses", "1 Kings 1:1,12; 2 Kings 1:1,12", [][]string{{"1 Kings 1:1", "1 Kings 1:1"}, {"1 Kings 1:12", "1 Kings 1:12"}, {"2 Kings 1:1", "2 Kings 1:1"}, {"2 Kings 1:12", "2 Kings 1:12"}}},
		{"multiple numbered book with white space between book number and book name, in range", "1 Kings 22:53-2 Kings 1:12", [][]string{{"1 Kings 22:53", "2 Kings 1:12"}}},
		{"book to book with no whitespace", "1 Kings-2 Kings", [][]string{{"1 Kings 1:1", "2 Kings 25:30"}}},
		{"en dash between verses", "1 Corinthians 13:1–3", [][]string{{"1 Corinthians 13:1", "1 Corinthians 13:3"}}},
		{"en dash between chapters", "1 Corinthians 12–13", [][]string{{"1 Corinthians 12:1", "1 Corinthians 13:13"}}},
		{"en dash between books", "1 Corinthians–2 Corinthians", [][]string{{"1 Corinthians 1:1", "2 Corinthians 13:14"}}},
		{"em dash between verses", "1 Corinthians 13:1—3", [][]string{{"1 Corinthians 13:1", "1 Corinthians 13:3"}}},
		{"em dash between chapters", "1 Corinthians 12—13", [][]string{{"1 Corinthians 12:1", "1 Corinthians 13:13"}}},
		{"em dash between books", "1 Corinthians—2 Corinthians", [][]string{{"1 Corinthians 1:1", "2 Corinthians 13:14"}}},
		{"double space and capital fragment issues/92", "Luke 24  36B-48", [][]string{{"Luke 24:36b", "Luke 24:48"}}},
		{"capital fragment issues/92", "Luke 24:36B-48", [][]string{{"Luke 24:36b", "Luke 24:48"}}},
		{"and separator issues/93", "John 1 and 2", [][]string{{"John 1:1", "John 1:51"}, {"John 2:1", "John 2:25"}}},
		{"ch standing for both chronicles and chapter", "2 CH ch 13  v 01", [][]string{{"2 Chronicles 13:1", "2 Chronicles 13:1"}}},
		{"ch standing for chronicles not chapter", "2 CH 13  v 01", [][]string{{"2 Chronicles 13:1", "2 Chronicles 13:1"}}},
		{"single chapter books can omit chapter", "Obadiah v 01", [][]string{{"Obadiah 1:1", "Obadiah 1:1"}}},
		{"single chapter books can omit chapter variant", "Obadiah verse 01", [][]string{{"Obadiah 1:1", "Obadiah 1:1"}}},
		{"single chapter books can omit chapter variant larger verse", "Obadiah verse 15", [][]string{{"Obadiah 1:15", "Obadiah 1:15"}}},
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
