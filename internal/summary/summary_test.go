package summary

import "testing"

func TestBuildCountsAndLines(t *testing.T) {
	s := Build([]ChangedFile{
		{Path: "cmd/main.go", Status: "added", Additions: 20, Deletions: 0},
		{Path: "internal/a/a.go", Status: "modified", Additions: 5, Deletions: 3},
		{Path: "internal/b/b.go", Status: "removed", Additions: 0, Deletions: 12},
		{Path: "go.mod", Status: "renamed", Additions: 1, Deletions: 1},
	})
	if s.Files != 4 {
		t.Errorf("Files = %d, want 4", s.Files)
	}
	if s.Added != 1 || s.Modified != 1 || s.Removed != 1 || s.Renamed != 1 {
		t.Errorf("status breakdown wrong: %+v", s)
	}
	if s.Additions != 26 || s.Deletions != 16 {
		t.Errorf("line totals wrong: +%d/-%d", s.Additions, s.Deletions)
	}
}

func TestBuildTopDirsSorted(t *testing.T) {
	s := Build([]ChangedFile{
		{Path: "internal/a.go"},
		{Path: "internal/b.go"},
		{Path: "cmd/main.go"},
		{Path: "README.md"},
	})
	if len(s.TopDirs) != 3 {
		t.Fatalf("TopDirs len = %d, want 3", len(s.TopDirs))
	}
	if s.TopDirs[0].Dir != "internal/" || s.TopDirs[0].Files != 2 {
		t.Errorf("top dir = %+v, want internal/ with 2 files", s.TopDirs[0])
	}
}

func TestTopDir(t *testing.T) {
	tests := map[string]string{
		"README.md":            "(root)",
		"internal/detect/x.go": "internal/",
		"./cmd/main.go":        "cmd/",
		"a/b/c.txt":            "a/",
	}
	for in, want := range tests {
		if got := topDir(in); got != want {
			t.Errorf("topDir(%q) = %q, want %q", in, got, want)
		}
	}
}
