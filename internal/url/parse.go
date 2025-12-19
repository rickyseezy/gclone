package url

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var scpLikeRe = regexp.MustCompile(`^([^@]+)@([^:]+):(.+)$`)

type RepoURL struct {
	Scheme string
	User   string
	Host   string
	Path   string
}

func Parse(input string) (RepoURL, error) {
	if strings.HasPrefix(input, "ssh://") || strings.HasPrefix(input, "https://") {
		parsed, err := url.Parse(input)
		if err != nil {
			return RepoURL{}, fmt.Errorf("parse url: %w", err)
		}
		if parsed.Host == "" {
			return RepoURL{}, fmt.Errorf("invalid url: missing host")
		}
		path := strings.TrimPrefix(parsed.Path, "/")
		if path == "" {
			return RepoURL{}, fmt.Errorf("invalid url: missing path")
		}
		if parsed.Scheme == "ssh" {
			user := ""
			if parsed.User != nil {
				user = parsed.User.Username()
			}
			return RepoURL{Scheme: "ssh", User: user, Host: parsed.Hostname(), Path: path}, nil
		}
		if parsed.Scheme == "https" {
			return RepoURL{Scheme: "https", Host: parsed.Hostname(), Path: path}, nil
		}
		return RepoURL{}, fmt.Errorf("unsupported scheme: %s", parsed.Scheme)
	}

	matches := scpLikeRe.FindStringSubmatch(input)
	if len(matches) != 4 {
		return RepoURL{}, fmt.Errorf("unsupported repo url format")
	}
	return RepoURL{Scheme: "ssh", User: matches[1], Host: matches[2], Path: matches[3]}, nil
}
