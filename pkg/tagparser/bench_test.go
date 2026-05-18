package tagparser_test

import (
	"strings"
	"testing"

	"github.com/Quad4-Software/tagparser/pkg/tagparser"
)

var benchSink *tagparser.Tag

func BenchmarkParse(b *testing.B) {
	cases := []struct {
		name string
		s    string
	}{
		{"empty", ""},
		{"simple", "hello"},
		{"kv", "hello:world,foo:bar,baz:qux"},
		{"quoted", "hello:'a,b',x:y"},
		{"parens", "hello:world('foo', 'bar')"},
		{"long", strings.Repeat("k:v,", 512) + "k:v"},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchSink = tagparser.Parse(tc.s)
			}
		})
	}
}

func BenchmarkParseParallel(b *testing.B) {
	s := "hello:world,foo:bar,baz:'x,y,z'"
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchSink = tagparser.Parse(s)
		}
	})
}

func BenchmarkParseParallel_long(b *testing.B) {
	s := strings.Repeat("a:1,", 256) + "z:9"
	b.ReportAllocs()
	b.SetBytes(int64(len(s)))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchSink = tagparser.Parse(s)
		}
	})
}

// BenchmarkParseThroughput measures ns/op and B/s for a medium-sized tag string.
func BenchmarkParseThroughput(b *testing.B) {
	s := strings.Repeat("key:value,", 128) + "last:end"
	b.SetBytes(int64(len(s)))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchSink = tagparser.Parse(s)
	}
}
