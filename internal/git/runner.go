package git

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

type CommandSpec struct {
	Name string
	Args []string
}

func (c CommandSpec) String() string {
	return strings.Join(append([]string{c.Name}, c.Args...), " ")
}

func BuildCloneCommand(repoURL, dest string) CommandSpec {
	args := []string{"clone", repoURL}
	if dest != "" {
		args = append(args, dest)
	}
	return CommandSpec{Name: "git", Args: args}
}

type Runner struct {
	Exec func(name string, args ...string) *exec.Cmd
}

func NewRunner() Runner {
	return Runner{Exec: exec.Command}
}

func (r Runner) Run(spec CommandSpec) (int, error) {
	execFn := r.Exec
	if execFn == nil {
		execFn = exec.Command
	}
	cmd := execFn(spec.Name, spec.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err == nil {
		return 0, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), err
	}
	return 1, err
}
