package tagparser_test

import (
	"maps"
	"testing"

	"github.com/Quad4-Software/pbt/pkg/pbt"
	"github.com/Quad4-Software/tagparser/pkg/tagparser"
)

func byteSliceToString(xs []int) string {
	b := make([]byte, len(xs))
	for i, v := range xs {
		b[i] = byte(v)
	}
	return string(b)
}

func TestPBTParseDeterministic(t *testing.T) {
	gen := pbt.Map("string",
		pbt.SliceOf(pbt.IntRange(0, 255), 0, 8192),
		byteSliceToString,
	)
	prop := pbt.ForAll(
		"two Parse calls yield identical Tag fields",
		gen,
		func(s string) bool {
			a := tagparser.Parse(s)
			b := tagparser.Parse(s)
			return a.Name == b.Name && maps.Equal(a.Options, b.Options)
		},
		pbt.WithShrinker[string](pbt.StringShrinker()),
	)
	pbt.Check(t, prop, pbt.WithRuns(300), pbt.WithSeed(45))
}

// absentKey returns a probe string guaranteed to be missing from m.
func absentKey(m map[string]string) string {
	probe := "q4\x00probe"
	for {
		if _, ok := m[probe]; !ok {
			return probe
		}
		probe += "x"
	}
}

func TestPBTHasOptionMatchesMap(t *testing.T) {
	gen := pbt.Map("string",
		pbt.SliceOf(pbt.IntRange(32, 126), 0, 2048),
		byteSliceToString,
	)
	prop := pbt.ForAll(
		"HasOption(k) matches Options membership in both directions",
		gen,
		func(s string) bool {
			tag := tagparser.Parse(s)
			for k := range tag.Options {
				if !tag.HasOption(k) {
					return false
				}
			}
			// A key absent from the map must report false.
			if tag.HasOption(absentKey(tag.Options)) {
				return false
			}
			return true
		},
		pbt.WithShrinker[string](pbt.StringShrinker()),
	)
	pbt.Check(t, prop, pbt.WithRuns(200), pbt.WithSeed(46))
}
