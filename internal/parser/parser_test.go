package parser_test

import (
	"bytes"
	"testing"

	"quad4/tagparser/internal/parser"
)

func TestParser_Read_exhausts(t *testing.T) {
	t.Parallel()
	p := parser.New([]byte("ab"))
	if got := p.Read(); got != 'a' {
		t.Fatalf("first Read: %q", got)
	}
	if got := p.Read(); got != 'b' {
		t.Fatalf("second Read: %q", got)
	}
	if p.Valid() {
		t.Fatal("expected exhausted")
	}
	if got := p.Read(); got != 0 {
		t.Fatalf("Read past end: %q", got)
	}
}

func TestParser_Peek(t *testing.T) {
	t.Parallel()
	p := parser.New([]byte("x"))
	if p.Peek() != 'x' {
		t.Fatalf("Peek: %q", p.Peek())
	}
	if p.Read() != 'x' {
		t.Fatal("Read")
	}
	if p.Peek() != 0 {
		t.Fatalf("Peek past end: %q", p.Peek())
	}
}

func TestParser_Skip(t *testing.T) {
	t.Parallel()
	p := parser.New([]byte("  x"))
	if !p.Skip(' ') {
		t.Fatal("first space")
	}
	if !p.Skip(' ') {
		t.Fatal("second space")
	}
	if p.Read() != 'x' {
		t.Fatal("x")
	}
}

func TestParser_SkipBytes(t *testing.T) {
	t.Parallel()
	p := parser.New([]byte("foobar"))
	if !p.SkipBytes([]byte("foo")) {
		t.Fatal("foo")
	}
	if b := p.Bytes(); string(b) != "bar" {
		t.Fatalf("Bytes: %q", b)
	}
	if p.SkipBytes([]byte("nomatch")) {
		t.Fatal("expected false")
	}
}

func TestParser_ReadSep(t *testing.T) {
	t.Parallel()
	p := parser.New([]byte("a,b:c"))
	part, ok := p.ReadSep(',')
	if !ok || string(part) != "a" {
		t.Fatalf("first ReadSep: %q %v", part, ok)
	}
	part, ok = p.ReadSep(':')
	if !ok || string(part) != "b" {
		t.Fatalf("second ReadSep: %q %v", part, ok)
	}
	part, ok = p.ReadSep(',')
	if ok {
		t.Fatalf("expected no sep, got %q", part)
	}
	if !bytes.Equal(part, []byte("c")) {
		t.Fatalf("tail: %q", part)
	}
	if p.Valid() {
		t.Fatal("expected consumed")
	}
}

func TestParser_SkipBytes_tooLong(t *testing.T) {
	t.Parallel()
	p := parser.New([]byte("hi"))
	if p.SkipBytes([]byte("hello")) {
		t.Fatal("expected false")
	}
}

func TestParser_NewString_empty(t *testing.T) {
	t.Parallel()
	p := parser.NewString("")
	if p.Valid() {
		t.Fatal("expected empty")
	}
}

func BenchmarkParser_ReadSep(b *testing.B) {
	data := bytes.Repeat([]byte("a,"), 1024)
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.New(data)
		for p.Valid() {
			_, _ = p.ReadSep(',')
		}
	}
}
