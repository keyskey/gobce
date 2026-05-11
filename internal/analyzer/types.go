package analyzer

type Input struct {
	CoverProfilePath string
}

type Result struct {
	Language                string
	StatementCoverage       float64
	EstimatedBranchCoverage float64
	UncoveredBranches       UncoveredBranchesReport
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

type coverageBlock struct {
	FilePath  string
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
	NumStmts  int
	Count     int
}

type branchCandidate struct {
	FilePath string
	Line     int
	Kind     string
	Covered  bool
}
