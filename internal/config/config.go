package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"strings"
)

func LoadProfilesFromConfigAndCredentials(credentialsFile *ini.File, configFile *ini.File) AWSProfiles {
	return AWSProfiles{
		CredentialsProfiles:   getProfilesFromCredentialsFile(credentialsFile),
		ConfigAssumedProfiles: getAssumedProfilesFromConfigFile(configFile),
	}
}

func getProfilesFromCredentialsFile(credentialsFile *ini.File) []AWSProfile {
	var profiles []AWSProfile

	if credentialsFile == nil {
		return profiles
	}

	for _, section := range credentialsFile.Sections() {
		if !strings.EqualFold(section.Name(), "default") {
			profiles = append(profiles, AWSProfile{
				ProfileName:        section.Name(),
				DisplayProfileName: section.Name(),
			})
		}
	}

	return profiles
}

func getAssumedProfilesFromConfigFile(configFile *ini.File) []AWSProfile {
	var profiles []AWSProfile

	if configFile == nil {
		return profiles
	}

	for _, section := range configFile.Sections() {
		if !strings.EqualFold(section.Name(), "default") &&
			section.HasKey("role_arn") &&
			section.HasKey("source_profile") {
			profile := AWSProfile{
				ProfileName:        section.Name(),
				DisplayProfileName: fmt.Sprintf("assume %s", section.Name()),
				RoleArn:            section.Key("role_arn").Value(),
				SourceProfile:      section.Key("source_profile").Value(),
			}

			if section.HasKey("mfa_serial") {
				profile.MFASerialNumber = section.Key("mfa_serial").Value()
			}

			if section.HasKey("region") {
				profile.Region = section.Key("region").Value()
			}

			profiles = append(profiles, profile)
		}
	}

	return profiles
}
