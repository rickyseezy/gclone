package url

import "fmt"

func Rewrite(u RepoURL, alias string) (string, error) {
	if alias == "" {
		return "", fmt.Errorf("ssh host alias is empty")
	}
	if u.Path == "" {
		return "", fmt.Errorf("repo path is empty")
	}

	switch u.Scheme {
	case "ssh":
		user := u.User
		if user == "" {
			user = "git"
		}
		return fmt.Sprintf("%s@%s:%s", user, alias, u.Path), nil
	case "https":
		return fmt.Sprintf("git@%s:%s", alias, u.Path), nil
	default:
		return "", fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}
}
