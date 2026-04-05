// Package lore provides functionality for recording architectural decisions
// as Architecture Decision Records (ADRs) in the docs/adr/ directory.
package lore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/benaskins/axon-snip/analysis"
)

// ADR represents an Architecture Decision Record.
type ADR struct {
	// Number is the ADR number (e.g., 1, 2, 3)
	Number int
	// Title is the short title of the decision
	Title string
	// Status is the current status: "proposed", "accepted", "rejected", "superseded"
	Status string
	// Context describes the context and forces that led to this decision
	Context string
	// Decision describes the decision that was made
	Decision string
	// Consequences are the results of the decision
	Consequences string
	// Sources indicate where this decision came from
	Sources []string
	// Date when the ADR was created
	Date string
}

// Writer handles writing ADRs to disk.
type Writer struct {
	// BaseDir is the base directory for the project
	BaseDir string
	// ADRDir is the directory where ADRs are stored (relative to BaseDir)
	ADRDir string
	// NextNumber is the next ADR number to use
	NextNumber int
}

// NewWriter creates a new ADR writer for the given project directory.
func NewWriter(projectDir string) *Writer {
	return &Writer{
		BaseDir: projectDir,
		ADRDir:  "docs/adr",
	}
}

// Write writes a single ADR to disk.
func (w *Writer) Write(adr ADR) error {
	// Ensure ADR directory exists
	adrPath := filepath.Join(w.BaseDir, w.ADRDir)
	if err := os.MkdirAll(adrPath, 0755); err != nil {
		return fmt.Errorf("creating ADR directory: %w", err)
	}

	// Format the ADR as markdown
	content := w.formatADR(adr)

	// Write to file
	filename := fmt.Sprintf("%03d-%s.md", adr.Number, sanitizeTitle(adr.Title))
	filePath := filepath.Join(adrPath, filename)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing ADR file: %w", err)
	}

	return nil
}

// WriteFromScaffold creates ADRs from a scaffold specification.
// It generates ADRs for:
// - Module selections
// - Boundary classifications
// - Constraint extractions
// - Gap resolutions
func (w *Writer) WriteFromScaffold(spec *analysis.ScaffoldSpec) error {
	// Find the next available ADR number
	nextNum, err := w.findNextNumber()
	if err != nil {
		return fmt.Errorf("finding next ADR number: %w", err)
	}
	w.NextNumber = nextNum

	// Generate ADR for module selections
	if len(spec.Modules) > 0 {
		adr := ADR{
			Number:   w.NextNumber,
			Title:    "Module Selection",
			Status:   "accepted",
			Context:  generateModuleContext(spec),
			Decision: generateModuleDecision(spec),
			Consequences: "The selected modules will be scaffolded with their dependencies. " +
				"This establishes the core architecture of the system.",
			Sources: []string{"PRD analysis", "Module catalogue"},
			Date:    time.Now().Format("2006-01-02"),
		}
		if err := w.Write(adr); err != nil {
			return fmt.Errorf("writing module selection ADR: %w", err)
		}
		w.NextNumber++
	}

	// Generate ADR for boundary classifications
	if len(spec.Boundaries) > 0 {
		adr := ADR{
			Number:   w.NextNumber,
			Title:    "Boundary Classifications",
			Status:   "accepted",
			Context:  generateBoundaryContext(spec.Boundaries),
			Decision: generateBoundaryDecision(spec.Boundaries),
			Consequences: "Boundaries define module interfaces and dependencies. " +
				"This ensures clear separation of concerns.",
			Sources: []string{"Module analysis", "Dependency graph"},
			Date:    time.Now().Format("2006-01-02"),
		}
		if err := w.Write(adr); err != nil {
			return fmt.Errorf("writing boundary ADR: %w", err)
		}
		w.NextNumber++
	}

	// Generate ADR for constraint extractions
	if len(spec.Constraints) > 0 {
		adr := ADR{
			Number:   w.NextNumber,
			Title:    "Architectural Constraints",
			Status:   "accepted",
			Context:  generateConstraintContext(spec.Constraints),
			Decision: generateConstraintDecision(spec.Constraints),
			Consequences: "Constraints guide implementation decisions and ensure " +
				"consistency across the codebase.",
			Sources: []string{"PRD requirements", "Catalogue patterns"},
			Date:    time.Now().Format("2006-01-02"),
		}
		if err := w.Write(adr); err != nil {
			return fmt.Errorf("writing constraint ADR: %w", err)
		}
		w.NextNumber++
	}

	// Generate ADR for gap resolutions
	if len(spec.Gaps) > 0 {
		adr := ADR{
			Number:   w.NextNumber,
			Title:    "Gap Resolutions",
			Status:   "accepted",
			Context:  generateGapContext(spec.Gaps),
			Decision: generateGapDecision(spec.Gaps),
			Consequences: "Resolved gaps provide the missing information needed " +
				"to complete the scaffold.",
			Sources: []string{"Gap analysis", "Expert resolution"},
			Date:    time.Now().Format("2006-01-02"),
		}
		if err := w.Write(adr); err != nil {
			return fmt.Errorf("writing gap resolution ADR: %w", err)
		}
		w.NextNumber++
	}

	return nil
}

// findNextNumber finds the next available ADR number.
func (w *Writer) findNextNumber() (int, error) {
	adrPath := filepath.Join(w.BaseDir, w.ADRDir)
	
	// Check if ADR directory exists
	if _, err := os.Stat(adrPath); os.IsNotExist(err) {
		return 1, nil
	} else if err != nil {
		return 0, err
	}

	// Find existing ADRs and get the max number
	maxNum := 0
	entries, err := os.ReadDir(adrPath)
	if err != nil {
		return 0, fmt.Errorf("reading ADR directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			var num int
			fmt.Sscanf(entry.Name(), "%d-", &num)
			if num > maxNum {
				maxNum = num
			}
		}
	}

	return maxNum + 1, nil
}

// formatADR formats an ADR as markdown.
func (w *Writer) formatADR(adr ADR) string {
	var sb strings.Builder

	sb.WriteString("# ")
	sb.WriteString(fmt.Sprintf("%03d. %s\n", adr.Number, adr.Title))
	sb.WriteString("\n")
	sb.WriteString("## Status\n\n")
	sb.WriteString(fmt.Sprintf("%s\n", adr.Status))
	sb.WriteString("\n")
	sb.WriteString("## Date\n\n")
	sb.WriteString(fmt.Sprintf("%s\n", adr.Date))
	sb.WriteString("\n")

	if len(adr.Sources) > 0 {
		sb.WriteString("## Sources\n\n")
		for _, src := range adr.Sources {
			sb.WriteString(fmt.Sprintf("- %s\n", src))
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

// sanitizeTitle converts a title to a filename-safe string.
func sanitizeTitle(title string) string {
	// Convert to lowercase
	title = strings.ToLower(title)
	// Replace spaces and special chars with hyphens
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

// Helper functions for generating ADR content

func generateModuleContext(spec *analysis.ScaffoldSpec) string {
	var sb strings.Builder
	sb.WriteString("The following modules were selected based on PRD analysis:\n\n")
	for _, mod := range spec.Modules {
		sb.WriteString(fmt.Sprintf("- **%s**: %s\n", mod.Name, mod.Reason))
	}
	return sb.String()
}

func generateModuleDecision(spec *analysis.ScaffoldSpec) string {
	var sb strings.Builder
	sb.WriteString("The following modules have been selected for scaffolding:\n\n")
	for _, mod := range spec.Modules {
		sb.WriteString(fmt.Sprintf("- %s\n", mod.Name))
	}
	return sb.String()
}

func generateBoundaryContext(boundaries []analysis.Boundary) string {
	var sb strings.Builder
	sb.WriteString("Module boundaries and dependencies identified:\n\n")
	for _, b := range boundaries {
		sb.WriteString(fmt.Sprintf("- %s -> %s (%s)\n", b.From, b.To, b.Type))
	}
	return sb.String()
}

func generateBoundaryDecision(boundaries []analysis.Boundary) string {
	return "Module boundaries have been classified based on dependency analysis. " +
		"Each module's dependencies are explicitly declared to ensure clear interfaces."
}

func generateConstraintContext(constraints []string) string {
	var sb strings.Builder
	sb.WriteString("Architectural constraints extracted from specification:\n\n")
	for _, c := range constraints {
		sb.WriteString(fmt.Sprintf("- %s\n", c))
	}
	return sb.String()
}

func generateConstraintDecision(constraints []string) string {
	return "The following constraints have been identified and will guide implementation."
}

func generateGapContext(gaps []analysis.Gap) string {
	var sb strings.Builder
	sb.WriteString("The following gaps were identified during scaffold analysis:\n\n")
	for _, gap := range gaps {
		sb.WriteString(fmt.Sprintf("- **%s**\n", gap.Question))
		if gap.Context != "" {
			sb.WriteString(fmt.Sprintf("  Context: %s\n", gap.Context))
		}
	}
	return sb.String()
}

func generateGapDecision(gaps []analysis.Gap) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d gap(s) were resolved through expert analysis:\n\n", len(gaps)))
	for _, gap := range gaps {
		sb.WriteString(fmt.Sprintf("- %s\n", gap.Question))
	}
	return sb.String()
}
