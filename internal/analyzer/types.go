package analyzer

type Input struct {
	CoverProfilePath string
}

type Result struct {
	Language                string
	StatementCoverage       float64
	EstimatedBranchCoverage float64
	UncoveredBranches       []UncoveredBranch
}

type UncoveredBranch struct {
	File           string
	Line           int
	Kind           string
	Recommendation string
}

type coverageBlock struct {
	FilePath  string
	StartLine int
	EndLine   int
	NumStmts  int
	Count     int
}

type branchCandidate struct {
	FilePath string
	Line     int
	Kind     string
	Covered  bool
}
