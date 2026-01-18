package executor

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"
)

var ErrTimeout = errors.New("command timed out")

type Executor interface {
	Execute(command string, timeout string, env map[string]string) (stdout, stderr string, exitCode int, err error)
	ExecuteWithStdin(command string, timeout string, env map[string]string, stdin string) (stdout, stderr string, exitCode int, err error)
}

type ShellExecutor struct{}

func NewShellExecutor() *ShellExecutor {
	return &ShellExecutor{}
}

func (e *ShellExecutor) Execute(command string, timeout string, env map[string]string) (string, string, int, error) {
	return e.ExecuteWithStdin(command, timeout, env, "")
}

func (e *ShellExecutor) ExecuteWithStdin(command string, timeout string, env map[string]string, stdin string) (string, string, int, error) {
	duration := parseDuration(timeout)
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	cmd := buildCommand(ctx, command, env)
	var stdout, stderr bytes.Buffer
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return buildResult(ctx, err, stdout.String(), stderr.String())
}

func parseDuration(timeout string) time.Duration {
	duration, err := time.ParseDuration(timeout)
	if err != nil {
		return time.Hour
	}
	return duration
}

func buildCommand(ctx context.Context, command string, env map[string]string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Env = os.Environ()
	for key, value := range env {
		cmd.Env = append(cmd.Env, key+"="+value)
	}
	return cmd
}

func buildResult(ctx context.Context, err error, stdout string, stderr string) (string, string, int, error) {
	if err == nil {
		return stdout, stderr, 0, nil
	}
	if ctx.Err() == context.DeadlineExceeded {
		return stdout, stderr, -1, ErrTimeout
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return stdout, stderr, exitErr.ExitCode(), nil
	}
	return stdout, stderr, -1, err
}
