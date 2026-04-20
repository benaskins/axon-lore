package lore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriter_Write_CreatesFileWithSlugifiedName(t *testing.T) {
	dir := t.TempDir()
	w := NewWriter(dir)

	err := w.Write(ADR{
		Number:   7,
		Title:    "Use SQLite for Local Storage",
		Status:   "accepted",
		Context:  "Embedded DB needed.",
		Decision: "Adopt SQLite.",
		Date:     "2026-04-20",
	})
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	want := filepath.Join(dir, "docs", "adr", "007-use-sqlite-for-local-storage.md")
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("expected %s: %v", want, err)
	}
}

func TestWriter_Write_FormatsAllSections(t *testing.T) {
	dir := t.TempDir()
	w := NewWriter(dir)

	err := w.Write(ADR{
		Number:       1,
		Title:        "Adopt X",
		Status:       "accepted",
		Context:      "Ctx body.",
		Decision:     "Do X.",
		Consequences: "Better Y.",
		Sources:      []string{"PRD section 3", "spike"},
		Date:         "2026-04-20",
	})
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	body := readFile(t, filepath.Join(dir, "docs", "adr", "001-adopt-x.md"))
	for _, want := range []string{
		"# 001. Adopt X",
		"## Status\n\naccepted",
		"## Date\n\n2026-04-20",
		"## Sources\n\n- PRD section 3\n- spike",
		"## Context\n\nCtx body.",
		"## Decision\n\nDo X.",
		"## Consequences\n\nBetter Y.",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("ADR body missing %q; got:\n%s", want, body)
		}
	}
}

func TestWriter_Write_OmitsEmptyOptionalSections(t *testing.T) {
	dir := t.TempDir()
	w := NewWriter(dir)

	err := w.Write(ADR{
		Number:   1,
		Title:    "Minimal",
		Status:   "proposed",
		Context:  "c",
		Decision: "d",
	})
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	body := readFile(t, filepath.Join(dir, "docs", "adr", "001-minimal.md"))
	for _, unwanted := range []string{"## Date", "## Sources", "## Consequences"} {
		if strings.Contains(body, unwanted) {
			t.Errorf("ADR should omit %q when empty; got:\n%s", unwanted, body)
		}
	}
}

func TestWriter_Write_AutoAssignsNumberWhenZero(t *testing.T) {
	dir := t.TempDir()
	w := NewWriter(dir)

	if err := w.Write(ADR{Number: 5, Title: "First", Status: "accepted", Context: "c", Decision: "d"}); err != nil {
		t.Fatalf("seed write: %v", err)
	}
	if err := w.Write(ADR{Title: "Second", Status: "accepted", Context: "c", Decision: "d"}); err != nil {
		t.Fatalf("auto-numbered write: %v", err)
	}

	want := filepath.Join(dir, "docs", "adr", "006-second.md")
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("expected auto-numbered ADR at %s: %v", want, err)
	}
}

func TestWriter_NextNumber_EmptyDirReturnsOne(t *testing.T) {
	w := NewWriter(t.TempDir())
	n, err := w.NextNumber()
	if err != nil {
		t.Fatalf("NextNumber: %v", err)
	}
	if n != 1 {
		t.Errorf("want 1, got %d", n)
	}
}

func TestWriter_NextNumber_SkipsNonADRFiles(t *testing.T) {
	dir := t.TempDir()
	adrDir := filepath.Join(dir, "docs", "adr")
	if err := os.MkdirAll(adrDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"001-first.md", "003-third.md", "README.md", "notes.txt"} {
		if err := os.WriteFile(filepath.Join(adrDir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	n, err := NewWriter(dir).NextNumber()
	if err != nil {
		t.Fatalf("NextNumber: %v", err)
	}
	if n != 4 {
		t.Errorf("want 4 (max numbered ADR is 003), got %d", n)
	}
}

func TestSanitizeTitle(t *testing.T) {
	cases := map[string]string{
		"Use SQLite":                  "use-sqlite",
		"Adopt OAuth 2.0":             "adopt-oauth-2-0",
		"  Leading/trailing slashes ": "leading-trailing-slashes",
		"ALL CAPS":                    "all-caps",
	}
	for in, want := range cases {
		if got := sanitizeTitle(in); got != want {
			t.Errorf("sanitizeTitle(%q) = %q, want %q", in, got, want)
		}
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}
