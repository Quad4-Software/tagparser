package tagparser_test

import (
	"maps"
	"strings"
	"testing"

	"github.com/Quad4-Software/pbt/pkg/pbt"
	"github.com/Quad4-Software/tagparser/pkg/tagparser"
)

// naiveTag is the oracle result type; it mirrors tagparser.Tag.
type naiveTag struct {
	name string
	opts map[string]string
}

// naiveParse is an independent implementation of the tag grammar used as a
// test oracle. It walks the input by index instead of the byte slice scanner
// used by the parser under test. Semantics:
//
//   - segments are separated by commas that are not inside quotes or parens;
//     one space after a separating comma is skipped
//   - a colon ends the key part of a segment; a single quote anywhere in the
//     key part starts a quoted value and discards the bytes read so far
//   - outside quotes, backslash escapes the next byte literally and an opening
//     paren starts a nested paren group; a lone trailing backslash is kept
//   - inside quotes the value runs to a quote preceded by an even number of
//     backslashes; a quote preceded by an odd number becomes a literal quote
//   - keys and values are space trimmed; the first segment becomes the name
//     when it has no key, otherwise it is an option like the rest
func naiveParse(s string) naiveTag {
	var t naiveTag
	hasName := false
	emit := func(key, value string) {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if !hasName {
			hasName = true
			if key == "" {
				t.name = value
				return
			}
		}
		if t.opts == nil {
			t.opts = make(map[string]string)
		}
		if key == "" {
			t.opts[value] = ""
		} else {
			t.opts[key] = value
		}
	}
	for i := 0; i < len(s); {
		i = naiveSegment(s, i, emit)
	}
	return t
}

// naiveSegment scans one segment starting at i, emits it, and returns the
// index where the next segment starts.
func naiveSegment(s string, i int, emit func(key, value string)) int {
	var b []byte
	for i < len(s) {
		switch s[i] {
		case ',':
			i++
			if i < len(s) && s[i] == ' ' {
				i++
			}
			emit("", string(b))
			return i
		case ':':
			key := string(b)
			return naiveValue(s, i+1, key, emit)
		case '\'':
			v, j := naiveQuoted(s, i+1)
			emit("", v)
			i = j
			if i < len(s) && s[i] == ',' {
				i++
				if i < len(s) && s[i] == ' ' {
					i++
				}
			}
			return i
		default:
			b = append(b, s[i])
			i++
		}
	}
	if len(b) > 0 {
		emit("", string(b))
	}
	return i
}

// naiveValue scans the value part of a keyed segment starting at i.
func naiveValue(s string, i int, key string, emit func(key, value string)) int {
	if i < len(s) && s[i] == '\'' {
		v, j := naiveQuoted(s, i+1)
		emit(key, v)
		i = j
		if i < len(s) && s[i] == ',' {
			i++
			if i < len(s) && s[i] == ' ' {
				i++
			}
		}
		return i
	}

	var b []byte
	for i < len(s) {
		switch s[i] {
		case '\\':
			if i+1 < len(s) {
				b = append(b, s[i+1])
				i += 2
			} else {
				b = append(b, s[i])
				i++
			}
		case '(':
			b = append(b, s[i])
			i++
			i, b = naiveBrackets(s, i, b)
		case ',':
			i++
			if i < len(s) && s[i] == ' ' {
				i++
			}
			emit(key, string(b))
			return i
		default:
			b = append(b, s[i])
			i++
		}
	}
	emit(key, string(b))
	return i
}

// naiveBrackets consumes a paren group after its opening paren, tracking
// nesting and backslash escapes. Returns the index after the group and the
// accumulated literal bytes.
func naiveBrackets(s string, i int, b []byte) (int, []byte) {
	lvl := 0
	for i < len(s) {
		switch s[i] {
		case '\\':
			if i+1 < len(s) {
				b = append(b, s[i+1])
				i += 2
			} else {
				b = append(b, s[i])
				i++
			}
		case '(':
			b = append(b, s[i])
			lvl++
			i++
		case ')':
			b = append(b, s[i])
			lvl--
			i++
			if lvl < 0 {
				return i, b
			}
		default:
			b = append(b, s[i])
			i++
		}
	}
	return i, b
}

// naiveQuoted consumes a quoted value starting just after the opening quote
// at i. Returns the unescaped value and the index after the closing quote, or
// the end of input when the quote is unterminated.
func naiveQuoted(s string, i int) (string, int) {
	var b []byte
	for i < len(s) {
		j := strings.IndexByte(s[i:], '\'')
		if j < 0 {
			b = append(b, s[i:]...)
			return string(b), len(s)
		}
		seg := s[i : i+j]
		i += j + 1
		bs := 0
		for k := len(seg) - 1; k >= 0 && seg[k] == '\\'; k-- {
			bs++
		}
		if bs%2 == 1 {
			b = append(b, seg[:len(seg)-1]...)
			b = append(b, '\'')
			continue
		}
		b = append(b, seg...)
		return string(b), i
	}
	return string(b), i
}

// assertMatchesOracle fails the test when Parse and naiveParse disagree on s.
func assertMatchesOracle(t *testing.T, s string) {
	t.Helper()
	tag := tagparser.Parse(s)
	want := naiveParse(s)
	if tag.Name != want.name || !maps.Equal(tag.Options, want.opts) {
		t.Fatalf("oracle mismatch on %q: got name=%q opts=%#v, want name=%q opts=%#v",
			truncInput(s), tag.Name, tag.Options, want.name, want.opts)
	}
}

// TestOracle_table pins oracle agreement on inputs that stress every branch of
// the grammar, including malformed ones.
func TestOracle_table(t *testing.T) {
	t.Parallel()
	inputs := []string{
		"", ",", ",,", "a,", ",a", "a,,b",
		"a", "a:", "a:b", "a:b:c", ":v", ":",
		"k:1,k:2", "k:v,k", "k,k:v",
		"k:'v", "'abc", "''", "k:''",
		"k:(a,b", "k:a)b", "k:(a(b)c),d:e", "k:(a)b)c",
		"k:'a,b',c:d", "k:'a:b'", "k:a'b",
		"k:'(a,b)',x:y", "k:('a,b'),x:y",
		`k:a\,b,c:d`, `k:a\:b`, `a\:b`,
		`k:v\`, `k:(v\`, `k:'a\'b'`, `k:'a\\',x:y`, `k:'a\\\'b'`,
		"a'b'c", "'a','b'", "k:'a'junk", "'a'b", "unun'",
		"名前:値", "キー,'値,1'", "a:\x00b", "\x00:\x00", "\x00",
		"  spaced  ,  k :  v  ", "a\t:\tb\n",
		"a(b:c", "'a' ,b", "k:(a\\)b),x:y",
		strings.Repeat("x", 1<<16),
		strings.Repeat(",", 64),
	}
	for _, s := range inputs {
		assertMatchesOracle(t, s)
	}
}

// oracleAlphabet is biased toward grammar metacharacters so short generated
// strings still exercise quotes, parens, escapes, and separators.
var oracleAlphabet = []byte("aA0 :,()'\\\t\x00")

func oracleAlphabetGen() pbt.Generator[string] {
	return pbt.Map("alphabet",
		pbt.SliceOf(pbt.IntRange(0, len(oracleAlphabet)-1), 0, 96),
		func(xs []int) string {
			b := make([]byte, len(xs))
			for i, v := range xs {
				b[i] = oracleAlphabet[v]
			}
			return string(b)
		})
}

// oraclePieces are segment building blocks; joining them with commas produces
// mostly well formed tags with occasional malformed combinations.
var oraclePieces = []string{
	"name", "k", "v", "a:b", "x:'p,q'", "f(g,h)", `esc\,ape`, "'q,n'",
	"k:", ":v", "", " ", "世界", "k:(a(b)c)", `k:'a\'b'`, `tr\`, "k:v\\",
}

func oracleStructuredGen() pbt.Generator[string] {
	return pbt.Map("structured",
		pbt.SliceOf(pbt.IntRange(0, len(oraclePieces)-1), 0, 12),
		func(xs []int) string {
			parts := make([]string, len(xs))
			for i, v := range xs {
				parts[i] = oraclePieces[v]
			}
			return strings.Join(parts, ",")
		})
}

func TestPBTMatchesOracle_alphabet(t *testing.T) {
	prop := pbt.ForAll(
		"Parse output equals naive oracle on metachar alphabet",
		oracleAlphabetGen(),
		func(s string) bool {
			tag := tagparser.Parse(s)
			want := naiveParse(s)
			return tag.Name == want.name && maps.Equal(tag.Options, want.opts)
		},
		pbt.WithShrinker[string](pbt.StringShrinker()),
	)
	pbt.Check(t, prop, pbt.WithRuns(2000), pbt.WithSeed(101))
}

func TestPBTMatchesOracle_structured(t *testing.T) {
	prop := pbt.ForAll(
		"Parse output equals naive oracle on structured tags",
		oracleStructuredGen(),
		func(s string) bool {
			tag := tagparser.Parse(s)
			want := naiveParse(s)
			return tag.Name == want.name && maps.Equal(tag.Options, want.opts)
		},
		pbt.WithShrinker[string](pbt.StringShrinker()),
	)
	pbt.Check(t, prop, pbt.WithRuns(2000), pbt.WithSeed(102))
}

func TestPBTMatchesOracle_arbitraryBytes(t *testing.T) {
	prop := pbt.ForAll(
		"Parse output equals naive oracle on arbitrary bytes",
		pbt.Map("string",
			pbt.SliceOf(pbt.IntRange(0, 255), 0, 512),
			byteSliceToString,
		),
		func(s string) bool {
			tag := tagparser.Parse(s)
			want := naiveParse(s)
			return tag.Name == want.name && maps.Equal(tag.Options, want.opts)
		},
		pbt.WithShrinker[string](pbt.StringShrinker()),
	)
	pbt.Check(t, prop, pbt.WithRuns(2000), pbt.WithSeed(103))
}
