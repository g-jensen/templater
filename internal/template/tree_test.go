package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderTree_Empty(t *testing.T) {
	assert.Equal(t, "", RenderTree([]string{}))
}

func TestRenderTree_SingleFeature(t *testing.T) {
	assert.Equal(t, "└── auth\n", RenderTree([]string{"auth"}))
}

func TestRenderTree_MultipleTopLevel(t *testing.T) {
	features := []string{"api", "auth", "database"}
	expected := "├── api\n" +
		"├── auth\n" +
		"└── database\n"
	assert.Equal(t, expected, RenderTree(features))
}

func TestRenderTree_Nested(t *testing.T) {
	features := []string{
		"auth",
		"auth/oauth",
		"auth/oauth/github",
		"auth/oauth/google",
		"database",
		"database/migrations",
	}
	expected := "├── auth\n" +
		"│   └── oauth\n" +
		"│       ├── github\n" +
		"│       └── google\n" +
		"└── database\n" +
		"    └── migrations\n"
	assert.Equal(t, expected, RenderTree(features))
}
