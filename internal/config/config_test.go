package config

import "testing"

func TestSelectProfile(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{
			"work":     {SSHHostAlias: "gitlab.com-work"},
			"personal": {SSHHostAlias: "github.com-personal"},
		},
		Defaults: Defaults{Profile: "personal"},
	}

	name, p, err := SelectProfile(cfg, "work")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if name != "work" || p.SSHHostAlias != "gitlab.com-work" {
		t.Fatalf("unexpected selection: %q %+v", name, p)
	}

	name, p, err = SelectProfile(cfg, "")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if name != "personal" || p.SSHHostAlias != "github.com-personal" {
		t.Fatalf("unexpected selection: %q %+v", name, p)
	}
}

func TestSelectProfileErrors(t *testing.T) {
	cfg := &Config{Profiles: map[string]Profile{}}

	if _, _, err := SelectProfile(cfg, ""); err == nil {
		t.Fatalf("expected error for missing selection")
	}
	cfg.Defaults.Profile = "missing"
	if _, _, err := SelectProfile(cfg, ""); err == nil {
		t.Fatalf("expected error for missing profile")
	}

	cfg.Profiles["bad"] = Profile{}
	cfg.Defaults.Profile = "bad"
	if _, _, err := SelectProfile(cfg, ""); err == nil {
		t.Fatalf("expected error for missing ssh_host_alias")
	}
}
