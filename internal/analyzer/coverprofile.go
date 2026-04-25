package analyzer

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func parseCoverProfile(path string) ([]coverageBlock, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open coverprofile: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	blocks := make([]coverageBlock, 0)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if lineNo == 1 && strings.HasPrefix(line, "mode:") {
			continue
		}

		b, err := parseCoverProfileLine(line)
		if err != nil {
			return nil, fmt.Errorf("parse coverprofile line %d: %w", lineNo, err)
		}
		blocks = append(blocks, b)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan coverprofile: %w", err)
	}

	if len(blocks) == 0 {
		return nil, errors.New("coverprofile contains no coverage blocks")
	}

	normalizedBlocks, err := normalizeBlockPaths(blocks, path)
	if err != nil {
		return nil, err
	}
	return normalizedBlocks, nil
}

func parseCoverProfileLine(line string) (coverageBlock, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return coverageBlock{}, fmt.Errorf("unexpected coverprofile entry: %q", line)
	}

	fileAndSpan := parts[0]
	stmts, err := strconv.Atoi(parts[1])
	if err != nil {
		return coverageBlock{}, fmt.Errorf("invalid num statements: %w", err)
	}
	count, err := strconv.Atoi(parts[2])
	if err != nil {
		return coverageBlock{}, fmt.Errorf("invalid execution count: %w", err)
	}

	colon := strings.LastIndex(fileAndSpan, ":")
	if colon == -1 {
		return coverageBlock{}, fmt.Errorf("missing ':' in span: %q", fileAndSpan)
	}
	filePath := fileAndSpan[:colon]
	span := fileAndSpan[colon+1:]

	spanParts := strings.Split(span, ",")
	if len(spanParts) != 2 {
		return coverageBlock{}, fmt.Errorf("invalid span: %q", span)
	}

	startLine, err := parseSpanLine(spanParts[0])
	if err != nil {
		return coverageBlock{}, err
	}
	endLine, err := parseSpanLine(spanParts[1])
	if err != nil {
		return coverageBlock{}, err
	}

	return coverageBlock{
		FilePath:  filepath.Clean(filePath),
		StartLine: startLine,
		EndLine:   endLine,
		NumStmts:  stmts,
		Count:     count,
	}, nil
}

func parseSpanLine(pos string) (int, error) {
	dot := strings.Index(pos, ".")
	if dot == -1 {
		return 0, fmt.Errorf("invalid position format: %q", pos)
	}
	line, err := strconv.Atoi(pos[:dot])
	if err != nil {
		return 0, fmt.Errorf("invalid line in position: %w", err)
	}
	return line, nil
}
