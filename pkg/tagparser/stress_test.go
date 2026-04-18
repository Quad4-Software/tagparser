package tagparser_test

import (
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"git.quad4.io/Go-Libs/tagparser/v2/pkg/tagparser"
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
	var wg sync.WaitGroup
	var nilTags atomic.Int32
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				for _, s := range inputs {
					tag := tagparser.Parse(s)
					if tag == nil {
						nilTags.Add(1)
					}
				}
			}
		}()
	}
	wg.Wait()
	if nilTags.Load() != 0 {
		t.Fatalf("nil tag results: %d", nilTags.Load())
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
	_ = tag.Name
	_ = tag.Options
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
}

func TestStressDeepParens(t *testing.T) {
	if testing.Short() {
		t.Skip("stress")
	}
	var b strings.Builder
	b.WriteString("k:")
	for i := 0; i < 2000; i++ {
		b.WriteByte('(')
	}
	b.WriteString("x")
	for i := 0; i < 2000; i++ {
		b.WriteByte(')')
	}
	tag := tagparser.Parse(b.String())
	if tag == nil {
		t.Fatal("nil")
	}
}
