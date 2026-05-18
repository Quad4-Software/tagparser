package internal_test

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/Quad4-Software/tagparser/internal"
)

func TestBytesToString_equivalence(t *testing.T) {
	t.Parallel()
	cases := [][]byte{
		nil,
		{},
		[]byte("ascii"),
		[]byte("caf\xe9"),
		bytes.Repeat([]byte("x"), 1024),
	}
	for i, b := range cases {
		b := b
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			got := internal.BytesToString(b)
			want := string(b)
			if got != want {
				t.Fatalf("BytesToString: got %q want %q", got, want)
			}
		})
	}
}

func TestStringToBytes_equivalence(t *testing.T) {
	t.Parallel()
	cases := []string{
		"",
		"a",
		"hello, 世界",
		strings.Repeat("z", 4096),
	}
	for i, s := range cases {
		s := s
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			got := internal.StringToBytes(s)
			want := []byte(s)
			if !bytes.Equal(got, want) {
				t.Fatalf("StringToBytes: got %q want %q", got, want)
			}
		})
	}
}

func TestStringToBytes_roundtripString(t *testing.T) {
	t.Parallel()
	s := "roundtrip"
	b := internal.StringToBytes(s)
	if string(b) != s {
		t.Fatalf("string(StringToBytes(s)): got %q want %q", string(b), s)
	}
}

func TestBytesToString_roundtripBytes(t *testing.T) {
	t.Parallel()
	b := []byte("bytes")
	if !bytes.Equal(internal.StringToBytes(internal.BytesToString(b)), b) {
		t.Fatal("StringToBytes(BytesToString(b)) != b")
	}
}
