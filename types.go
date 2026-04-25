package gobce

type AnalyzeInput struct {
	CoverProfilePath string
}

type Result struct {
	Language                string            `json:"language"`
	StatementCoverage       float64           `json:"statementCoverage"`
	EstimatedBranchCoverage float64           `json:"estimatedBranchCoverage"`
	UncoveredBranches       []UncoveredBranch `json:"uncoveredBranches"`
}

type UncoveredBranch struct {
	File           string `json:"file"`
	Line           int    `json:"line"`
	Kind           string `json:"kind"`
	Recommendation string `json:"recommendation,omitempty"`
}
