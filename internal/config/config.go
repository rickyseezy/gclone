package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	appDirName = "gclone"
	configFile = "config.yaml"
)

type Config struct {
	Profiles map[string]Profile `yaml:"profiles"`
	Defaults Defaults           `yaml:"defaults"`
}

type Profile struct {
	SSHHostAlias string `yaml:"ssh_host_alias"`
}

type Defaults struct {
	Profile string `yaml:"profile"`
}

// ConfigPath resolves the OS-appropriate configuration path.
func ConfigPath() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve config dir: %w", err)
	}
	return filepath.Join(cfgDir, appDirName, configFile), nil
}

func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, &MissingConfigError{Path: path}
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if cfg.Profiles == nil {
		cfg.Profiles = map[string]Profile{}
	}
	return &cfg, nil
}

// SelectProfile returns the chosen profile name and profile definition.
func SelectProfile(cfg *Config, requested string) (string, Profile, error) {
	if cfg == nil {
		return "", Profile{}, errors.New("config is nil")
	}
	name := requested
	if name == "" {
		name = cfg.Defaults.Profile
	}
	if name == "" {
		return "", Profile{}, errors.New("no profile selected; set --profile or defaults.profile")
	}
	profile, ok := cfg.Profiles[name]
	if !ok {
		return "", Profile{}, fmt.Errorf("profile %q not found in config", name)
	}
	if profile.SSHHostAlias == "" {
		return "", Profile{}, fmt.Errorf("profile %q missing ssh_host_alias", name)
	}
	return name, profile, nil
}

type MissingConfigError struct {
	Path string
}

func (e *MissingConfigError) Error() string {
	return fmt.Sprintf("config file not found at %s", e.Path)
}

func (e *MissingConfigError) Example() string {
	return `profiles:\n  profile1:\n    ssh_host_alias: "gitlab.com-work"\n  personal:\n    ssh_host_alias: "github.com-personal"\ndefaults:\n  profile: "personal"\n`
}
