package config

import (
	"gopkg.in/ini.v1"
	"strings"
)

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
	copyValueToDefaultProfileIfAvailable(defaultProfileInConfig, selectedProfileInConfig, "region", "mfa_serial")
}

func SetSelectedAssumedProfileAsDefault(selectedAssumedProfileName string, configFile *ini.File) {
	selectedProfile := configFile.Section(selectedAssumedProfileName)
	selectedRoleArn := selectedProfile.Key("role_arn").Value()
	selectedSourceProfile := selectedProfile.Key("source_profile").Value()

	defaultProfile := configFile.Section("default")
	defaultProfile.Key("role_arn").SetValue(selectedRoleArn)
	defaultProfile.Key("source_profile").SetValue(selectedSourceProfile)

	copyValueToDefaultProfileIfAvailable(defaultProfile, selectedProfile, "region", "mfa_serial")
}

func findConfigSectionByName(name string, configFile *ini.File) *ini.Section {
	for _, section := range configFile.Sections() {
		if strings.EqualFold(strings.TrimPrefix(section.Name(), "profile "), name) {
			return section
		}
	}

	return nil
}

func copyValueToDefaultProfileIfAvailable(defaultProfile *ini.Section, selectedProfile *ini.Section, keys ...string) {
	for _, key := range keys {
		if selectedProfile != nil &&
			selectedProfile.HasKey(key) {
			defaultProfile.Key(key).SetValue(selectedProfile.Key(key).Value())
		} else {
			defaultProfile.DeleteKey(key)
		}
	}
}