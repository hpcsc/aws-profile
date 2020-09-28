package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
	"strings"
)

type AWSSharedCredentialsProcessor struct {
	CredentialsFile *ini.File
	ConfigFile      *ini.File
}

func (processor AWSSharedCredentialsProcessor) getProfilesFromCredentialsFile() []AWSProfile {
	var profiles []AWSProfile

	for _, section := range processor.CredentialsFile.Sections() {
		if !strings.EqualFold(section.Name(), "default") {
			profiles = append(profiles, AWSProfile{
				ProfileName:        section.Name(),
				DisplayProfileName: section.Name(),
			})
		}
	}

	return profiles
}

func (processor AWSSharedCredentialsProcessor) getAssumedProfilesFromConfigFile() []AWSProfile {
	var profiles []AWSProfile

	for _, section := range processor.ConfigFile.Sections() {
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

func (processor AWSSharedCredentialsProcessor) GetProfilesFromCredentialsAndConfig() AWSProfiles {
	return AWSProfiles{
		CredentialsProfiles:   processor.getProfilesFromCredentialsFile(),
		ConfigAssumedProfiles: processor.getAssumedProfilesFromConfigFile(),
	}
}

func (processor AWSSharedCredentialsProcessor) findRegionInConfigFile(selectedProfileName string) string {
	for _, section := range processor.ConfigFile.Sections() {
		if strings.EqualFold(strings.TrimPrefix(section.Name(), "profile "), selectedProfileName) &&
			section.HasKey("region") {
			return section.Key("region").Value()
		}
	}

	return ""
}

func (processor AWSSharedCredentialsProcessor) SetSelectedProfileAsDefault(selectedProfileName string) {
	credentialsFile := processor.CredentialsFile

	selectedProfile := credentialsFile.Section(selectedProfileName)
	selectedKeyId := selectedProfile.Key("aws_access_key_id").Value()
	selectedAccessKey := selectedProfile.Key("aws_secret_access_key").Value()

	defaultProfileInCredentials := credentialsFile.Section("default")
	defaultProfileInCredentials.Key("aws_access_key_id").SetValue(selectedKeyId)
	defaultProfileInCredentials.Key("aws_secret_access_key").SetValue(selectedAccessKey)

	defaultProfileInConfig := processor.ConfigFile.Section("default")
	defaultProfileInConfig.DeleteKey("role_arn")
	defaultProfileInConfig.DeleteKey("source_profile")

	selectedRegion := processor.findRegionInConfigFile(selectedProfileName)
	if selectedRegion != "" {
		defaultProfileInConfig.Key("region").SetValue(selectedRegion)
	} else {
		defaultProfileInConfig.DeleteKey("region")
	}
}

func (processor AWSSharedCredentialsProcessor) SetSelectedAssumedProfileAsDefault(selectedAssumedProfileName string) {
	configFile := processor.ConfigFile

	selectedProfile := configFile.Section(selectedAssumedProfileName)
	selectedRoleArn := selectedProfile.Key("role_arn").Value()
	selectedSourceProfile := selectedProfile.Key("source_profile").Value()

	defaultProfile := configFile.Section("default")
	defaultProfile.Key("role_arn").SetValue(selectedRoleArn)
	defaultProfile.Key("source_profile").SetValue(selectedSourceProfile)

	if selectedProfile.HasKey("region") {
		defaultProfile.Key("region").SetValue(selectedProfile.Key("region").Value())
	} else {
		defaultProfile.DeleteKey("region")
	}

	if selectedProfile.HasKey("mfa_serial") {
		defaultProfile.Key("mfa_serial").SetValue(selectedProfile.Key("mfa_serial").Value())
	} else {
		defaultProfile.DeleteKey("mfa_serial")
	}
}
