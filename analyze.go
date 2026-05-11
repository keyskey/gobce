package gobce

import "github.com/keyskey/gobce/internal/analyzer"

func Analyze(input AnalyzeInput) (Result, error) {
	analysisResult, err := analyzer.Analyze(analyzer.Input{
		CoverProfilePath: input.CoverProfilePath,
	})
	if err != nil {
		return Result{}, err
	}

	return Result{
		Language:                analysisResult.Language,
		StatementCoverage:       analysisResult.StatementCoverage,
		EstimatedBranchCoverage: analysisResult.EstimatedBranchCoverage,
		UncoveredBranches:       convertUncoveredRollup(analysisResult.UncoveredBranches),
	}, nil
}

func convertUncoveredRollup(in analyzer.UncoveredBranchesReport) UncoveredBranchesReport {
	pkgs := make([]UncoveredPackage, len(in.Packages))
	for i, p := range in.Packages {
		files := make([]UncoveredFile, len(p.Files))
		for j, f := range p.Files {
			br := make([]UncoveredBranch, len(f.Branches))
			for k, b := range f.Branches {
				br[k] = UncoveredBranch{
					Line:           b.Line,
					Kind:           b.Kind,
					Recommendation: b.Recommendation,
				}
			}
			files[j] = UncoveredFile{Path: f.Path, Branches: br}
		}
		pkgs[i] = UncoveredPackage{ImportPath: p.ImportPath, Files: files}
	}
	return UncoveredBranchesReport{Packages: pkgs}
}
