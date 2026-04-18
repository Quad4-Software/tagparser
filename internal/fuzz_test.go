package internal_test

import (
	"bytes"
	"testing"

	"git.quad4.io/Go-Libs/tagparser/v2/internal"
)

func FuzzBytesToString(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte("hello"))
	f.Add([]byte{0, 1, 255})
	f.Fuzz(func(t *testing.T, b []byte) {
		got := internal.BytesToString(b)
		want := string(b)
		if got != want {
			t.Fatalf("BytesToString: got %q want %q", got, want)
		}
	})
}

func FuzzStringToBytes(f *testing.F) {
	f.Add("")
	f.Add("hello")
	f.Add("\x00\xff\x7f")
	f.Fuzz(func(t *testing.T, s string) {
		got := internal.StringToBytes(s)
		want := []byte(s)
		if !bytes.Equal(got, want) {
			t.Fatalf("StringToBytes: mismatch len(got)=%d len(want)=%d", len(got), len(want))
		}
	})
}

// FuzzConvertRoundtrip checks StringToBytes and BytesToString stay consistent with copy semantics.
func FuzzConvertRoundtrip(f *testing.F) {
	f.Add("x")
	f.Add("")
	f.Fuzz(func(t *testing.T, s string) {
		b := internal.StringToBytes(s)
		if string(b) != s {
			t.Fatalf("string(StringToBytes(s)) != s")
		}
		if internal.BytesToString(b) != s {
			t.Fatalf("BytesToString(StringToBytes(s)) != s")
		}
		if !bytes.Equal(b, []byte(s)) {
			t.Fatalf("byte content mismatch")
		}
	})
}
