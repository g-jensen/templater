package executor

import (
	"templater/internal/executor"
)

type ExecutedCommand struct {
	Command string
	Timeout string
	Env     map[string]string
}

type FakeExecutor struct {
	Commands         []ExecutedCommand
	Stdout           string
	Stderr           string
	DefaultExitCode  int
	ExitCodes        map[string]int
	TimeoutCommands  map[string]bool
	TimeoutExitCodes map[string]int
	StdinReceived    string
}

func (fake *FakeExecutor) Execute(command string, timeout string, env map[string]string) (stdout, stderr string, exitCode int, err error) {
	fake.Commands = append(fake.Commands, ExecutedCommand{Command: command, Timeout: timeout, Env: env})
	if fake.shouldTimeout(command) {
		return "", "", fake.timeoutExitCode(command), executor.ErrTimeout
	}
	return fake.Stdout, fake.Stderr, fake.exitCodeFor(command), nil
}

func (fake *FakeExecutor) shouldTimeout(command string) bool {
	return fake.TimeoutCommands != nil && fake.TimeoutCommands[command]
}

func (fake *FakeExecutor) timeoutExitCode(command string) int {
	if fake.TimeoutExitCodes == nil {
		return -1
	}
	if code, ok := fake.TimeoutExitCodes[command]; ok {
		return code
	}
	return -1
}

func (fake *FakeExecutor) exitCodeFor(command string) int {
	if fake.ExitCodes == nil {
		return fake.DefaultExitCode
	}
	if code, ok := fake.ExitCodes[command]; ok {
		return code
	}
	return fake.DefaultExitCode
}

func (fake *FakeExecutor) ExecuteWithStdin(command string, timeout string, env map[string]string, stdin string) (stdout, stderr string, exitCode int, err error) {
	fake.StdinReceived = stdin
	return fake.Execute(command, timeout, env)
}
