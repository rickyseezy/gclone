package url

import "testing"

func TestRewriteSSH(t *testing.T) {
	u := RepoURL{Scheme: "ssh", User: "git", Host: "github.com", Path: "org/repo.git"}
	out, err := Rewrite(u, "github.com-work")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if out != "git@github.com-work:org/repo.git" {
		t.Fatalf("unexpected out: %s", out)
	}
}

func TestRewriteSSHDefaultUser(t *testing.T) {
	u := RepoURL{Scheme: "ssh", Host: "github.com", Path: "org/repo.git"}
	out, err := Rewrite(u, "github.com-work")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if out != "git@github.com-work:org/repo.git" {
		t.Fatalf("unexpected out: %s", out)
	}
}

func TestRewriteHTTPS(t *testing.T) {
	u := RepoURL{Scheme: "https", Host: "github.com", Path: "org/repo.git"}
	out, err := Rewrite(u, "github.com-work")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if out != "git@github.com-work:org/repo.git" {
		t.Fatalf("unexpected out: %s", out)
	}
}

func TestRewriteErrors(t *testing.T) {
	if _, err := Rewrite(RepoURL{Scheme: "ssh"}, ""); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := Rewrite(RepoURL{Scheme: "ftp"}, "alias"); err == nil {
		t.Fatalf("expected error")
	}
}
