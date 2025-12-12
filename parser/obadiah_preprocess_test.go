package parser

import (
	"regexp"
	"testing"
)

func TestObadiahPreprocess(t *testing.T) {
	s := "Obadiah verse 01"
	substitutions := []struct{ re, rep string }{
		{`(?i)([A-Za-z])([0-9])`, `$1 $2`},
		{`(?i)([0-9])([d-z])`, `$1 $2`},
		{`(?i)(—|–)`, `-`},
		{`(?i)[^a-z]to[^a-z]`, `-`},
		{`(?i)([^a-z])chapter([^a-z])`, `$1ch$2`},
		{`(?i)([^a-z])c([^a-z])`, `$1ch$2`},
		{`(?i)([^a-z])verses?([^a-z])`, `$1v$2`},
		{`(?i)(^|;\ *|\-)([\d])([a-zA-Z])`, `$1 $2 $3`},
	}
	t.Logf("Orig: %q", s)
	for _, sub := range substitutions {
		r := regexp.MustCompile(sub.re)
		s = r.ReplaceAllString(s, sub.rep)
		t.Logf(" after %s -> %q", sub.re, s)
	}
}
