package utils

import "strings"

type AWSProfiles struct {
	CredentialsProfiles []string
	ConfigAssumedProfiles []string
}

func (profiles AWSProfiles) CredentialsFileContains(selected string) bool {
	for _, profile := range profiles.CredentialsProfiles {
		if strings.EqualFold(profile, selected) {
			return true
		}
	}
	return false
}

func (profiles AWSProfiles) ConfigFileContains(selected string) bool {
	for _, profile := range profiles.ConfigAssumedProfiles {
		if strings.EqualFold(profile, selected) {
			return true
		}
	}
	return false
}
