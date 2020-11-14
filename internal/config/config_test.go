package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path"
	"runtime"
	"strings"
	"testing"
)

func TestFromFile(t *testing.T) {
	t.Run("return config with default values if config file not exists", func(t *testing.T) {
		c, err := FromFile("not-exist-config.yaml")

		require.NoError(t, err)
		expectedConfig := &config{
			HighlightColor: DefaultHighlightColor(),
			Regions:        DefaultRegions(),
		}
		require.Equal(t, expectedConfig, c)
	})

	t.Run("return unmarshalled config with default highlight color if highlight color is missing", func(t *testing.T) {
		c, err := FromFile("testdata/missing-highlight-color.yaml")

		require.NoError(t, err)
		expectedConfig := &config{
			HighlightColor: DefaultHighlightColor(),
			Regions: []string{
				"ap-southeast-2",
				"us-west-2",
				"us-east-1",
			},
		}
		require.Equal(t, expectedConfig, c)
	})

	t.Run("return unmarshalled config with default regions if regions are missing", func(t *testing.T) {
		c, err := FromFile("testdata/missing-regions.yaml")

		require.NoError(t, err)
		expectedConfig := &config{
			HighlightColor: "yellow",
			Regions:        DefaultRegions(),
		}
		require.Equal(t, expectedConfig, c)
	})

	t.Run("return unmarshalled config for sample config file", func(t *testing.T) {
		// this test makes sure sample config file is always in sync with code changes
		c, err := FromFile(sampleConfigPath(t))

		require.NoError(t, err)
		expectedConfig := &config{
			HighlightColor: "red",
			Regions: []string{
				"ap-southeast-2",
				"us-west-2",
				"us-east-1",
			},
		}
		require.Equal(t, expectedConfig, c)
	})

	t.Run("return error if highlight color is not in the predefined list", func(t *testing.T) {
		_, err := FromFile("testdata/invalid-highlight-color.yaml")

		require.Error(t, err)
		require.Equal(t, "valid values for highlight color are: black, red, green, yellow, blue, magenta, cyan, white", err.Error())
	})

	t.Run("return error if failed to unmarshal config file", func (t *testing.T) {
		_, err := FromFile("testdata/invalid-config.yaml")

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to unmarshal config file")
	})

}

func sampleConfigPath(t *testing.T) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		assert.Fail(t, "failed to get project root path")
	}

	projectRoot := strings.ReplaceAll(path.Dir(filename), "/internal/config", "")
	return path.Join(projectRoot, "configs/sample-config.yaml")
}
