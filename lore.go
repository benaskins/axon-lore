// Package lore provides functionality for recording architectural decisions
// as Architecture Decision Records (ADRs) in the docs/adr/ directory.
package lore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ADR represents an Architecture Decision Record.
type ADR struct {
	// Number is the ADR number (e.g., 1, 2, 3). When zero, Writer.Write assigns
	// the next available number automatically.
	Number int
	// Title is the short title of the decision.
	Title string
	// Status is the current status: "proposed", "accepted", "rejected", "superseded".
	Status string
	// Context describes the context and forces that led to this decision.
	Context string
	// Decision describes the decision that was made.
	Decision string
	// Consequences are the results of the decision.
	Consequences string
	// Sources indicate where this decision came from.
	Sources []string
	// Date when the ADR was created (YYYY-MM-DD).
	Date string
}

// Writer handles writing ADRs to disk.
type Writer struct {
	// BaseDir is the base directory for the project.
	BaseDir string
	// ADRDir is the directory where ADRs are stored (relative to BaseDir).
	ADRDir string
}

// NewWriter creates a new ADR writer for the given project directory.
func NewWriter(projectDir string) *Writer {
	return &Writer{
		BaseDir: projectDir,
		ADRDir:  "docs/adr",
	}
}

// Write writes a single ADR to disk. If adr.Number is zero, the next available
// number is assigned by scanning the ADR directory.
func (w *Writer) Write(adr ADR) error {
	adrPath := filepath.Join(w.BaseDir, w.ADRDir)
	if err := os.MkdirAll(adrPath, 0o755); err != nil {
		return fmt.Errorf("creating ADR directory: %w", err)
	}

	if adr.Number == 0 {
		next, err := w.NextNumber()
		if err != nil {
			return fmt.Errorf("finding next ADR number: %w", err)
		}
		adr.Number = next
	}

	content := formatADR(adr)

	filename := fmt.Sprintf("%03d-%s.md", adr.Number, sanitizeTitle(adr.Title))
	filePath := filepath.Join(adrPath, filename)
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing ADR file: %w", err)
	}

	return nil
}

// NextNumber returns the next available ADR number by scanning the ADR
// directory for existing numbered markdown files.
func (w *Writer) NextNumber() (int, error) {
	adrPath := filepath.Join(w.BaseDir, w.ADRDir)

	if _, err := os.Stat(adrPath); os.IsNotExist(err) {
		return 1, nil
	} else if err != nil {
		return 0, err
	}

	maxNum := 0
	entries, err := os.ReadDir(adrPath)
	if err != nil {
		return 0, fmt.Errorf("reading ADR directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		var num int
		if _, err := fmt.Sscanf(entry.Name(), "%d-", &num); err == nil && num > maxNum {
			maxNum = num
		}
	}

	return maxNum + 1, nil
}

// formatADR formats an ADR as markdown.
func formatADR(adr ADR) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "# %03d. %s\n\n", adr.Number, adr.Title)
	sb.WriteString("## Status\n\n")
	sb.WriteString(adr.Status)
	sb.WriteString("\n\n")

	if adr.Date != "" {
		sb.WriteString("## Date\n\n")
		sb.WriteString(adr.Date)
		sb.WriteString("\n\n")
	}

	if len(adr.Sources) > 0 {
		sb.WriteString("## Sources\n\n")
		for _, src := range adr.Sources {
			fmt.Fprintf(&sb, "- %s\n", src)
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Context\n\n")
	sb.WriteString(adr.Context)
	sb.WriteString("\n\n")

	sb.WriteString("## Decision\n\n")
	sb.WriteString(adr.Decision)
	sb.WriteString("\n\n")

	if adr.Consequences != "" {
		sb.WriteString("## Consequences\n\n")
		sb.WriteString(adr.Consequences)
		sb.WriteString("\n\n")
	}

	return sb.String()
}

// sanitizeTitle converts a title to a filename-safe slug.
func sanitizeTitle(title string) string {
	title = strings.ToLower(title)
	var sb strings.Builder
	for _, r := range title {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		} else {
			sb.WriteRune('-')
		}
	}
	return strings.Trim(sb.String(), "-")
}
