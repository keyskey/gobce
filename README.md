# gobce

gobce (Go Branch Coverage Estimator) is a standalone CLI that estimates C1 branch coverage for Go projects.
It combines `go test -coverprofile` output with lightweight AST/CFG analysis to produce:
- statement coverage
- estimated C1 branch coverage
- uncovered branch findings

> Note: gobce reports **estimated C1**, not exact branch coverage, because Go coverprofiles are statement/block-oriented.