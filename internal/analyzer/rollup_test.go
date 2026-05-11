package analyzer

import (
	"path/filepath"
	"testing"
)

func TestClassifySourceFileUnderModule(t *testing.T) {
	root := filepath.Join("testdata", "modroot")
	mod := "github.com/example/lib"
	abs := filepath.Join(root, "internal", "x", "a.go")

	gotImport, gotPath := classifySourceFile(abs, root, mod)
	wantImport := "github.com/example/lib/internal/x"
	wantPath := "internal/x/a.go"
	if gotImport != wantImport {
		t.Fatalf("importPath: got %q want %q", gotImport, wantImport)
	}
	if gotPath != wantPath {
		t.Fatalf("displayPath: got %q want %q", gotPath, wantPath)
	}
}

func TestGroupUncoveredSourcesSorts(t *testing.T) {
	rows := []uncoveredSource{
		{file: "/p/a.go", line: 2, kind: "b"},
		{file: "/p/a.go", line: 2, kind: "a"},
		{file: "/p/b.go", line: 1, kind: "x"},
	}
	got := groupUncoveredSources(rows, "", "")
	if len(got.Packages) != 1 {
		t.Fatalf("packages: %d", len(got.Packages))
	}
	p := got.Packages[0]
	if len(p.Files) != 2 {
		t.Fatalf("files: %d", len(p.Files))
	}
	if p.Files[0].Path != "/p/a.go" || p.Files[1].Path != "/p/b.go" {
		t.Fatalf("file order: %#v", p.Files)
	}
	br := p.Files[0].Branches
	if len(br) != 2 || br[0].Kind != "a" || br[1].Kind != "b" {
		t.Fatalf("branch sort: %#v", br)
	}
}
