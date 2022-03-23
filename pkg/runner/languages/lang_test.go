package languages

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseLanguage(t *testing.T) {
	expected := LanguageGo
	actual, err := ParseLanguage("Golang")
	require.Nil(t, err)
	require.Equal(t, expected, actual)
	_, err = ParseLanguage("notalang")
	require.NotNil(t, err)
}

func TestBuildCommandLine(t *testing.T) {
	expectedCmd := "echo"
	expectedArgs := []string{"hello", "world"}
	actualCmd, actualArgs := buildCommandLine("echo hello world ")
	require.Equal(t, expectedCmd, actualCmd)
	require.Equal(t, expectedArgs, actualArgs)
}
