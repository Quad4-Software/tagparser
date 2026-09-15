package tagparser_test

import (
	"maps"
	"strconv"
	"testing"

	"github.com/Quad4-Software/tagparser/pkg/tagparser"
)

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
