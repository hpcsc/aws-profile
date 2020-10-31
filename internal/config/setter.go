package config

import (
	"gopkg.in/ini.v1"
	"strings"
)

func findConfigSectionByName(name string, configFile *ini.File) *ini.Section {
	for _, section := range configFile.Sections() {
		if strings.EqualFold(strings.TrimPrefix(section.Name(), "profile "), name) {
			return section
		}
	}

	return nil
}

func SetSelectedProfileAsDefault(selectedProfileName string, credentialsFile *ini.File, configFile *ini.File) {
	selectedProfileInCredentials := credentialsFile.Section(selectedProfileName)
	selectedKeyId := selectedProfileInCredentials.Key("aws_access_key_id").Value()
	selectedAccessKey := selectedProfileInCredentials.Key("aws_secret_access_key").Value()

	defaultProfileInCredentials := credentialsFile.Section("default")
	defaultProfileInCredentials.Key("aws_access_key_id").SetValue(selectedKeyId)
	defaultProfileInCredentials.Key("aws_secret_access_key").SetValue(selectedAccessKey)

	defaultProfileInConfig := configFile.Section("default")
	defaultProfileInConfig.DeleteKey("role_arn")
	defaultProfileInConfig.DeleteKey("source_profile")

	selectedProfileInConfig := findConfigSectionByName(selectedProfileName, configFile)

	if selectedProfileInConfig != nil &&
		selectedProfileInConfig.HasKey("region") {
		defaultProfileInConfig.Key("region").SetValue(selectedProfileInConfig.Key("region").Value())
	} else {
		defaultProfileInConfig.DeleteKey("region")
	}

	if selectedProfileInConfig != nil &&
		selectedProfileInConfig.HasKey("mfa_serial") {
		defaultProfileInConfig.Key("mfa_serial").SetValue(selectedProfileInConfig.Key("mfa_serial").Value())
	} else {
		defaultProfileInConfig.DeleteKey("mfa_serial")
	}
}

func SetSelectedAssumedProfileAsDefault(selectedAssumedProfileName string, configFile *ini.File) {
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
