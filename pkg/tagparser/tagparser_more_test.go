package tagparser_test

import (
	"testing"

	"github.com/Quad4-Software/tagparser/pkg/tagparser"
)

func TestTag_HasOption(t *testing.T) {
	tag := tagparser.Parse("n,a:1,b:,c:value")
	if tag.Name != "n" {
		t.Fatalf("name: got %q", tag.Name)
	}
	for _, key := range []string{"a", "b", "c"} {
		if !tag.HasOption(key) {
			t.Fatalf("expected HasOption(%q)", key)
		}
	}
	if tag.HasOption("missing") {
		t.Fatal("unexpected HasOption(missing)")
	}
	if tag.Options["a"] != "1" || tag.Options["b"] != "" || tag.Options["c"] != "value" {
		t.Fatalf("options: %#v", tag.Options)
	}
}

func TestTag_HasOption_nilOptions(t *testing.T) {
	for _, s := range []string{"", "name", "'quoted name'"} {
		tag := tagparser.Parse(s)
		if tag.Options != nil {
			t.Fatalf("expected nil options for %q, got %#v", s, tag.Options)
		}
		if tag.HasOption("x") {
			t.Fatalf("unexpected HasOption on nil options (tag=%q)", s)
		}
	}
}

func TestVersion(t *testing.T) {
	if tagparser.Version() == "" {
		t.Fatal("empty version")
	}
}
