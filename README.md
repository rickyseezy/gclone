# gclone 🚀

**gclone** is a CLI that lets developers clone Git repositories using the correct SSH identity in one command. If you juggle multiple Git accounts (work/personal/clients), gclone eliminates mistakes and saves time by selecting the right SSH host alias automatically. ⚡️

## Why gclone ✅

- **Fast, reliable cloning** with the right SSH profile every time 🚄
- **Works with GitHub, GitLab, Bitbucket** (any host) 🌍
- **Zero guesswork**: pick a profile and go 🎯
- **Clean, predictable URLs** (HTTPS → SSH conversion supported) 🔁
- **Great for teams** with multiple orgs and identities 👥

## Install 🛠️

### Homebrew (recommended) 🍺

```bash
brew tap rickyseezy/gclone
brew install gclone
```

> Note: gclone is currently distributed via the official tap. Homebrew/core requires a notability threshold; we’ll submit once the repo meets it.
> If Homebrew reports outdated Command Line Tools, use the one‑line installer below as a fallback.

### One‑line install (no Homebrew / no CLT required) ⚡

```bash
curl -fsSL https://raw.githubusercontent.com/rickyseezy/gclone/main/install.sh | sh
```

If you want a custom location:

```bash
GCLONE_INSTALL_DIR=~/.local/bin curl -fsSL https://raw.githubusercontent.com/rickyseezy/gclone/main/install.sh | sh
```

### Go (direct) 🧰

```bash
go install github.com/rickyseezy/gclone/cmd/gclone@latest
```

### From source 🧱

```bash
go build -o gclone ./cmd/gclone
```

## Configuration ⚙️

gclone reads a YAML config from your OS‑specific config directory:

- macOS/Linux: `~/.config/gclone/config.yaml`
- Windows: `%APPDATA%\gclone\config.yaml`

Example:

```yaml
profiles:
  work:
    ssh_host_alias: "gitlab.com-work"
  personal:
    ssh_host_alias: "github.com-personal"
defaults:
  profile: "personal"
```

### SSH config 🔐

Ensure your SSH config defines the host aliases you reference:

```ssh-config
Host gitlab.com-work
  HostName gitlab.com
  User git

Host github.com-personal
  HostName github.com
  User git
```

## Usage ✨

```bash
gclone <repo_url> --profile <name> [--dest <path>] [--dry-run] [--verbose]
```

Examples (GitHub shown; works with GitHub and GitLab):

```bash
gclone git@github.com:octo-org/octo-repo.git --profile work

gclone https://github.com/octo-org/octo-repo.git --profile work
```

## How it works 🧠

gclone supports these input formats:

- SSH (scp‑like): `git@github.com:org/repo.git`
- SSH (url): `ssh://git@github.com/org/repo.git`
- HTTPS: `https://github.com/org/repo.git`

Rules:

- SSH inputs are rewritten to the selected profile’s **ssh_host_alias**.
- HTTPS inputs are converted to SSH with **user `git`**.
- The `.git` suffix is preserved if present.

## Troubleshooting 🧯

- **Missing config:** create `config.yaml` in the OS config directory.
- **Missing SSH alias:** add a `Host <alias>` entry in `~/.ssh/config`.
- **Check behavior:** use `--dry-run` to preview the rewritten URL and command.
- **Debug:** use `--verbose` to see selection details.

## Commercial use 💼

gclone is designed for production teams and multi‑account workflows. It’s lightweight, scriptable, and ideal for company‑wide developer onboarding. If you want priority support, custom features, or team onboarding help, open an issue or reach out via GitHub.

## License

MIT
