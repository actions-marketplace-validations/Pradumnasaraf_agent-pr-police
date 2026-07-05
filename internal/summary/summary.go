package summary

import (
	"path"
	"sort"
	"strings"
)

type ChangedFile struct {
	Path      string
	Status    string
	Additions int
	Deletions int
}

type DirCount struct {
	Dir   string
	Files int
}

type Summary struct {
	Files     int
	Added     int
	Modified  int
	Removed   int
	Renamed   int
	Additions int
	Deletions int
	TopDirs   []DirCount
}

const maxTopDirs = 5

func Build(files []ChangedFile) Summary {
	var s Summary
	dirs := map[string]int{}
	for _, f := range files {
		s.Files++
		s.Additions += f.Additions
		s.Deletions += f.Deletions
		switch f.Status {
		case "added":
			s.Added++
		case "removed":
			s.Removed++
		case "renamed":
			s.Renamed++
		default:
			s.Modified++
		}
		dirs[topDir(f.Path)]++
	}
	for d, n := range dirs {
		s.TopDirs = append(s.TopDirs, DirCount{Dir: d, Files: n})
	}
	sort.Slice(s.TopDirs, func(i, j int) bool {
		if s.TopDirs[i].Files != s.TopDirs[j].Files {
			return s.TopDirs[i].Files > s.TopDirs[j].Files
		}
		return s.TopDirs[i].Dir < s.TopDirs[j].Dir
	})
	if len(s.TopDirs) > maxTopDirs {
		s.TopDirs = s.TopDirs[:maxTopDirs]
	}
	return s
}

func topDir(p string) string {
	p = strings.TrimPrefix(path.Clean(p), "./")
	if dir, _, found := strings.Cut(p, "/"); found {
		return dir + "/"
	}
	return "(root)"
}
