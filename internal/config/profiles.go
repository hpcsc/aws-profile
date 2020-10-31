package config

import "strings"

type Profiles struct {
	CredentialsProfiles   []Profile
	ConfigAssumedProfiles []Profile
}

func (profiles Profiles) FindProfileInCredentialsFile(selected string) *Profile {
	for _, profile := range profiles.CredentialsProfiles {
		if strings.EqualFold(profile.ProfileName, selected) {
			return &profile
		}
	}

	return nil
}

func (profiles Profiles) FindProfileInConfigFile(selected string) *Profile {
	for _, profile := range profiles.ConfigAssumedProfiles {
		if strings.EqualFold(profile.ProfileName, selected) {
			return &profile
		}
	}
	return nil
}

func (profiles Profiles) GetAllDisplayProfileNames() []string {
	var displayProfileNames []string

	for _, profile := range profiles.CredentialsProfiles {
		displayProfileNames = append(displayProfileNames, profile.DisplayProfileName)
	}

	for _, profile := range profiles.ConfigAssumedProfiles {
		displayProfileNames = append(displayProfileNames, profile.DisplayProfileName)
	}

	return displayProfileNames
}

func (profiles Profiles) Filter(pattern string) []Profile {
	var filteredProfiles []Profile

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
