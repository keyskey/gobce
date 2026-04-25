package gobce

import "github.com/keyskey/gobce/internal/analyzer"

func Analyze(input AnalyzeInput) (Result, error) {
	analysisResult, err := analyzer.Analyze(analyzer.Input{
		CoverProfilePath: input.CoverProfilePath,
	})
	if err != nil {
		return Result{}, err
	}

	uncovered := make([]UncoveredBranch, 0, len(analysisResult.UncoveredBranches))
	for _, b := range analysisResult.UncoveredBranches {
		uncovered = append(uncovered, UncoveredBranch{
			File:           b.File,
			Line:           b.Line,
			Kind:           b.Kind,
			Recommendation: b.Recommendation,
		})
	}

	return Result{
		Language:                analysisResult.Language,
		StatementCoverage:       analysisResult.StatementCoverage,
		EstimatedBranchCoverage: analysisResult.EstimatedBranchCoverage,
		UncoveredBranches:       uncovered,
	}, nil
}
