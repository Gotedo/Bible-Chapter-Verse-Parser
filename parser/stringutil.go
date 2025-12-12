package parser

import (
	"regexp"
	"strings"
)

func StandardiseString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = regexp.MustCompile("[^a-z0-9 ]").ReplaceAllString(s, "")
	return s
}

func SplitOnSeparators(separators []string, text string) []string {
	normalised := strings.ReplaceAll(text, separators[0], separators[0])
	for _, sep := range separators[1:] {
		normalised = strings.ReplaceAll(normalised, sep, separators[0])
	}
	return strings.Split(normalised, separators[0])
}
