# axon-lore

Architecture Decision Record (ADR) writer. Plain-data, zero scaffolding-dialect dependencies.

Import: `github.com/benaskins/axon-lore`

## What it does

axon-lore writes ADRs to `docs/adr/` as numbered markdown files. It exposes
two surfaces:

- A direct Go API (`Writer.Write`) for deterministic callers.
- A `ToolDef` (`NewRecordDecisionTool`) for LLMs and callers that dispatch
  through axon-tool.

axon-lore knows nothing about PRDs, scaffolds, or modules. Callers shape their
domain data into ADR fields and hand them over.

## Usage

### Direct

```go
import "github.com/benaskins/axon-lore"

w := lore.NewWriter(projectDir)
err := w.Write(lore.ADR{
    Title:        "Use SQLite for local storage",
    Status:       "accepted",
    Context:      "We need an embedded database with zero ops cost.",
    Decision:     "Adopt SQLite with WAL journaling.",
    Consequences: "No external DB; single-writer concurrency.",
    Date:         "2026-04-20",
})
```

ADR numbers auto-increment when `Number` is zero. Numbers are assigned by
scanning existing files in the ADR directory.

### Via tool dispatch

```go
import (
    lore "github.com/benaskins/axon-lore"
    tool "github.com/benaskins/axon-tool"
)

w := lore.NewWriter(projectDir)
def := lore.NewRecordDecisionTool(w)

// Callers (axon-loop, a direct dispatcher) invoke def.Execute with a map[string]any.
res := def.Execute(&tool.ToolContext{}, map[string]any{
    "title":    "Adopt event sourcing",
    "status":   "accepted",
    "context":  "Audit trail is a hard requirement.",
    "decision": "Use event sourcing for state changes.",
    "sources":  []any{"PRD §4", "spike notes"},
})
```

## Dependencies

- `axon-tool` (for `ToolDef`)

## Build & Test

```bash
go test ./...
go vet ./...
```
