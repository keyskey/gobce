package analyzer

import (
	"os"
	"path/filepath"
	"strings"
)

func normalizeBlockPaths(blocks []coverageBlock, coverProfilePath string) ([]coverageBlock, error) {
	coverDir := filepath.Dir(coverProfilePath)
	moduleRoot, modulePath := discoverModuleContext(coverDir)

	normalized := make([]coverageBlock, 0, len(blocks))
	for _, b := range blocks {
		resolved := resolveSourcePath(b.FilePath, coverDir, moduleRoot, modulePath)
		b.FilePath = filepath.Clean(resolved)
		normalized = append(normalized, b)
	}
	return normalized, nil
}

func discoverModuleContext(startDir string) (moduleRoot string, modulePath string) {
	dir := filepath.Clean(startDir)
	for {
		goModPath := filepath.Join(dir, "go.mod")
		raw, err := os.ReadFile(goModPath)
		if err == nil {
			return dir, parseModulePath(string(raw))
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ""
		}
		dir = parent
	}
}

func parseModulePath(goModContent string) string {
	for _, rawLine := range strings.Split(goModContent, "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if !strings.HasPrefix(line, "module ") {
			continue
		}
		modulePath := strings.TrimSpace(strings.TrimPrefix(line, "module "))
		modulePath = strings.Trim(modulePath, "\"")
		return modulePath
	}
	return ""
}

func resolveSourcePath(rawPath string, coverDir string, moduleRoot string, modulePath string) string {
	if filepath.IsAbs(rawPath) {
		return rawPath
	}

	relativeCandidate := filepath.Join(coverDir, rawPath)
	if _, err := os.Stat(relativeCandidate); err == nil {
		return relativeCandidate
	}

	if moduleRoot != "" && modulePath != "" {
		modulePrefix := modulePath + "/"
		if rawPath == modulePath {
			return moduleRoot
		}
		if strings.HasPrefix(rawPath, modulePrefix) {
			suffix := strings.TrimPrefix(rawPath, modulePrefix)
			moduleCandidate := filepath.Join(moduleRoot, filepath.FromSlash(suffix))
			if _, err := os.Stat(moduleCandidate); err == nil {
				return moduleCandidate
			}
		}
	}

	return rawPath
}
