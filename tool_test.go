package lore

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	tool "github.com/benaskins/axon-tool"
)

func TestNewRecordDecisionTool_Schema(t *testing.T) {
	def := NewRecordDecisionTool(NewWriter(t.TempDir()))

	if def.Name != "record_architectural_decision" {
		t.Errorf("Name: got %q", def.Name)
	}
	for _, req := range []string{"title", "status", "context", "decision"} {
		if !slices.Contains(def.Parameters.Required, req) {
			t.Errorf("Required missing %q", req)
		}
	}
	for _, field := range []string{"title", "status", "context", "decision", "consequences", "sources"} {
		if _, ok := def.Parameters.Properties[field]; !ok {
			t.Errorf("Properties missing %q", field)
		}
	}
	if def.Parameters.Properties["sources"].Items == nil {
		t.Error("sources.Items not set")
	}
}

func TestNewRecordDecisionTool_WritesADR(t *testing.T) {
	dir := t.TempDir()
	def := NewRecordDecisionTool(NewWriter(dir))

	res := def.Execute(&tool.ToolContext{}, map[string]any{
		"title":        "Adopt Event Sourcing",
		"status":       "accepted",
		"context":      "We need an audit trail.",
		"decision":     "Use event sourcing.",
		"consequences": "Harder queries, full history.",
		"sources":      []any{"PRD §4", "RFC 12"},
	})

	if !strings.Contains(res.Content, "recorded ADR") {
		t.Errorf("Content: %q", res.Content)
	}

	path := filepath.Join(dir, "docs", "adr", "001-adopt-event-sourcing.md")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read ADR: %v", err)
	}
	for _, want := range []string{
		"# 001. Adopt Event Sourcing",
		"accepted",
		"We need an audit trail.",
		"Use event sourcing.",
		"- PRD §4",
		"- RFC 12",
	} {
		if !strings.Contains(string(body), want) {
			t.Errorf("ADR missing %q", want)
		}
	}
}

func TestNewRecordDecisionTool_HandlesStringSlice(t *testing.T) {
	dir := t.TempDir()
	def := NewRecordDecisionTool(NewWriter(dir))

	res := def.Execute(&tool.ToolContext{}, map[string]any{
		"title":    "Plain String Slice",
		"status":   "proposed",
		"context":  "c",
		"decision": "d",
		"sources":  []string{"a", "b"},
	})
	if strings.HasPrefix(res.Content, "error") {
		t.Fatalf("unexpected error: %s", res.Content)
	}

	body, err := os.ReadFile(filepath.Join(dir, "docs", "adr", "001-plain-string-slice.md"))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(body), "- a\n- b") {
		t.Errorf("sources not rendered; got:\n%s", body)
	}
}

func TestNewRecordDecisionTool_ReturnsErrorOnWriteFailure(t *testing.T) {
	w := NewWriter("/nonexistent/\x00/bad")
	def := NewRecordDecisionTool(w)

	res := def.Execute(&tool.ToolContext{}, map[string]any{
		"title":    "Doomed",
		"status":   "proposed",
		"context":  "c",
		"decision": "d",
	})
	if !strings.HasPrefix(res.Content, "error") {
		t.Errorf("expected error result, got %q", res.Content)
	}
}

