# gclone

Clone git repositories using SSH host aliases defined in your SSH config. Useful when you have multiple git accounts and aliases like `github.com-work` or `gitlab.com-personal`.

## Install

```bash
go build -o gclone ./cmd/gclone
```

Or install into your GOPATH/bin:

```bash
go install ./cmd/gclone
```

## Homebrew

Once releases are published, you can install via Homebrew:

```bash
brew tap rickyseezy/gclone
brew install gclone
```

## Config

Config path is resolved using Go's OS-specific config directory:

- macOS/Linux: `~/.config/gclone/config.yaml`
- Windows: `%APPDATA%\gclone\config.yaml`

Example config:

```yaml
profiles:
  profile1:
    ssh_host_alias: "gitlab.com-work"
  personal:
    ssh_host_alias: "github.com-personal"
defaults:
  profile: "personal"
```

## SSH config alias validation

Before cloning, `gclone` checks that the chosen alias exists as a `Host` entry in your SSH config.

Example SSH config snippet:

```ssh-config
Host gitlab.com-work
  HostName gitlab.com
  User git
```

## Usage

```bash
gclone <repo_url> --profile <name> [--dest <path>] [--dry-run] [--verbose]
```

Examples:

```bash
gclone git@gitlab.com:stream-flow/backend/streamflow-api.git --profile profile1
gclone https://gitlab.com/stream-flow/backend/streamflow-api.git --profile profile1
```

## URL behavior

Supported inputs:

- SCP-like SSH: `git@github.com:org/repo.git`
- SSH URL: `ssh://git@github.com/org/repo.git`
- HTTPS: `https://github.com/org/repo.git`

Rewriting rules:

- SSH inputs: host is replaced with the selected profile's `ssh_host_alias`.
- HTTPS inputs: converted to SSH SCP-like format with user `git`.
- `.git` is preserved if present and never forced.

## Troubleshooting

- Missing config: ensure `config.yaml` exists in the OS config directory.
- Missing SSH alias: add a `Host <alias>` entry in `~/.ssh/config`.
- Use `--dry-run` to see the rewritten URL and command without cloning.
- Use `--verbose` to see profile and rewrite details.
