package source

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseCompromisedPackagesFromMarkdownFixture(t *testing.T) {
	data := mustReadFixture(t, "compromised-list.md")
	packages, err := ParseCompromisedPackages(data)
	if err != nil {
		t.Fatalf("ParseCompromisedPackages error: %v", err)
	}

	want := []string{"foo"}
	assertContainsExactly(t, packages, want)
}

func TestParseCompromisedPackagesFromHTMLFixture(t *testing.T) {
	data := mustReadFixture(t, "compromised-page.html")
	packages, err := ParseCompromisedPackages(data)
	if err != nil {
		t.Fatalf("ParseCompromisedPackages error: %v", err)
	}

	want := []string{"foo"}
	assertContainsExactly(t, packages, want)
}

func TestParseCompromisedPackagesRejectsEmptyInput(t *testing.T) {
	if _, err := ParseCompromisedPackages([]byte("")); err == nil {
		t.Fatal("expected an error for empty input")
	}
}

func mustReadFixture(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "..", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", name, err)
	}
	return data
}

func assertContainsExactly(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d packages, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
