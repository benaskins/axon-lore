@AGENTS.md

## Conventions
- ADRs are numbered markdown files in docs/adr/ with sanitised title slugs
- Writer.Write assigns the next number when ADR.Number is zero
- Optional sections (Date, Sources, Consequences) are omitted when empty

## Constraints
- Plain data only: no knowledge of PRDs, scaffolds, or any caller domain
- Zero dependency on axon-snip or any scaffolding-specific package
- axon-tool is the only axon-* dependency (for ToolDef)
- Do not add LLM calls; this package is entirely deterministic

## Surfaces
- `Writer.Write(ADR)` for direct/deterministic callers
- `NewRecordDecisionTool(*Writer) *tool.ToolDef` for LLM or direct dispatch

## Testing
- `go test ./...` runs all tests
- `go vet ./...` for lint
