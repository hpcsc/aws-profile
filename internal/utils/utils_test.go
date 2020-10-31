package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetEnvVariableOrDefault(t *testing.T) {
	t.Run("return environment variable value if it is set", func (t *testing.T) {
		os.Setenv("MY_VARIABLE", "my-variable")
		defer os.Unsetenv("MY_VARIABLE")

	    value := GetEnvVariableOrDefault("MY_VARIABLE", "default-value")

	    assert.Equal(t, "my-variable", value)
	})

	t.Run("return default value if environment variable is not set", func (t *testing.T) {
		value := GetEnvVariableOrDefault("MY_VARIABLE", "default-value")

		assert.Equal(t, "default-value", value)
	})
}
