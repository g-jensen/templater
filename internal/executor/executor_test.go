package executor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShellExecutor_CapturesStdout(t *testing.T) {
	executor := NewShellExecutor()

	stdout, _, _, err := executor.Execute("echo hello", "10s", nil)

	require.NoError(t, err)
	assert.Equal(t, "hello\n", stdout)
}

func TestShellExecutor_CapturesStderr(t *testing.T) {
	executor := NewShellExecutor()

	_, stderr, _, err := executor.Execute("echo error >&2", "10s", nil)

	require.NoError(t, err)
	assert.Equal(t, "error\n", stderr)
}

func TestShellExecutor_ReturnsExitCode(t *testing.T) {
	executor := NewShellExecutor()

	_, _, exitCode, err := executor.Execute("exit 42", "10s", nil)

	assert.NoError(t, err)
	assert.Equal(t, 42, exitCode)
}

func TestShellExecutor_PassesEnvVars(t *testing.T) {
	executor := NewShellExecutor()
	env := map[string]string{"MY_VAR": "hello"}

	stdout, _, _, err := executor.Execute("echo $MY_VAR", "10s", env)

	require.NoError(t, err)
	assert.Equal(t, "hello\n", stdout)
}

func TestShellExecutor_ReturnsErrTimeoutWhenCommandExceedsTimeout(t *testing.T) {
	executor := NewShellExecutor()

	_, _, _, err := executor.Execute("sleep 10", "100ms", nil)

	assert.True(t, errors.Is(err, ErrTimeout))
}

func TestExecuteWithStdin_PassesStdinToCommand(t *testing.T) {
	exec := NewShellExecutor()

	stdout, stderr, exitCode, err := exec.ExecuteWithStdin("cat", "5s", nil, "hello from stdin")

	require.NoError(t, err)
	assert.Equal(t, 0, exitCode)
	assert.Equal(t, "hello from stdin", stdout)
	assert.Empty(t, stderr)
}

func TestExecuteWithStdin_CommandCanProcessStdin(t *testing.T) {
	exec := NewShellExecutor()

	stdout, stderr, exitCode, err := exec.ExecuteWithStdin("wc -c", "5s", nil, "12345")

	require.NoError(t, err)
	assert.Equal(t, 0, exitCode)
	assert.Contains(t, stdout, "5")
	assert.Empty(t, stderr)
}
