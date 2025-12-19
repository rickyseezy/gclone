package sshconfig

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var ErrNoSSHConfig = errors.New("ssh config not found")

type Config struct {
	Aliases map[string]struct{}
}

func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".ssh", "config"), nil
}

func Load() (*Config, string, error) {
	path, err := DefaultPath()
	if err != nil {
		return nil, "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if runtime.GOOS == "windows" {
				return nil, path, ErrNoSSHConfig
			}
			return nil, path, ErrNoSSHConfig
		}
		return nil, path, fmt.Errorf("read ssh config: %w", err)
	}
	cfg := Parse(string(data))
	return cfg, path, nil
}

func Parse(content string) *Config {
	aliases := make(map[string]struct{})
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if strings.EqualFold(fields[0], "Host") {
			for _, host := range fields[1:] {
				aliases[host] = struct{}{}
			}
		}
	}
	return &Config{Aliases: aliases}
}

func (c *Config) AliasExists(alias string) bool {
	if c == nil {
		return false
	}
	_, ok := c.Aliases[alias]
	return ok
}

func MissingAliasError(profile, alias, configPath string) error {
	return fmt.Errorf("ssh config missing alias %q for profile %q (checked %s). Example:\nHost %s\n  HostName gitlab.com\n  User git", alias, profile, configPath, alias)
}
