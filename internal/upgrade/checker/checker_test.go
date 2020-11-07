package checker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewChecker(t *testing.T) {
	t.Run("should return bintray checker if include prerelease is true", func(t *testing.T) {
		c := NewChecker("linux", true)

		_, ok := c.(bintrayChecker)

		assert.True(t, ok)
	})

	t.Run("should return github checker if include prerelease is false", func(t *testing.T) {
		c := NewChecker("linux", false)

		_, ok := c.(githubChecker)

		assert.True(t, ok)
	})
}
