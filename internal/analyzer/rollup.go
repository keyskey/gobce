package analyzer

import (
	"path/filepath"
	"sort"
	"strings"
)

type uncoveredSource struct {
	file           string
	line           int
	kind           string
	recommendation string
}

func groupUncoveredSources(rows []uncoveredSource, moduleRoot, modulePath string) UncoveredBranchesReport {
	if len(rows) == 0 {
		return UncoveredBranchesReport{Packages: []UncoveredPackage{}}
	}

	type fileKey string
	pkgFiles := map[string]map[fileKey][]UncoveredBranch{}

	for _, row := range rows {
		importPath, displayPath := classifySourceFile(row.file, moduleRoot, modulePath)
		if pkgFiles[importPath] == nil {
			pkgFiles[importPath] = map[fileKey][]UncoveredBranch{}
		}
		k := fileKey(displayPath)
		pkgFiles[importPath][k] = append(pkgFiles[importPath][k], UncoveredBranch{
			Line:           row.line,
			Kind:           row.kind,
			Recommendation: row.recommendation,
		})
	}

	pkgKeys := make([]string, 0, len(pkgFiles))
	for k := range pkgFiles {
		pkgKeys = append(pkgKeys, k)
	}
	sort.Strings(pkgKeys)

	out := UncoveredBranchesReport{Packages: make([]UncoveredPackage, 0, len(pkgKeys))}
	for _, importPath := range pkgKeys {
		filesMap := pkgFiles[importPath]
		fileKeys := make([]string, 0, len(filesMap))
		for fk := range filesMap {
			fileKeys = append(fileKeys, string(fk))
		}
		sort.Strings(fileKeys)

		ufiles := make([]UncoveredFile, 0, len(fileKeys))
		for _, p := range fileKeys {
			branches := filesMap[fileKey(p)]
			sort.Slice(branches, func(i, j int) bool {
				if branches[i].Line != branches[j].Line {
					return branches[i].Line < branches[j].Line
				}
				return branches[i].Kind < branches[j].Kind
			})
			ufiles = append(ufiles, UncoveredFile{
				Path:     p,
				Branches: branches,
			})
		}
		out.Packages = append(out.Packages, UncoveredPackage{
			ImportPath: importPath,
			Files:      ufiles,
		})
	}
	return out
}

func classifySourceFile(absFile, moduleRoot, modulePath string) (importPath, displayPath string) {
	absFile = filepath.Clean(absFile)
	moduleRoot = filepath.Clean(moduleRoot)

	if moduleRoot == "" || modulePath == "" {
		return fallbackPackageKey(absFile), filepath.ToSlash(absFile)
	}

	rel, err := filepath.Rel(moduleRoot, absFile)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fallbackPackageKey(absFile), filepath.ToSlash(absFile)
	}

	rel = filepath.Clean(rel)
	displayPath = filepath.ToSlash(rel)
	dir := filepath.Dir(rel)
	if dir == "." {
		importPath = modulePath
	} else {
		importPath = modulePath + "/" + filepath.ToSlash(dir)
	}
	return importPath, displayPath
}

func fallbackPackageKey(absFile string) string {
	return "_" + filepath.ToSlash(filepath.Dir(absFile))
}
