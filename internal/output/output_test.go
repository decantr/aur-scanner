package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteTextSafe(t *testing.T) {
	var buf bytes.Buffer
	report := Report{Source: "fixture.md", CheckedPackageCount: 2, CompromisedPackageCount: 1, Safe: true}
	if err := WriteText(&buf, report, TextOptions{}); err != nil {
		t.Fatalf("WriteText returned error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "No compromised installed packages found") {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestWriteTextMatchesQuiet(t *testing.T) {
	var buf bytes.Buffer
	report := Report{Matches: []string{"foo", "bar"}}
	if err := WriteText(&buf, report, TextOptions{Quiet: true}); err != nil {
		t.Fatalf("WriteText returned error: %v", err)
	}
	if got, want := buf.String(), "foo\nbar\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
