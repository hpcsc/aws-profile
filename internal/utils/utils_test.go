package utils

import (
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestExpandHomeDirectory(t *testing.T) {
	t.Run("return original path if it does not start with tilde", func(t *testing.T) {
		output := ExpandHomeDirectory("/home/root/some-file")

		require.Equal(t, "/home/root/some-file", output)
	})

	t.Run("replace tilde with user home directory when path starts with tilde", func(t *testing.T) {
		output := ExpandHomeDirectory("~/.aws/config")

		require.False(t, strings.HasPrefix(output, "~/"))
	})
}

func TestGetEnvVariableOrDefault(t *testing.T) {
	t.Run("return environment variable value if it is set", func(t *testing.T) {
		os.Setenv("MY_VARIABLE", "my-variable")
		defer os.Unsetenv("MY_VARIABLE")

		value := GetEnvVariableOrDefault("MY_VARIABLE", "default-value")

		require.Equal(t, "my-variable", value)
	})

	t.Run("return default value if environment variable is not set", func(t *testing.T) {
		value := GetEnvVariableOrDefault("MY_VARIABLE", "default-value")

		require.Equal(t, "default-value", value)
	})
}
