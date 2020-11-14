package config

import (
	"fmt"
	"github.com/hpcsc/aws-profile/internal/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	HighlightColor string   `yaml:"highlightColor"`
	Regions        []string `yaml:"regions"`
}

const defaultHighlightColor = "green"

var defaultRegions = []string{
	"af-south-1",
	"ap-east-1",
	"ap-northeast-1",
	"ap-northeast-2",
	"ap-northeast-3",
	"ap-south-1",
	"ap-southeast-1",
	"ap-southeast-2",
	"ca-central-1",
	"cn-north-1",
	"cn-northwest-1",
	"eu-central-1",
	"eu-north-1",
	"eu-south-1",
	"eu-west-1",
	"eu-west-2",
	"eu-west-3",
	"me-south-1",
	"sa-east-1",
	"us-east-1",
	"us-east-2",
	"us-gov-east-1",
	"us-gov-west-1",
	"us-west-1",
	"us-west-2",
}

var allowedColors = []string{
	"black",
	"red",
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
	"white",
}

func Load() (*Config, error) {
	configPath := utils.GetEnvVariableOrDefault("AWS_PROFILE_CONFIG", "~/.aws-profile/config.yaml")
	return FromFile(utils.ExpandHomeDirectory(configPath))
}

func FromFile(path string) (*Config, error) {
	cleanedPath := filepath.Clean(path)
	if !fileExists(cleanedPath) {
		return defaultConfig(), nil
	}

	fileContent, err := ioutil.ReadFile(cleanedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file at %s: %v", path, err)
	}

	c := &Config{}
	if err = yaml.Unmarshal(fileContent, c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %s: %v", path, err)
	}

	if c.HighlightColor == "" {
		c.HighlightColor = defaultHighlightColor
	}

	if c.Regions == nil || len(c.Regions) == 0 {
		c.Regions = defaultRegions
	}

	if !isValidHighlightColor(c.HighlightColor) {
		return nil, fmt.Errorf("valid values for highlight color are: %s", strings.Join(allowedColors, ", "))
	}

	return c, nil
}

func DefaultHighlightColor() string {
	return defaultHighlightColor
}

func DefaultRegions() []string {
	return defaultRegions
}

func isValidHighlightColor(color string) bool {
	for _, allowedColor := range allowedColors {
		if strings.EqualFold(color, allowedColor) {
			return true
		}
	}

	return false
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func defaultConfig() *Config {
	return &Config{
		HighlightColor: defaultHighlightColor,
		Regions:        defaultRegions,
	}
}
