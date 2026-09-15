package tagparser_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/Quad4-Software/tagparser/pkg/tagparser"
)

// edgeCases pins exact expected output for boundary and adversarial inputs.
// Expected values are hand derived from the grammar, independent of the
// implementation under test.
var edgeCases = []struct {
	tag  string
	name string
	opts map[string]string
}{
	// Empty and separator only inputs.
	{"", "", nil},
	{",", "", nil},
	{",,", "", map[string]string{"": ""}},
	{",,,", "", map[string]string{"": ""}},
	{"a,", "a", nil},
	{",a", "", map[string]string{"a": ""}},
	{"a,,b", "a", map[string]string{"": "", "b": ""}},

	// Bare name vs empty value vs valued option.
	{"a", "a", nil},
	{"a:", "", map[string]string{"a": ""}},
	{"a:b", "", map[string]string{"a": "b"}},
	{"a:b:c", "", map[string]string{"a": "b:c"}},
	{":v", "v", nil},
	{":", "", nil},

	// Repeated keys: the last occurrence wins, and a bare segment writes an
	// option keyed by its value.
	{"k:1,k:2", "", map[string]string{"k": "2"}},
	{"k:v,k", "", map[string]string{"k": ""}},
	{"k,k:v", "k", map[string]string{"k": "v"}},

	// Unbalanced and empty quotes.
	{"k:'v", "", map[string]string{"k": "v"}},
	{"'abc", "abc", nil},
	{"''", "", nil},
	{"k:''", "", map[string]string{"k": ""}},

	// Unbalanced, nested, and over closed parens.
	{"k:(a,b", "", map[string]string{"k": "(a,b"}},
	{"k:a)b", "", map[string]string{"k": "a)b"}},
	{"k:(a(b)c),d:e", "", map[string]string{"k": "(a(b)c)", "d": "e"}},
	{"k:(a)b)c", "", map[string]string{"k": "(a)b)c"}},

	// Quoted commas and colons do not split.
	{"k:'a,b',c:d", "", map[string]string{"k": "a,b", "c": "d"}},
	{"k:'a:b'", "", map[string]string{"k": "a:b"}},
	{"k:a'b", "", map[string]string{"k": "a'b"}},
	{"k:'(a,b)',x:y", "", map[string]string{"k": "(a,b)", "x": "y"}},
	{"k:('a,b'),x:y", "", map[string]string{"k": "('a,b')", "x": "y"}},

	// Escapes: the next byte is literal outside quotes; a lone trailing
	// backslash is kept rather than producing a NUL byte.
	{`k:a\,b,c:d`, "", map[string]string{"k": "a,b", "c": "d"}},
	{`k:a\:b`, "", map[string]string{"k": "a:b"}},
	{`a\:b`, "", map[string]string{`a\`: "b"}},
	{`k:v\`, "", map[string]string{`k`: `v\`}},
	{`k:(v\`, "", map[string]string{"k": `(v\`}},
	{`k:(a\)b),x:y`, "", map[string]string{"k": "(a)b)", "x": "y"}},

	// Escaped quotes inside quoted values. A quote preceded by an odd run of
	// backslashes is literal; an even run leaves it as the terminator, so a
	// double backslash must not swallow the rest of the input.
	{`k:'a\'b'`, "", map[string]string{"k": "a'b"}},
	{`k:'a\\',x:y`, "", map[string]string{"k": `a\\`, "x": "y"}},
	{`k:'a\\\'b'`, "", map[string]string{"k": `a\\'b`}},

	// Quotes in the key part start a quoted value and discard what was read.
	{"a'b'c", "b", map[string]string{"c": ""}},
	{"'a','b'", "a", map[string]string{"b": ""}},
	{"k:'a'junk", "", map[string]string{"k": "a", "junk": ""}},
	{"'a'b", "a", map[string]string{"b": ""}},
	{"unun'", "", nil},
	{"'a' ,b", "a", map[string]string{"": "", "b": ""}},

	// Unicode and NUL bytes pass through unmodified.
	{"名前:値", "", map[string]string{"名前": "値"}},
	{"キー,'値,1'", "キー", map[string]string{"値,1": ""}},
	{"a:\x00b", "", map[string]string{"a": "\x00b"}},
	{"\x00:\x00", "", map[string]string{"\x00": "\x00"}},
	{"\x00", "\x00", nil},

	// Keys and values are space trimmed, including quoted values.
	{"  spaced  ,  k :  v  ", "spaced", map[string]string{"k": "v"}},
	{"a\t:\tb\n", "", map[string]string{"a": "b"}},
	{"' '", "", nil},
	{"k:' '", "", map[string]string{"k": ""}},
}

func TestParse_edgeCases(t *testing.T) {
	for i, tc := range edgeCases {
		tc := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tag := tagparser.Parse(tc.tag)
			if tag.Name != tc.name {
				t.Fatalf("name: got %q want %q (tag=%q)", tag.Name, tc.name, tc.tag)
			}
			if len(tag.Options) != len(tc.opts) {
				t.Fatalf("options: got %#v want %#v (tag=%q)", tag.Options, tc.opts, tc.tag)
			}
			for k, v := range tc.opts {
				got, ok := tag.Options[k]
				if !ok {
					t.Fatalf("option %q missing (tag=%q)", k, tc.tag)
				}
				if got != v {
					t.Fatalf("option %q: got %q want %q (tag=%q)", k, got, v, tc.tag)
				}
			}
		})
	}
}

// TestParse_edgeCasesOracle cross checks the same inputs against the naive
// oracle so a pinned typo cannot hide a shared misreading.
func TestParse_edgeCasesOracle(t *testing.T) {
	for _, tc := range edgeCases {
		assertMatchesOracle(t, tc.tag)
	}
}

func TestParse_longSingleToken(t *testing.T) {
	t.Parallel()
	long := strings.Repeat("x", 1<<20)
	tag := tagparser.Parse(long)
	if tag.Name != long {
		t.Fatalf("name length: got %d want %d", len(tag.Name), len(long))
	}
	if tag.Options != nil {
		t.Fatalf("unexpected options: %#v", tag.Options)
	}

	key := "k"
	tag = tagparser.Parse(key + ":" + long)
	if tag.Options[key] != long {
		t.Fatalf("value length: got %d want %d", len(tag.Options[key]), len(long))
	}
}
