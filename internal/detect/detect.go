package detect

import (
	"regexp"
	"strings"
)

const DefaultLabel = "pr-by-ai"

type PullRequest struct {
	AuthorLogin    string
	AuthorType     string
	Labels         []string
	CommitMessages []string
}

type Config struct {
	Label            string
	TreatAllAsAgent  bool
	ExtraIdentifiers []string
}

type Result struct {
	IsAgent bool
	Agent   string
	Signals []string
}

type Agent struct {
	Name            string
	LoginPatterns   []string
	TrailerPatterns []string
}

var Registry = []Agent{
	{
		Name:            "GitHub Copilot",
		LoginPatterns:   []string{"copilot"},
		TrailerPatterns: []string{"copilot"},
	},
	{
		Name:            "Claude Code",
		LoginPatterns:   []string{"claude-bot", "anthropic"},
		TrailerPatterns: []string{"claude", "anthropic"},
	},
	{
		Name:            "Devin",
		LoginPatterns:   []string{"devin-ai-integration", "devin-ai"},
		TrailerPatterns: []string{"devin"},
	},
	{
		Name:            "Cursor",
		LoginPatterns:   []string{"cursoragent", "cursor-com", "cursor[bot]"},
		TrailerPatterns: []string{"cursor"},
	},
	{
		Name:            "OpenAI Codex",
		LoginPatterns:   []string{"chatgpt-codex", "codex-connector", "openai-codex"},
		TrailerPatterns: []string{"codex"},
	},
	{
		Name:            "Aider",
		LoginPatterns:   nil,
		TrailerPatterns: []string{"aider"},
	},
	{
		Name:            "Google Jules",
		LoginPatterns:   []string{"google-labs-jules"},
		TrailerPatterns: []string{"jules"},
	},
	{
		Name:            "Sourcegraph Cody",
		LoginPatterns:   []string{"sourcegraph-cody", "sourcegraph-bot"},
		TrailerPatterns: []string{"sourcegraph cody"},
	},
	{
		Name:            "Sweep",
		LoginPatterns:   []string{"sweep-ai"},
		TrailerPatterns: []string{"sweep-ai"},
	},
}

var coAuthorRe = regexp.MustCompile(`(?im)^\s*co-authored-by:\s*(.+)$`)

func Detect(pr PullRequest, cfg Config) Result {
	var res Result

	if cfg.TreatAllAsAgent {
		res.IsAgent = true
		res.Signals = append(res.Signals, "treat-all-prs-as-agent is enabled")
	}

	if name, ok := matchLogin(pr.AuthorLogin, cfg.ExtraIdentifiers); ok {
		res.IsAgent = true
		setAgent(&res, name)
		res.Signals = append(res.Signals, "PR author login matches "+name+": "+pr.AuthorLogin)
	}

	label := cfg.Label
	if label == "" {
		label = DefaultLabel
	}
	for _, l := range pr.Labels {
		if strings.EqualFold(strings.TrimSpace(l), label) {
			res.IsAgent = true
			res.Signals = append(res.Signals, "PR carries the agent label: "+label)
			break
		}
	}

	if name, trailer, ok := matchTrailers(pr.CommitMessages, cfg.ExtraIdentifiers); ok {
		res.IsAgent = true
		setAgent(&res, name)
		res.Signals = append(res.Signals, "a commit has an agent co-author trailer ("+name+"): "+trailer)
	}

	return res
}

func setAgent(res *Result, name string) {
	if res.Agent == "" && name != "" && name != "custom agent" {
		res.Agent = name
	}
}

func matchLogin(login string, extra []string) (string, bool) {
	l := strings.ToLower(strings.TrimSpace(login))
	if l == "" {
		return "", false
	}
	for _, a := range Registry {
		for _, p := range a.LoginPatterns {
			if strings.Contains(l, strings.ToLower(p)) {
				return a.Name, true
			}
		}
	}
	for _, e := range extra {
		if e = strings.ToLower(strings.TrimSpace(e)); e != "" && strings.Contains(l, e) {
			return "custom agent", true
		}
	}
	return "", false
}

func matchTrailers(messages, extra []string) (name, trailer string, ok bool) {
	for _, m := range messages {
		for _, match := range coAuthorRe.FindAllStringSubmatch(m, -1) {
			value := strings.TrimSpace(match[1])
			lower := strings.ToLower(value)
			for _, a := range Registry {
				for _, p := range a.TrailerPatterns {
					if strings.Contains(lower, strings.ToLower(p)) {
						return a.Name, value, true
					}
				}
			}
			for _, e := range extra {
				if e = strings.ToLower(strings.TrimSpace(e)); e != "" && strings.Contains(lower, e) {
					return "custom agent", value, true
				}
			}
		}
	}
	return "", "", false
}

func IdentifyLogin(login string) string {
	name, _ := matchLogin(login, nil)
	return name
}
