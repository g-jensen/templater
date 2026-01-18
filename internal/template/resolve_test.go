package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveDependencies_NestedFeature(t *testing.T) {
	available := []string{"auth", "auth/oauth", "auth/oauth/google", "database"}
	got := ResolveDependencies("auth/oauth/google", available, false)
	assert.Equal(t, []string{"auth", "auth/oauth", "auth/oauth/google"}, got)
}

func TestResolveDependencies_TopLevelWithRoot(t *testing.T) {
	available := []string{"auth", "database"}
	got := ResolveDependencies("database", available, true)
	assert.Equal(t, []string{"", "database"}, got)
}

func TestResolveDependencies_TopLevelWithoutRoot(t *testing.T) {
	available := []string{"auth", "database"}
	got := ResolveDependencies("database", available, false)
	assert.Equal(t, []string{"database"}, got)
}

func TestResolveDependencies_SkipsMissingIntermediate(t *testing.T) {
	available := []string{"auth", "auth/oauth/google"}
	got := ResolveDependencies("auth/oauth/google", available, false)
	assert.Equal(t, []string{"auth", "auth/oauth/google"}, got)
}
