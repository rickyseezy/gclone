package url

import "testing"

func TestParseSCP(t *testing.T) {
	u, err := Parse("git@github.com:org/repo.git")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if u.Scheme != "ssh" || u.User != "git" || u.Host != "github.com" || u.Path != "org/repo.git" {
		t.Fatalf("unexpected parsed: %+v", u)
	}
}

func TestParseSSHURL(t *testing.T) {
	u, err := Parse("ssh://git@github.com/org/repo.git")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if u.Scheme != "ssh" || u.User != "git" || u.Host != "github.com" || u.Path != "org/repo.git" {
		t.Fatalf("unexpected parsed: %+v", u)
	}
}

func TestParseHTTPS(t *testing.T) {
	u, err := Parse("https://github.com/org/repo.git")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if u.Scheme != "https" || u.Host != "github.com" || u.Path != "org/repo.git" {
		t.Fatalf("unexpected parsed: %+v", u)
	}
}

func TestParseInvalid(t *testing.T) {
	if _, err := Parse("git@github.com/org/repo.git"); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := Parse("ssh://github.com"); err == nil {
		t.Fatalf("expected error")
	}
}
