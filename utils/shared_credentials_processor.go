package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
	"strings"
)

type AWSSharedCredentialsProcessor struct {
	CredentialsFile *ini.File
	ConfigFile *ini.File
}

func (processor AWSSharedCredentialsProcessor) getProfilesFromCredentialsFile() []string {
	var profiles []string

	for _, section := range processor.CredentialsFile.Sections() {
		if !strings.EqualFold(section.Name(), "default") {
			profiles = append(profiles, section.Name())
		}
	}

	return profiles
}

func (processor AWSSharedCredentialsProcessor) getAssumedProfilesFromConfigFile() []string {
	var profiles []string

	for _, section := range processor.ConfigFile.Sections() {
		if !strings.EqualFold(section.Name(), "default") &&
			section.HasKey("role_arn") &&
			section.HasKey("source_profile") {
			profiles = append(profiles, fmt.Sprintf("assume %s:%s", section.Name(), section.Name()))
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

func (processor AWSSharedCredentialsProcessor) SetSelectedProfileAsDefault(selectedProfile string) {
	credentialsFile := processor.CredentialsFile
	configFile := processor.ConfigFile

	selectedKeyId := credentialsFile.Section(selectedProfile).Key("aws_access_key_id").Value()
	selectedAccessKey := credentialsFile.Section(selectedProfile).Key("aws_secret_access_key").Value()

	credentialsFile.Section("default").Key("aws_access_key_id").SetValue(selectedKeyId)
	credentialsFile.Section("default").Key("aws_secret_access_key").SetValue(selectedAccessKey)
	configFile.Section("default").DeleteKey("role_arn")
	configFile.Section("default").DeleteKey("source_profile")
}

func (processor AWSSharedCredentialsProcessor) SetSelectedAssumedProfileAsDefault(selectedAssumedProfile string) {
	configFile := processor.ConfigFile


	selectedRoleArn := configFile.Section(selectedAssumedProfile).Key("role_arn").Value()
	selectedSourceProfile := configFile.Section(selectedAssumedProfile).Key("source_profile").Value()

	configFile.Section("default").Key("role_arn").SetValue(selectedRoleArn)
	configFile.Section("default").Key("source_profile").SetValue(selectedSourceProfile)
}

func (processor AWSSharedCredentialsProcessor) GetProfileSectionNameFromColonDelimited(selected string) string {
	if ! strings.Contains(selected, ":") {
		return selected
	}

	elements := strings.Split(selected, ":")
	return elements[1]
}
