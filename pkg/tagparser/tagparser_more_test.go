package tagparser_test

import (
	"testing"

	"quad4/tagparser/pkg/tagparser"
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
