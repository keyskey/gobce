package gobce

type AnalyzeInput struct {
	CoverProfilePath string
}

type Result struct {
	Language                string                  `json:"language"`
	StatementCoverage       float64                 `json:"statementCoverage"`
	EstimatedBranchCoverage float64                 `json:"estimatedBranchCoverage"`
	UncoveredBranches       UncoveredBranchesReport `json:"uncoveredBranches"`
}

// UncoveredBranchesReport groups uncovered branch findings by import path and file path.
type UncoveredBranchesReport struct {
	Packages []UncoveredPackage `json:"packages"`
}

type UncoveredPackage struct {
	ImportPath string          `json:"importPath"`
	Files      []UncoveredFile `json:"files"`
}

type UncoveredFile struct {
	Path     string            `json:"path"`
	Branches []UncoveredBranch `json:"branches"`
}

type UncoveredBranch struct {
	Line           int    `json:"line"`
	Kind           string `json:"kind"`
	Recommendation string `json:"recommendation,omitempty"`
}
