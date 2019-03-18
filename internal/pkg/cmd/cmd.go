package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	cmd  string
	args []string
	env  []string
}

type Option func(*Command)

func New(cmd string, opts ...Option) *Command {
	c := &Command{
		cmd:  cmd,
		args: make([]string, 0),
		env:  make([]string, 0),
	}
	for _, option := range opts {
		option(c)
	}
	return c
}

func (cmd *Command) Execute(ctx context.Context) error {
	errors := make(chan error, 1)
	go func() { errors <- cmd.execute() }()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errors:
		return err
	}
}

func (cmd *Command) execute() error {
	if len(strings.TrimSpace(cmd.cmd)) == 0 {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Assemble command line
	cmdLine := []string{cmd.cmd}
	for _, arg := range cmd.args {
		cmdLine = append(cmdLine, arg)
	}

	// Wrap in shell
	shell := []string{"sh", "-exc"}
	shell = append(shell, strings.Join(cmdLine, " "))

	proc := exec.Command(shell[0], shell[1:]...)
	proc.Dir = cwd
	proc.Env = cmd.env
	proc.Stdin = os.Stdin
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	return proc.Run()
}

func WithArgs(args ...string) Option {
	return func(cmd *Command) {
		for _, arg := range args {
			cmd.args = append(cmd.args, arg)
		}
	}
}

func WithImplicitEnv() Option {
	return func(cmd *Command) {
		for _, pair := range os.Environ() {
			cmd.env = append(cmd.env, pair)
		}
	}
}

func WithEnv(pair string) Option {
	return func(cmd *Command) {
		cmd.env = append(cmd.env, pair)
	}
}

func WithEnvPair(key, value string) Option {
	return WithEnv(fmt.Sprintf("%s=%s", key, value))
}
