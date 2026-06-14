package tagparser_test

import (
	"strings"
	"testing"

	"quad4/tagparser/pkg/tagparser"
)

// FuzzUntrustedTag feeds potentially hostile inputs (control chars, long runs, unicode).
// Name avoids the "FuzzParse" prefix so -fuzz=FuzzParse does not match two tests.
func FuzzUntrustedTag(f *testing.F) {
	seeds := []string{
		"\x00",
		"\x00\x00\x00",
		"\xff\xfe",
		"\u202eRTL",
		"\ufeff",
		strings.Repeat("a,", 2000),
		strings.Repeat(":", 500),
		strings.Repeat("'", 300),
		strings.Repeat("\\", 400),
		strings.Repeat("(", 200) + strings.Repeat(")", 200),
		"key:" + strings.Repeat("x", 8000),
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic: %v (input len=%d)", r, len(s))
			}
		}()
		tag := tagparser.Parse(s)
		if tag == nil {
			t.Fatal("nil Tag")
		}
		_ = tag.Name
		if tag.Options != nil {
			for k, v := range tag.Options {
				_ = k
				_ = v
			}
		}
	})
}
