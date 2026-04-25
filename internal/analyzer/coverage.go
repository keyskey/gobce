package analyzer

func computeStatementCoverage(blocks []coverageBlock) float64 {
	total := 0
	covered := 0
	for _, b := range blocks {
		total += b.NumStmts
		if b.Count > 0 {
			covered += b.NumStmts
		}
	}
	if total == 0 {
		return 0
	}
	return percent(float64(covered), float64(total))
}

func percent(numerator, denominator float64) float64 {
	return (numerator / denominator) * 100
}
