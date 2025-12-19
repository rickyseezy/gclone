package sshconfig

import "testing"

func TestParseAndAliasExists(t *testing.T) {
	content := `# Comment
Host github.com-personal
  HostName github.com

Host gitlab.com-work gitlab.com-alt
  HostName gitlab.com
`
	cfg := Parse(content)
	if !cfg.AliasExists("github.com-personal") {
		t.Fatalf("expected alias")
	}
	if !cfg.AliasExists("gitlab.com-work") {
		t.Fatalf("expected alias")
	}
	if !cfg.AliasExists("gitlab.com-alt") {
		t.Fatalf("expected alias")
	}
	if cfg.AliasExists("missing") {
		t.Fatalf("did not expect alias")
	}
}
