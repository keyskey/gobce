package analyzer

import "errors"

func Analyze(input Input) (Result, error) {
	if input.CoverProfilePath == "" {
		return Result{}, errors.New("coverprofile path is required")
	}

	blocks, err := parseCoverProfile(input.CoverProfilePath)
	if err != nil {
		return Result{}, err
	}

	statementCoverage := computeStatementCoverage(blocks)
	candidates, err := collectBranchCandidates(blocks)
	if err != nil {
		return Result{}, err
	}

	covered := 0
	uncovered := make([]UncoveredBranch, 0)
	for _, c := range candidates {
		if c.Covered {
			covered++
			continue
		}

		uncovered = append(uncovered, UncoveredBranch{
			File: c.FilePath,
			Line: c.Line,
			Kind: c.Kind,
		})
	}

	estimated := 100.0
	if len(candidates) > 0 {
		estimated = percent(float64(covered), float64(len(candidates)))
	}

	return Result{
		Language:                "go",
		StatementCoverage:       statementCoverage,
		EstimatedBranchCoverage: estimated,
		UncoveredBranches:       uncovered,
	}, nil
}
