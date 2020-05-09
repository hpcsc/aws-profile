package utils

import "strings"

type AWSProfile struct {
	ProfileName        string
	DisplayProfileName string
	RoleArn            string
	MFASerialNumber    string
	Region             string
}

type AWSProfiles struct {
	CredentialsProfiles   []AWSProfile
	ConfigAssumedProfiles []AWSProfile
}

func (profiles AWSProfiles) FindProfileInCredentialsFile(selected string) *AWSProfile {
	for _, profile := range profiles.CredentialsProfiles {
		if strings.EqualFold(profile.ProfileName, selected) {
			return &profile
		}
	}

	return nil
}

func (profiles AWSProfiles) FindProfileInConfigFile(selected string) *AWSProfile {
	for _, profile := range profiles.ConfigAssumedProfiles {
		if strings.EqualFold(profile.ProfileName, selected) {
			return &profile
		}
	}
	return nil
}

func (profiles AWSProfiles) GetAllDisplayProfileNames() []string {
	var displayProfileNames []string

	for _, profile := range profiles.CredentialsProfiles {
		displayProfileNames = append(displayProfileNames, profile.DisplayProfileName)
	}

	for _, profile := range profiles.ConfigAssumedProfiles {
		displayProfileNames = append(displayProfileNames, profile.DisplayProfileName)
	}

	return displayProfileNames
}

func (profiles AWSProfiles) Filter(pattern string) []AWSProfile {
	var filteredProfiles []AWSProfile

	for _, profile := range profiles.CredentialsProfiles {
		if pattern == "" || strings.Contains(profile.ProfileName, pattern) {
			filteredProfiles = append(filteredProfiles, profile)
		}
	}

	for _, profile := range profiles.ConfigAssumedProfiles {
		if pattern == "" || strings.Contains(profile.ProfileName, pattern) {
			filteredProfiles = append(filteredProfiles, profile)
		}
	}

	return filteredProfiles
}
