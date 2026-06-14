package tagparser_test

import (
	"maps"
	"strconv"
	"strings"
	"testing"

	"quad4/tagparser/pkg/tagparser"
)

func TestParse_neverNil(t *testing.T) {
	t.Parallel()
	for _, s := range []string{"", "x", strings.Repeat("a", 10000)} {
		tag := tagparser.Parse(s)
		if tag == nil {
			t.Fatalf("nil for %q", trunc(s))
		}
	}
}

func TestParse_idempotent(t *testing.T) {
	t.Parallel()
	for i, tc := range tagTests {
		tc := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			a := tagparser.Parse(tc.tag)
			b := tagparser.Parse(tc.tag)
			if a.Name != b.Name {
				t.Fatalf("name mismatch: %q vs %q", a.Name, b.Name)
			}
			if !maps.Equal(a.Options, b.Options) {
				t.Fatalf("options mismatch: %#v vs %#v", a.Options, b.Options)
			}
		})
	}
}

func trunc(s string) string {
	if len(s) <= 40 {
		return s
	}
	return s[:37] + "..."
}
