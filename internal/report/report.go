package report

import (
	"fmt"
	"strings"

	"github.com/pradumnasaraf/agent-pr-police/internal/summary"
)

const Marker = "<!-- agent-pr-police:sticky-comment -->"

type Input struct {
	Agent   string
	Signals []string
	Summary summary.Summary
}

func Build(in Input) string {
	var b strings.Builder
	b.WriteString(Marker)
	b.WriteString("\n## 🤖 Agent PR Police\n\n")

	who := "an AI coding agent"
	if in.Agent != "" {
		who = "the **" + in.Agent + "** coding agent"
	}
	fmt.Fprintf(&b, "This PR was opened by %s.\n\n", who)

	s := in.Summary
	fmt.Fprintf(&b, "**Changes:** %s, +%d / -%d lines.\n\n", filesLabel(s.Files), s.Additions, s.Deletions)

	if bd := breakdown(s); bd != "" {
		fmt.Fprintf(&b, "**Breakdown:** %s.\n\n", bd)
	}

	if len(s.TopDirs) > 0 {
		b.WriteString("**Top areas**\n\n")
		b.WriteString("| Area | Files |\n| --- | --- |\n")
		for _, d := range s.TopDirs {
			fmt.Fprintf(&b, "| `%s` | %d |\n", d.Dir, d.Files)
		}
		b.WriteString("\n")
	}

	if len(in.Signals) > 0 {
		b.WriteString("<details><summary>Why this was flagged as an agent PR</summary>\n\n")
		for _, sig := range in.Signals {
			b.WriteString("- " + sig + "\n")
		}
		b.WriteString("\n</details>\n")
	}

	return b.String()
}

func filesLabel(n int) string {
	if n == 1 {
		return "1 file"
	}
	return fmt.Sprintf("%d files", n)
}

func breakdown(s summary.Summary) string {
	var parts []string
	if s.Added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", s.Added))
	}
	if s.Modified > 0 {
		parts = append(parts, fmt.Sprintf("%d modified", s.Modified))
	}
	if s.Removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", s.Removed))
	}
	if s.Renamed > 0 {
		parts = append(parts, fmt.Sprintf("%d renamed", s.Renamed))
	}
	return strings.Join(parts, ", ")
}
