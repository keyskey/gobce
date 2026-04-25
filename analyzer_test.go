package gobce

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAnalyze(t *testing.T) {
	tmp := t.TempDir()
	srcPath := filepath.Join(tmp, "sample.go")
	src := `package sample

func score(v int) int {
	if v > 10 {
		return 1
	} else {
		return 2
	}
}
`
	if err := os.WriteFile(srcPath, []byte(src), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	coverPath := filepath.Join(tmp, "coverage.out")
	coverage := strings.Join([]string{
		"mode: set",
		srcPath + ":3.23,4.13 1 1",
		srcPath + ":4.13,5.3 1 1",
		srcPath + ":6.8,8.3 1 0",
	}, "\n")
	if err := os.WriteFile(coverPath, []byte(coverage), 0o644); err != nil {
		t.Fatalf("write coverage: %v", err)
	}

	result, err := Analyze(AnalyzeInput{CoverProfilePath: coverPath})
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}

	if result.Language != "go" {
		t.Fatalf("language: got %q", result.Language)
	}
	if len(result.UncoveredBranches) == 0 {
		t.Fatalf("expected uncovered branches")
	}

	var hasIfFalse bool
	for _, b := range result.UncoveredBranches {
		if b.Kind == "if_false_path" {
			hasIfFalse = true
			break
		}
	}
	if !hasIfFalse {
		t.Fatalf("expected if_false_path in uncovered branches")
	}
}
