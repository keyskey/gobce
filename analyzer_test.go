package gobce

import (
	"os"
	"path/filepath"
	"strconv"
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

func TestAnalyzeOneLineIfElseColumnAwareCoverage(t *testing.T) {
	tmp := t.TempDir()
	srcPath := filepath.Join(tmp, "oneline.go")
	src := `package sample

func score(v int) int { if v > 10 { return 1 } else { return 2 } }
`
	if err := os.WriteFile(srcPath, []byte(src), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	coverPath := filepath.Join(tmp, "coverage.out")
	coverage := strings.Join([]string{
		"mode: set",
		srcPath + ":3.31,3.41 1 1",
		srcPath + ":3.50,3.65 1 0",
	}, "\n")
	if err := os.WriteFile(coverPath, []byte(coverage), 0o644); err != nil {
		t.Fatalf("write coverage: %v", err)
	}

	result, err := Analyze(AnalyzeInput{CoverProfilePath: coverPath})
	if err != nil {
		t.Fatalf("analyze: %v", err)
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

func TestAnalyzeSkipsUnresolvableSourcePaths(t *testing.T) {
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
		"github.com/example/dependency/pkg/generated.go:1.1,1.10 1 1",
	}, "\n")
	if err := os.WriteFile(coverPath, []byte(coverage), 0o644); err != nil {
		t.Fatalf("write coverage: %v", err)
	}

	result, err := Analyze(AnalyzeInput{CoverProfilePath: coverPath})
	if err != nil {
		t.Fatalf("analyze should not fail on unresolved source path: %v", err)
	}

	if len(result.UncoveredBranches) == 0 {
		t.Fatalf("expected uncovered branches from resolvable source file")
	}
}

func TestAnalyzeIfWithoutElseTerminatingBody(t *testing.T) {
	tests := []struct {
		name                       string
		ifBodyCount                int
		fallthroughCount           int
		wantUncoveredTruePath      bool
		wantUncoveredImplicitFalse bool
	}{
		{
			name:                       "only_true_path_covered",
			ifBodyCount:                1,
			fallthroughCount:           0,
			wantUncoveredImplicitFalse: true,
		},
		{
			name:                  "only_false_path_covered",
			ifBodyCount:           0,
			fallthroughCount:      1,
			wantUncoveredTruePath: true,
		},
		{
			name:             "both_paths_covered",
			ifBodyCount:      1,
			fallthroughCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			srcPath := filepath.Join(tmp, "sample.go")
			src := `package sample

func score(v int) int {
	if v > 10 {
		return 1
	}
	return 2
}
`
			if err := os.WriteFile(srcPath, []byte(src), 0o644); err != nil {
				t.Fatalf("write source: %v", err)
			}

			coverPath := filepath.Join(tmp, "coverage.out")
			coverage := strings.Join([]string{
				"mode: set",
				srcPath + ":3.23,4.13 1 1",
				srcPath + ":4.13,6.3 1 " + strconv.Itoa(tt.ifBodyCount),
				srcPath + ":6.3,7.10 1 " + strconv.Itoa(tt.fallthroughCount),
			}, "\n")
			if err := os.WriteFile(coverPath, []byte(coverage), 0o644); err != nil {
				t.Fatalf("write coverage: %v", err)
			}

			result, err := Analyze(AnalyzeInput{CoverProfilePath: coverPath})
			if err != nil {
				t.Fatalf("analyze: %v", err)
			}

			gotUncoveredTruePath := hasUncoveredBranchKind(result, "if_true_path")
			if gotUncoveredTruePath != tt.wantUncoveredTruePath {
				t.Fatalf("if_true_path uncovered: got %v, want %v", gotUncoveredTruePath, tt.wantUncoveredTruePath)
			}

			gotUncoveredImplicitFalse := hasUncoveredBranchKind(result, "if_implicit_false_path")
			if gotUncoveredImplicitFalse != tt.wantUncoveredImplicitFalse {
				t.Fatalf("if_implicit_false_path uncovered: got %v, want %v", gotUncoveredImplicitFalse, tt.wantUncoveredImplicitFalse)
			}
		})
	}
}

func TestAnalyzeIfWithoutElseNonTerminatingBody(t *testing.T) {
	tmp := t.TempDir()
	srcPath := filepath.Join(tmp, "sample.go")
	src := `package sample

func score(v int) int {
	score := 0
	if v > 10 {
		score = 1
	}
	return score
}
`
	if err := os.WriteFile(srcPath, []byte(src), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	coverPath := filepath.Join(tmp, "coverage.out")
	coverage := strings.Join([]string{
		"mode: set",
		srcPath + ":3.23,5.13 2 1",
		srcPath + ":5.13,7.3 1 0",
		srcPath + ":7.3,8.14 1 1",
	}, "\n")
	if err := os.WriteFile(coverPath, []byte(coverage), 0o644); err != nil {
		t.Fatalf("write coverage: %v", err)
	}

	result, err := Analyze(AnalyzeInput{CoverProfilePath: coverPath})
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}

	if hasUncoveredBranchKind(result, "if_implicit_false_path") {
		t.Fatalf("did not expect if_implicit_false_path for non-terminating if body")
	}
}

func TestAnalyzeIfWithoutElsePanicBody(t *testing.T) {
	tmp := t.TempDir()
	srcPath := filepath.Join(tmp, "sample.go")
	src := `package sample

func score(v int) int {
	if v > 10 {
		panic("too high")
	}
	return 2
}
`
	if err := os.WriteFile(srcPath, []byte(src), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	coverPath := filepath.Join(tmp, "coverage.out")
	coverage := strings.Join([]string{
		"mode: set",
		srcPath + ":3.23,4.13 1 1",
		srcPath + ":4.13,6.3 1 1",
		srcPath + ":6.3,7.10 1 0",
	}, "\n")
	if err := os.WriteFile(coverPath, []byte(coverage), 0o644); err != nil {
		t.Fatalf("write coverage: %v", err)
	}

	result, err := Analyze(AnalyzeInput{CoverProfilePath: coverPath})
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}

	if !hasUncoveredBranchKind(result, "if_implicit_false_path") {
		t.Fatalf("expected if_implicit_false_path in uncovered branches")
	}
}

func hasUncoveredBranchKind(result Result, kind string) bool {
	for _, b := range result.UncoveredBranches {
		if b.Kind == kind {
			return true
		}
	}
	return false
}
