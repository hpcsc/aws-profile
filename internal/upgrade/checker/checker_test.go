package checker

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewChecker(t *testing.T) {
	t.Run("should return artifactory checker if include prerelease is true", func(t *testing.T) {
		c := NewChecker("linux", true)

		_, ok := c.(artifactoryChecker)

		require.True(t, ok)
	})

	t.Run("should return github checker if include prerelease is false", func(t *testing.T) {
		c := NewChecker("linux", false)

		_, ok := c.(githubChecker)

		require.True(t, ok)
	})
}
