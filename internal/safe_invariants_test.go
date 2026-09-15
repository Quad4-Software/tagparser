//go:build appengine
// +build appengine

package internal_test

import (
	"testing"

	"github.com/Quad4-Software/tagparser/v2/internal"
)

// The safe (appengine) variant must copy in both directions: mutating one
// side of a conversion must not leak into the other.
func TestSafeStringToBytes_copies(t *testing.T) {
	t.Parallel()
	s := "hello"
	b := internal.StringToBytes(s)
	b[0] = 'X'
	if s != "hello" {
		t.Fatalf("StringToBytes did not copy: s=%q", s)
	}
}

func TestSafeBytesToString_copies(t *testing.T) {
	t.Parallel()
	b := []byte("hello")
	s := internal.BytesToString(b)
	b[0] = 'X'
	if s != "hello" {
		t.Fatalf("BytesToString did not copy: s=%q", s)
	}
}
