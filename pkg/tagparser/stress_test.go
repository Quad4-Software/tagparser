package tagparser_test

import (
	"maps"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Quad4-Software/tagparser/pkg/tagparser"
)

func TestStressConcurrentParse(t *testing.T) {
	if testing.Short() {
		t.Skip("stress")
	}
	const workers = 64
	const iters = 2000
	inputs := []string{
		"",
		"a:b,c:d",
		strings.Repeat("x:y,", 128) + "z:9",
		`msg:'a,b',k:v`,
	}
	want := make([]naiveTag, len(inputs))
	for i, s := range inputs {
		want[i] = naiveParse(s)
	}
	var wg sync.WaitGroup
	var mismatches atomic.Int32
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				for j, s := range inputs {
					tag := tagparser.Parse(s)
					if tag == nil || tag.Name != want[j].name || !maps.Equal(tag.Options, want[j].opts) {
						mismatches.Add(1)
					}
				}
			}
		}()
	}
	wg.Wait()
	if mismatches.Load() != 0 {
		t.Fatalf("concurrent Parse mismatches: %d", mismatches.Load())
	}
	runtime.GC()
}

func TestStressLongInput(t *testing.T) {
	if testing.Short() {
		t.Skip("stress")
	}
	const n = 256 * 1024
	var b strings.Builder
	b.Grow(n)
	for i := 0; i < n; i++ {
		switch i % 4 {
		case 0:
			b.WriteByte('a')
		case 1:
			b.WriteByte(',')
		case 2:
			b.WriteByte(':')
		default:
			b.WriteByte('\'')
		}
	}
	s := b.String()
	tag := tagparser.Parse(s)
	if tag == nil {
		t.Fatal("nil")
	}
	want := naiveParse(s)
	if tag.Name != want.name || !maps.Equal(tag.Options, want.opts) {
		t.Fatalf("oracle mismatch: got name=%q opts=%#v, want name=%q opts=%#v",
			tag.Name, tag.Options, want.name, want.opts)
	}
}

func TestStressLongRepeatedSegments(t *testing.T) {
	if testing.Short() {
		t.Skip("stress")
	}
	s := strings.Repeat("key:value,", 4096) + "last:1"
	tag := tagparser.Parse(s)
	if tag == nil {
		t.Fatal("nil")
	}
	if !tag.HasOption("last") {
		t.Fatal("expected last option")
	}
	if tag.Options["key"] != "value" || tag.Options["last"] != "1" {
		t.Fatalf("options: %#v", tag.Options)
	}
	if len(tag.Options) != 2 {
		t.Fatalf("expected 2 options, got %#v", tag.Options)
	}
}

// TestStressManySegments keeps the segment scan iterative: a recursive parser
// overflows the goroutine stack well below this segment count.
func TestStressManySegments(t *testing.T) {
	if testing.Short() {
		t.Skip("stress")
	}
	const n = 1 << 24
	s := strings.Repeat(",", n)
	tag := tagparser.Parse(s)
	if tag == nil {
		t.Fatal("nil")
	}
	want := naiveParse(s)
	if tag.Name != want.name || !maps.Equal(tag.Options, want.opts) {
		t.Fatalf("oracle mismatch: got name=%q opts=%#v, want name=%q opts=%#v",
			tag.Name, tag.Options, want.name, want.opts)
	}
}

func TestStressDeepParens(t *testing.T) {
	if testing.Short() {
		t.Skip("stress")
	}
	const depth = 2000
	var b strings.Builder
	b.WriteString("k:")
	for i := 0; i < depth; i++ {
		b.WriteByte('(')
	}
	b.WriteString("x")
	for i := 0; i < depth; i++ {
		b.WriteByte(')')
	}
	tag := tagparser.Parse(b.String())
	if tag == nil {
		t.Fatal("nil")
	}
	want := strings.Repeat("(", depth) + "x" + strings.Repeat(")", depth)
	if tag.Options["k"] != want {
		t.Fatalf("value mismatch: got len=%d want len=%d", len(tag.Options["k"]), len(want))
	}
}
