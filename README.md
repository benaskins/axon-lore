# axon-lore

Architecture Decision Record (ADR) generation from scaffold specifications.

Import: `github.com/benaskins/axon-lore`

## What it does

axon-lore writes ADRs to `docs/adr/` as numbered markdown files. It integrates with axon-snip's analysis pipeline to automatically document architectural decisions made during scaffolding: module selections, boundary classifications, constraints, and gap resolutions.

## Usage

### From a scaffold spec

```go
import "github.com/benaskins/axon-lore"

writer := lore.NewWriter(projectDir)
err := writer.WriteFromScaffold(spec)
```

This generates separate ADRs for each category of decision in the spec.

### Manual ADR

```go
writer := lore.NewWriter(projectDir)
err := writer.Write(lore.ADR{
    Number: 1,
    Title:  "Use event sourcing for audit trail",
    Status: "accepted",
    Context:    "...",
    Decision:   "...",
    Consequences: "...",
})
```

ADR numbers are auto-incremented by scanning existing files.

## Dependencies

- axon-snip (for `analysis.ScaffoldSpec`)

## Build & Test

```bash
go test ./...
go vet ./...
```
