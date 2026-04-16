@AGENTS.md

## Conventions
- ADRs are numbered markdown files in docs/adr/ with sanitised title slugs
- WriteFromScaffold generates multiple ADRs from a single ScaffoldSpec
- ADR numbers auto-increment by scanning existing files

## Constraints
- Depends on axon-snip only (for analysis.ScaffoldSpec); no other axon-* deps
- ADR format follows the standard template: Status, Date, Sources, Context, Decision, Consequences
- Do not add LLM calls; this package is entirely deterministic

## Testing
- `go test ./...` runs all tests
- `go vet ./...` for lint
