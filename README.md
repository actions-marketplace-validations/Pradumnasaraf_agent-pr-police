# Agent PR Police

**Agent PR Police** is a GitHub Action that flags pull requests opened by AI coding agents like GitHub Copilot, Claude Code, Devin, Cursor, and Codex. When it detects an agent PR, it adds a label and posts a single comment summarizing which agent opened it and what the PR touched, so reviewers know at a glance that a bot wrote the change.

It does not gate or block anything. It is a transparency layer, not a security scanner.

## Features

- **Agent detection**: Detects agent-authored PRs from the author login, a label, or `Co-authored-by` commit trailers. Add your own identifiers for agents that aren't built in.
- **Auto label**: Applies a `pr-by-ai` label (configurable) to detected agent PRs. If you add the label by hand, the PR is treated as agent-authored too.
- **Sticky summary comment**: Posts one comment with the agent name and a summary of what the PR changed (files, added and removed lines, top areas). It updates in place on every run, so no duplicates.
- **Opt in or out**: The label and the comment are each independently toggleable.

## Detected agents

GitHub Copilot, Claude Code, Devin, Cursor, OpenAI Codex, Aider, Google Jules, Sourcegraph Cody, and Sweep. Agents whose tooling commits under a human account (like Claude Code and Aider) are matched by their co-author trailer to avoid false positives on human names.

## Usage

All inputs are optional.

| Input | Default | Description |
| --- | --- | --- |
| `label` | `pr-by-ai` | Label applied to detected agent PRs. Also recognized as an agent marker if added manually. |
| `add-label` | `true` | Apply the label to detected agent PRs. |
| `comment` | `true` | Post and update a single sticky comment summarizing the PR. |
| `treat-all-prs-as-agent` | `false` | Skip detection and treat every PR as agent-authored. |
| `extra-agent-identifiers` | `` | Newline-separated substrings matched against the author login and co-author trailers, for agents not in the built-in registry. |
| `github-token` | `${{ github.token }}` | Token used to read the PR, add the label, and post the comment. |

Outputs: `is-agent-pr` (`true` or `false`) and `agent` (the detected agent name, empty if none).

### Event trigger

Agent PR Police runs on **pull request events** only.

```yaml
on:
  pull_request:
```

### Permissions

To add the label and post the comment, the job needs write access to pull requests:

```yaml
permissions:
  contents: read
  pull-requests: write
```

### Example workflow

> [!IMPORTANT]
> Before using the snippet below, check the latest version in the `uses` field from the [GitHub Marketplace](https://github.com/marketplace/actions/agent-pr-police).

```yaml
name: Agent PR Police

on:
  pull_request:

jobs:
  police:
    runs-on: ubuntu-latest

    permissions:
      contents: read
      pull-requests: write

    steps:
      - name: Running Agent PR Police
        uses: Pradumnasaraf/agent-pr-police@v1
```

## Contributing

If you have suggestions for improving Agent PR Police or want to report a bug, feel free to open an issue. All contributions are welcome. For more details, check out the [Contributing Guide](CONTRIBUTING.md).

## License

This project is licensed under the [Apache License 2.0](LICENSE).

## Security

For information on reporting security vulnerabilities, please refer to the [Security Policy](SECURITY.md).

[build-ci]: https://github.com/Pradumnasaraf/agent-pr-police/actions/workflows/ci.yml
[build-ci-badge]: https://github.com/Pradumnasaraf/agent-pr-police/actions/workflows/ci.yml/badge.svg
[release]: https://github.com/Pradumnasaraf/agent-pr-police/releases
[release-badge]: https://img.shields.io/github/v/release/Pradumnasaraf/agent-pr-police
[actions-marketplace]: https://github.com/marketplace/actions/agent-pr-police
[actions-marketplace-badge]: https://img.shields.io/badge/marketplace-Agent%20PR%20Police-blue?&logo=github
