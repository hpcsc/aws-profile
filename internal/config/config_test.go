package config

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/ini.v1"
	"testing"
)

func AddCredentialsSection(file *ini.File, sectionName string) *ini.Section {
	section, _ := file.NewSection(sectionName)
	section.Key("aws_access_key_id").SetValue(sectionName + "-id")
	section.Key("aws_secret_access_key").SetValue(sectionName + "-secret")
	return section
}

func AddConfigSection(file *ini.File, sectionName string) *ini.Section {
	section, _ := file.NewSection(sectionName)
	section.Key("role_arn").SetValue(sectionName + "-role-arn")
	section.Key("source_profile").SetValue(sectionName + "-source-profile")
	return section
}

func TestLoadProfilesFromConfigAndCredentials(t *testing.T) {
	t.Run("return non default profiles from credentials file", func(t *testing.T) {
		credentialsFile := ini.Empty()
		AddCredentialsSection(credentialsFile, "default")
		AddCredentialsSection(credentialsFile, "profile-1")

		result := LoadProfilesFromConfigAndCredentials(credentialsFile, nil)

		assert.Equal(t, 0, len(result.ConfigAssumedProfiles))

		credentialsProfiles := result.CredentialsProfiles
		assert.Equal(t, 1, len(credentialsProfiles))
		assert.Equal(t, "profile-1", credentialsProfiles[0].ProfileName)
	})

	t.Run("return non default profiles from config file with role arn and source profile attributes", func(t *testing.T) {
		configFile := ini.Empty()
		AddConfigSection(configFile, "default")
		profile1Section, _ := configFile.NewSection("profile-1")
		profile1Section.Key("some_attribute").SetValue("some-value")
		AddConfigSection(configFile, "profile-2")

		result := LoadProfilesFromConfigAndCredentials(nil, configFile)

		assert.Equal(t, 0, len(result.CredentialsProfiles))

		configProfiles := result.ConfigAssumedProfiles
		assert.Equal(t, 1, len(configProfiles))
		assert.Equal(t, "profile-2", configProfiles[0].ProfileName)
		assert.Equal(t, "profile-2-role-arn", configProfiles[0].RoleArn)
		assert.Equal(t, "profile-2-source-profile", configProfiles[0].SourceProfile)
	})

	t.Run("return role arn and mfa serial with config profiles if available", func(t *testing.T) {
		configFile := ini.Empty()
		AddConfigSection(configFile, "default")
		AddConfigSection(configFile, "profile-1")
		profile2Section := AddConfigSection(configFile, "profile-2")
		profile2Section.Key("mfa_serial").SetValue("12345")

		result := LoadProfilesFromConfigAndCredentials(nil, configFile)
		configProfiles := result.ConfigAssumedProfiles

		assert.Equal(t, 2, len(configProfiles))
		assert.Equal(t, "profile-1-role-arn", configProfiles[0].RoleArn)
		assert.Equal(t, "profile-2-role-arn", configProfiles[1].RoleArn)
		assert.Equal(t, "12345", configProfiles[1].MFASerialNumber)
	})

	t.Run("return region with config profiles if available", func(t *testing.T) {
		configFile := ini.Empty()
		AddConfigSection(configFile, "default")
		AddConfigSection(configFile, "profile-1")
		profile2Section := AddConfigSection(configFile, "profile-2")
		profile2Section.Key("region").SetValue("ap-southeast-2")

		result := LoadProfilesFromConfigAndCredentials(nil, configFile)
		configProfiles := result.ConfigAssumedProfiles

		assert.Equal(t, 2, len(configProfiles))
		assert.Equal(t, "profile-1-role-arn", configProfiles[0].RoleArn)
		assert.Equal(t, "profile-2-role-arn", configProfiles[1].RoleArn)
		assert.Equal(t, "ap-southeast-2", configProfiles[1].Region)
	})
}
