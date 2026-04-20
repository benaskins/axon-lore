package lore

import (
	"fmt"
	"time"

	tool "github.com/benaskins/axon-tool"
)

// NewRecordDecisionTool returns a ToolDef that writes an ADR through the given
// Writer. LLM harnesses (via axon-loop) and direct callers can invoke it
// through the standard axon-tool dispatch surface.
//
// The tool assigns the next available ADR number and today's date when the
// caller omits them.
func NewRecordDecisionTool(w *Writer) *tool.ToolDef {
	return &tool.ToolDef{
		Name:        "record_architectural_decision",
		Description: "Record an architectural decision as an ADR in docs/adr/",
		Parameters: tool.ParameterSchema{
			Type:     "object",
			Required: []string{"title", "status", "context", "decision"},
			Properties: map[string]tool.PropertySchema{
				"title": {
					Type:        "string",
					Description: "Short title of the decision",
				},
				"status": {
					Type:        "string",
					Description: "ADR lifecycle state",
					Enum:        []any{"proposed", "accepted", "rejected", "superseded"},
				},
				"context": {
					Type:        "string",
					Description: "Forces and constraints that led to this decision",
				},
				"decision": {
					Type:        "string",
					Description: "The decision itself",
				},
				"consequences": {
					Type:        "string",
					Description: "Results and trade-offs of the decision",
				},
				"sources": {
					Type:        "array",
					Description: "References that informed the decision",
					Items:       &tool.PropertySchema{Type: "string"},
				},
			},
		},
		Execute: func(_ *tool.ToolContext, args map[string]any) tool.ToolResult {
			adr := ADR{
				Title:        stringArg(args, "title"),
				Status:       stringArg(args, "status"),
				Context:      stringArg(args, "context"),
				Decision:     stringArg(args, "decision"),
				Consequences: stringArg(args, "consequences"),
				Sources:      stringSliceArg(args, "sources"),
				Date:         time.Now().UTC().Format("2006-01-02"),
			}
			if err := w.Write(adr); err != nil {
				return tool.ToolResult{Content: fmt.Sprintf("error: %s", err)}
			}
			return tool.ToolResult{
				Content: fmt.Sprintf("recorded ADR: %s", adr.Title),
			}
		},
	}
}

func stringArg(args map[string]any, key string) string {
	v, ok := args[key]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

func stringSliceArg(args map[string]any, key string) []string {
	v, ok := args[key]
	if !ok {
		return nil
	}
	switch xs := v.(type) {
	case []string:
		return xs
	case []any:
		out := make([]string, 0, len(xs))
		for _, x := range xs {
			if s, ok := x.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}
