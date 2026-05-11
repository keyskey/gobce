package analyzer

import "errors"

func Analyze(input Input) (Result, error) {
	if input.CoverProfilePath == "" {
		return Result{}, errors.New("coverprofile path is required")
	}

	parsed, err := parseCoverProfile(input.CoverProfilePath)
	if err != nil {
		return Result{}, err
	}
	blocks := parsed.Blocks

	statementCoverage := computeStatementCoverage(blocks)
	candidates, err := collectBranchCandidates(blocks)
	if err != nil {
		return Result{}, err
	}

	covered := 0
	sources := make([]uncoveredSource, 0)
	for _, c := range candidates {
		if c.Covered {
			covered++
			continue
		}

		sources = append(sources, uncoveredSource{
			file: c.FilePath,
			line: c.Line,
			kind: c.Kind,
		})
	}

	rollup := groupUncoveredSources(sources, parsed.ModuleRoot, parsed.ModulePath)

	estimated := 100.0
	if len(candidates) > 0 {
		estimated = percent(float64(covered), float64(len(candidates)))
	}

	return Result{
		Language:                "go",
		StatementCoverage:       statementCoverage,
		EstimatedBranchCoverage: estimated,
		UncoveredBranches:       rollup,
	}, nil
}
