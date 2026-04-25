package analyzer

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func collectBranchCandidates(blocks []coverageBlock) ([]branchCandidate, error) {
	byFile := map[string][]coverageBlock{}
	for _, b := range blocks {
		byFile[b.FilePath] = append(byFile[b.FilePath], b)
	}

	result := make([]branchCandidate, 0)
	for filePath, fileBlocks := range byFile {
		candidates, err := collectFileBranchCandidates(filePath, fileBlocks)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, err
		}
		result = append(result, candidates...)
	}
	return result, nil
}

func collectFileBranchCandidates(filePath string, blocks []coverageBlock) ([]branchCandidate, error) {
	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parse source %q: %w", filePath, err)
	}

	candidates := make([]branchCandidate, 0)
	ast.Inspect(parsed, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.IfStmt:
			candidates = append(candidates, branchCandidate{
				FilePath: filePath,
				Line:     fset.Position(n.Body.Lbrace).Line,
				Kind:     "if_true_path",
				Covered:  spanCovered(blocks, fset, n.Body.Pos(), n.Body.End()),
			})

			if n.Else != nil {
				candidates = append(candidates, branchCandidate{
					FilePath: filePath,
					Line:     fset.Position(n.Else.Pos()).Line,
					Kind:     "if_false_path",
					Covered:  spanCovered(blocks, fset, n.Else.Pos(), n.Else.End()),
				})
			}

		case *ast.SwitchStmt:
			candidates = append(candidates, collectCaseClauseCandidates(filePath, fset, blocks, n.Body.List, "switch_case_path", "switch_default_path")...)
		case *ast.TypeSwitchStmt:
			candidates = append(candidates, collectCaseClauseCandidates(filePath, fset, blocks, n.Body.List, "type_switch_case_path", "type_switch_default_path")...)
		case *ast.ForStmt:
			bodyCovered := spanCovered(blocks, fset, n.Body.Pos(), n.Body.End())
			stmtCovered := spanCovered(blocks, fset, n.Pos(), n.Body.Lbrace)
			candidates = append(candidates, branchCandidate{
				FilePath: filePath,
				Line:     fset.Position(n.Body.Lbrace).Line,
				Kind:     "for_body_entered",
				Covered:  bodyCovered,
			})
			candidates = append(candidates, branchCandidate{
				FilePath: filePath,
				Line:     fset.Position(n.For).Line,
				Kind:     "for_body_not_entered",
				Covered:  stmtCovered && !bodyCovered,
			})
		case *ast.RangeStmt:
			bodyCovered := spanCovered(blocks, fset, n.Body.Pos(), n.Body.End())
			stmtCovered := spanCovered(blocks, fset, n.Pos(), n.Body.Lbrace)
			candidates = append(candidates, branchCandidate{
				FilePath: filePath,
				Line:     fset.Position(n.Body.Lbrace).Line,
				Kind:     "range_body_entered",
				Covered:  bodyCovered,
			})
			candidates = append(candidates, branchCandidate{
				FilePath: filePath,
				Line:     fset.Position(n.For).Line,
				Kind:     "range_body_not_entered",
				Covered:  stmtCovered && !bodyCovered,
			})
		}
		return true
	})

	return candidates, nil
}

func collectCaseClauseCandidates(filePath string, fset *token.FileSet, blocks []coverageBlock, clauses []ast.Stmt, caseKind string, defaultKind string) []branchCandidate {
	candidates := make([]branchCandidate, 0, len(clauses))
	for _, stmt := range clauses {
		clause, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}
		kind := caseKind
		if len(clause.List) == 0 {
			kind = defaultKind
		}
		candidates = append(candidates, branchCandidate{
			FilePath: filePath,
			Line:     fset.Position(clause.Case).Line,
			Kind:     kind,
			Covered:  spanCovered(blocks, fset, clause.Pos(), clause.End()),
		})
	}
	return candidates
}

func spanCovered(blocks []coverageBlock, fset *token.FileSet, start token.Pos, end token.Pos) bool {
	startPos := fset.Position(start)
	endPos := fset.Position(end)

	for _, b := range blocks {
		if b.Count <= 0 {
			continue
		}
		if spansOverlap(startPos.Line, startPos.Column, endPos.Line, endPos.Column, b.StartLine, b.StartCol, b.EndLine, b.EndCol) {
			return true
		}
	}
	return false
}

func spansOverlap(aStartLine, aStartCol, aEndLine, aEndCol, bStartLine, bStartCol, bEndLine, bEndCol int) bool {
	return positionLess(aStartLine, aStartCol, bEndLine, bEndCol) &&
		positionLess(bStartLine, bStartCol, aEndLine, aEndCol)
}

func positionLess(lineA, colA, lineB, colB int) bool {
	if lineA != lineB {
		return lineA < lineB
	}
	return colA < colB
}
