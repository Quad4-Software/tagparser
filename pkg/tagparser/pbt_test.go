package tagparser_test

import (
	"maps"
	"testing"

	"git.quad4.io/Go-Libs/tagparser/v2/pkg/tagparser"
	"git.quad4.io/Go-Libs/pbt/pkg/pbt"
)

func byteSliceToString(xs []int) string {
	b := make([]byte, len(xs))
	for i, v := range xs {
		b[i] = byte(v)
	}
	return string(b)
}

func TestPBTParseNoPanic(t *testing.T) {
	gen := pbt.Map("string",
		pbt.SliceOf(pbt.IntRange(0, 255), 0, 4096),
		byteSliceToString,
	)
	prop := pbt.ForAll(
		"Parse accepts arbitrary strings without panic",
		gen,
		func(s string) bool {
			_ = tagparser.Parse(s)
			return true
		},
		pbt.WithShrinker[string](pbt.StringShrinker()),
	)
	pbt.Check(t, prop, pbt.WithRuns(200), pbt.WithSeed(44))
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

func TestPBTHasOptionMatchesMap(t *testing.T) {
	gen := pbt.Map("string",
		pbt.SliceOf(pbt.IntRange(32, 126), 0, 2048),
		byteSliceToString,
	)
	prop := pbt.ForAll(
		"HasOption(k) matches Options[k]",
		gen,
		func(s string) bool {
			tag := tagparser.Parse(s)
			if tag.Options == nil {
				return true
			}
			for k := range tag.Options {
				if !tag.HasOption(k) {
					return false
				}
			}
			return true
		},
		pbt.WithShrinker[string](pbt.StringShrinker()),
	)
	pbt.Check(t, prop, pbt.WithRuns(200), pbt.WithSeed(46))
}
