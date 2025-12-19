package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"

	"gclone/internal/config"
	"gclone/internal/git"
	"gclone/internal/sshconfig"
	"gclone/internal/url"
)

type Options struct {
	RepoURL string
	Profile string
	Dest    string
	DryRun  bool
	Verbose bool
	Init    bool
}

type App struct {
	ConfigLoader    func() (*config.Config, error)
	SSHConfigLoader func() (*sshconfig.Config, string, error)
	Runner          git.Runner
	Out             io.Writer
	Err             io.Writer
}

func New() *App {
	return &App{
		ConfigLoader:    config.Load,
		SSHConfigLoader: sshconfig.Load,
		Runner:          git.NewRunner(),
		Out:             os.Stdout,
		Err:             os.Stderr,
	}
}

func (a *App) Run(opts Options) (int, error) {
	if opts.Init {
		path, err := config.Init()
		if err != nil {
			return 1, err
		}
		fmt.Fprintf(a.Out, "Created config at %s\n", path)
		return 0, nil
	}

	cfg, err := a.ConfigLoader()
	if err != nil {
		var missing *config.MissingConfigError
		if errors.As(err, &missing) {
			return 1, fmt.Errorf("%s\n\nExample config:\n%s", missing.Error(), missing.Example())
		}
		return 1, err
	}

	profileName, profile, err := config.SelectProfile(cfg, opts.Profile)
	if err != nil {
		return 1, err
	}

	repoURL, err := url.Parse(opts.RepoURL)
	if err != nil {
		return 1, err
	}
	finalURL, err := url.Rewrite(repoURL, profile.SSHHostAlias)
	if err != nil {
		return 1, err
	}

	sshCfg, sshPath, err := a.SSHConfigLoader()
	if err != nil {
		if errors.Is(err, sshconfig.ErrNoSSHConfig) {
			if runtime.GOOS == "windows" {
				if opts.Verbose {
					fmt.Fprintf(a.Err, "Warning: ssh config not found at %s; skipping alias validation on Windows.\n", sshPath)
				}
			} else {
				return 1, fmt.Errorf("ssh config not found at %s", sshPath)
			}
		} else {
			return 1, err
		}
	} else {
		if !sshCfg.AliasExists(profile.SSHHostAlias) {
			return 1, sshconfig.MissingAliasError(profileName, profile.SSHHostAlias, sshPath)
		}
	}

	if opts.Verbose {
		fmt.Fprintf(a.Err, "Profile: %s (alias %s)\n", profileName, profile.SSHHostAlias)
		fmt.Fprintf(a.Err, "Repo URL: %s\n", opts.RepoURL)
		fmt.Fprintf(a.Err, "Rewritten URL: %s\n", finalURL)
	}

	spec := git.BuildCloneCommand(finalURL, opts.Dest)
	if opts.DryRun {
		fmt.Fprintf(a.Out, "Rewritten URL: %s\n", finalURL)
		fmt.Fprintf(a.Out, "Command: %s\n", spec.String())
		return 0, nil
	}

	code, err := a.Runner.Run(spec)
	if err != nil {
		return code, err
	}
	return 0, nil
}
