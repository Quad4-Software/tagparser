//go:build !appengine
// +build !appengine

package internal_test

import (
	"strconv"
	"testing"

	"git.quad4.io/Go-Libs/tagparser/v2/internal"
)

// StringToBytes for the non-appengine build uses a slice header with len==cap==len(s).
func TestUnsafeStringToBytes_lenEqualsCap(t *testing.T) {
	t.Parallel()
	for i, s := range []string{"", "a", "hello", string(rune(0x1F600))} {
		s := s
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			b := internal.StringToBytes(s)
			if len(b) != len(s) {
				t.Fatalf("len %d != len(s) %d", len(b), len(s))
			}
			if cap(b) != len(s) {
				t.Fatalf("cap %d != len(s) %d (unsafe path should set cap=len)", cap(b), len(s))
			}
		})
	}
}
