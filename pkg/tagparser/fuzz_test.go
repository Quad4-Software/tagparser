package tagparser_test

import (
	"maps"
	"strings"
	"testing"

	"github.com/Quad4-Software/tagparser/pkg/tagparser"
)

func FuzzParse(f *testing.F) {
	seeds := []string{
		"",
		"hello",
		"hello:world",
		"hello:'a,b'",
		`x:'D\'Angelo'`,
		"a,b,c:d",
		strings.Repeat(",", 500),
	}
	for _, s := range seeds {
		f.Add(s)
	}
	for _, tc := range tagTests {
		f.Add(tc.tag)
	}
	f.Fuzz(func(t *testing.T, s string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic on %q: %v", truncInput(s), r)
			}
		}()
		tag := tagparser.Parse(s)
		if tag == nil {
			t.Fatal("nil Tag")
		}
		_ = tag.Name
		_ = tag.Options
	})
}

func FuzzDeterminism(f *testing.F) {
	for _, tc := range tagTests {
		f.Add(tc.tag)
	}
	f.Add("k:v,x:y")
	f.Fuzz(func(t *testing.T, s string) {
		a := tagparser.Parse(s)
		b := tagparser.Parse(s)
		if a.Name != b.Name {
			t.Fatalf("name: %q vs %q", a.Name, b.Name)
		}
		if !maps.Equal(a.Options, b.Options) {
			t.Fatalf("options: %#v vs %#v", a.Options, b.Options)
		}
	})
}

func truncInput(s string) string {
	if len(s) <= 64 {
		return s
	}
	return s[:61] + "..."
}
