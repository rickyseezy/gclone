package git

import (
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestBuildCloneCommand(t *testing.T) {
	spec := BuildCloneCommand("git@github.com-work:org/repo.git", "")
	if spec.Name != "git" {
		t.Fatalf("unexpected name: %s", spec.Name)
	}
	if len(spec.Args) != 2 || spec.Args[0] != "clone" || spec.Args[1] != "git@github.com-work:org/repo.git" {
		t.Fatalf("unexpected args: %#v", spec.Args)
	}

	spec = BuildCloneCommand("git@github.com-work:org/repo.git", "dest")
	if len(spec.Args) != 3 || spec.Args[2] != "dest" {
		t.Fatalf("unexpected args: %#v", spec.Args)
	}
}

func TestRunnerRun(t *testing.T) {
	r := Runner{Exec: fakeExecCommand}
	spec := CommandSpec{Name: "git", Args: []string{"clone", "repo"}}
	code, err := r.Run(spec)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if code != 0 {
		t.Fatalf("unexpected code: %d", code)
	}
}

func TestRunnerRunExitCode(t *testing.T) {
	r := Runner{Exec: fakeExecCommandWithExit(3)}
	spec := CommandSpec{Name: "git", Args: []string{"clone", "repo"}}
	code, err := r.Run(spec)
	if err == nil {
		t.Fatalf("expected err")
	}
	if code != 3 {
		t.Fatalf("unexpected code: %d", code)
	}
}

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cmd := exec.Command(os.Args[0], append([]string{"-test.run=TestHelperProcess", "--"}, args...)...)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1", "HELPER_EXIT_CODE=0")
	return cmd
}

func fakeExecCommandWithExit(code int) func(string, ...string) *exec.Cmd {
	return func(command string, args ...string) *exec.Cmd {
		cmd := exec.Command(os.Args[0], append([]string{"-test.run=TestHelperProcess", "--"}, args...)...)
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1", "HELPER_EXIT_CODE="+strconv.Itoa(code))
		return cmd
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	codeStr := os.Getenv("HELPER_EXIT_CODE")
	if codeStr == "" || codeStr == "0" {
		os.Exit(0)
	}
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(code)
}
