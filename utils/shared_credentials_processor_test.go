package utils

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

func TestGetProfilesFromCredentialsFile_ReturnNonDefaultProfiles(t *testing.T) {
	credentialsFile := ini.Empty()
	AddCredentialsSection(credentialsFile, "default")
	AddCredentialsSection(credentialsFile, "profile-1")

	processor := AWSSharedCredentialsProcessor{
		CredentialsFile: credentialsFile,
		ConfigFile:      nil,
	}

	result := processor.getProfilesFromCredentialsFile()

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "profile-1", result[0].ProfileName)
}

func TestGetAssumedProfilesFromConfigFile_ReturnNonDefaultProfilesWithRoleArnAndSourceProfileAttributes(t *testing.T) {
	configFile := ini.Empty()
	AddConfigSection(configFile, "default")
	profile1Section, _ := configFile.NewSection("profile-1")
	profile1Section.Key("some_attribute").SetValue("some-value")
	AddConfigSection(configFile, "profile-2")

	processor := AWSSharedCredentialsProcessor{
		CredentialsFile: nil,
		ConfigFile:      configFile,
	}

	result := processor.getAssumedProfilesFromConfigFile()

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "profile-2", result[0].ProfileName)
}

func TestGetAssumedProfilesFromConfigFile_ReturnRoleArnAndMFASerialWithProfilesIfAvailable(t *testing.T) {
	configFile := ini.Empty()
	AddConfigSection(configFile, "default")
	AddConfigSection(configFile, "profile-1")
	profile2Section := AddConfigSection(configFile, "profile-2")
	profile2Section.Key("mfa_serial").SetValue("12345")

	processor := AWSSharedCredentialsProcessor{
		CredentialsFile: nil,
		ConfigFile:      configFile,
	}

	result := processor.getAssumedProfilesFromConfigFile()

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "profile-1-role-arn", result[0].RoleArn)
	assert.Equal(t, "profile-2-role-arn", result[1].RoleArn)
	assert.Equal(t, "12345", result[1].MFASerialNumber)
}

func TestGetAssumedProfilesFromConfigFile_ReturnRegionWithProfilesIfAvailable(t *testing.T) {
	configFile := ini.Empty()
	AddConfigSection(configFile, "default")
	AddConfigSection(configFile, "profile-1")
	profile2Section := AddConfigSection(configFile, "profile-2")
	profile2Section.Key("region").SetValue("ap-southeast-2")

	processor := AWSSharedCredentialsProcessor{
		CredentialsFile: nil,
		ConfigFile:      configFile,
	}

	result := processor.getAssumedProfilesFromConfigFile()

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "profile-1-role-arn", result[0].RoleArn)
	assert.Equal(t, "profile-2-role-arn", result[1].RoleArn)
	assert.Equal(t, "ap-southeast-2", result[1].Region)
}

func TestFindRegionInConfigFile_ReturnEmptyIfNoProfileWithGivenNameFound(t *testing.T) {
	configFile := ini.Empty()
	AddConfigSection(configFile, "default")
	AddConfigSection(configFile, "profile-1")
	profile2Section := AddConfigSection(configFile, "profile-2")
	profile2Section.Key("region").SetValue("ap-southeast-2")

	processor := AWSSharedCredentialsProcessor{
		CredentialsFile: nil,
		ConfigFile:      configFile,
	}

	result := processor.findRegionInConfigFile("profile-3")

	assert.Empty(t, result)
}

func TestFindRegionInConfigFile_ReturnEmptyIfSelectedProfileHasNoRegion(t *testing.T) {
	configFile := ini.Empty()
	AddConfigSection(configFile, "default")
	AddConfigSection(configFile, "profile-1")
	profile2Section := AddConfigSection(configFile, "profile-2")
	profile2Section.Key("region").SetValue("ap-southeast-2")

	processor := AWSSharedCredentialsProcessor{
		CredentialsFile: nil,
		ConfigFile:      configFile,
	}

	result := processor.findRegionInConfigFile("profile-1")

	assert.Empty(t, result)
}

func TestFindRegionInConfigFile_ReturnRegionIfSelectedProfileFoundAndHasRegion(t *testing.T) {
	configFile := ini.Empty()
	AddConfigSection(configFile, "default")
	AddConfigSection(configFile, "profile-1")
	profile2Section := AddConfigSection(configFile, "profile-2")
	profile2Section.Key("region").SetValue("ap-southeast-2")

	processor := AWSSharedCredentialsProcessor{
		CredentialsFile: nil,
		ConfigFile:      configFile,
	}

	result := processor.findRegionInConfigFile("profile-2")

	assert.Equal(t, "ap-southeast-2", result)
}

func TestSetSelectedProfileAsDefault_SetDefaultProfileAccessKeyIdAndSecretAccessKey(t *testing.T) {
	credentialsFile := ini.Empty()
	AddCredentialsSection(credentialsFile, "default")
	AddCredentialsSection(credentialsFile, "profile-1")
	AddCredentialsSection(credentialsFile, "profile-2")

	configFile := ini.Empty()

	processor := AWSSharedCredentialsProcessor{
		CredentialsFile: credentialsFile,
		ConfigFile:      configFile,
	}

	processor.SetSelectedProfileAsDefault("profile-2")

	defaultSection := processor.CredentialsFile.Section("default")
	assert.Equal(t, "profile-2-id", defaultSection.Key("aws_access_key_id").Value())
	assert.Equal(t, "profile-2-secret", defaultSection.Key("aws_secret_access_key").Value())
}

func TestSetSelectedProfileAsDefault_ResetConfigFileDefaultProfile(t *testing.T) {
	credentialsFile := ini.Empty()
	AddCredentialsSection(credentialsFile, "default")
	AddCredentialsSection(credentialsFile, "profile-1")
	AddCredentialsSection(credentialsFile, "profile-2")

	configFile := ini.Empty()
	AddConfigSection(configFile, "default")
	AddConfigSection(configFile, "profile-1")

	processor := AWSSharedCredentialsProcessor{
		CredentialsFile: credentialsFile,
		ConfigFile:      configFile,
	}

	defaultSection := processor.ConfigFile.Section("default")
	assert.NotEmpty(t, defaultSection.Key("role_arn").Value())
	assert.NotEmpty(t, defaultSection.Key("source_profile").Value())

	processor.SetSelectedProfileAsDefault("profile-2")

	defaultSection = processor.ConfigFile.Section("default")
	assert.Empty(t, defaultSection.Key("role_arn").Value())
	assert.Empty(t, defaultSection.Key("source_profile").Value())
}

func TestSetSelectedProfileAsDefault_SetDefaultProfileRegionInConfigFileIfAvailable(t *testing.T) {
	credentialsFile := ini.Empty()
	AddCredentialsSection(credentialsFile, "default")
	AddCredentialsSection(credentialsFile, "profile-1")
	AddCredentialsSection(credentialsFile, "profile-2")

	configFile := ini.Empty()
	AddConfigSection(configFile, "default")
	profile1Section := AddConfigSection(configFile, "profile profile-1")
	profile1Section.Key("region").SetValue("us-east-2")

	processor := AWSSharedCredentialsProcessor{
		CredentialsFile: credentialsFile,
		ConfigFile:      configFile,
	}

	defaultSection := processor.ConfigFile.Section("default")
	assert.Empty(t, defaultSection.Key("region").Value())

	processor.SetSelectedProfileAsDefault("profile-1")

	defaultSection = processor.ConfigFile.Section("default")
	assert.Equal(t, "us-east-2", defaultSection.Key("region").Value())
}

func TestSetSelectedProfileAsDefault_ClearDefaultRegionIfSelectedProfileHasNoRegion(t *testing.T) {
	credentialsFile := ini.Empty()
	AddCredentialsSection(credentialsFile, "default")
	AddCredentialsSection(credentialsFile, "profile-1")
	AddCredentialsSection(credentialsFile, "profile-2")

	configFile := ini.Empty()
	defaultSection := AddConfigSection(configFile, "default")
	defaultSection.Key("region").SetValue("ap-southeast-2")
	AddConfigSection(configFile, "profile-3")

	processor := AWSSharedCredentialsProcessor{
		CredentialsFile: credentialsFile,
		ConfigFile:      configFile,
	}

	assert.Equal(t, "ap-southeast-2", processor.ConfigFile.Section("default").Key("region").Value())

	processor.SetSelectedProfileAsDefault("profile-1")

	assert.Empty(t, processor.ConfigFile.Section("default").Key("region").Value())
}

func TestSetSelectedAssumedProfileAsDefault_SetConfigFileDefaultProfileRoleArnAndSourceProfile(t *testing.T) {
	credentialsFile := ini.Empty()

	configFile := ini.Empty()
	AddConfigSection(configFile, "default")
	AddConfigSection(configFile, "profile-1")
	AddConfigSection(configFile, "profile-2")

	processor := AWSSharedCredentialsProcessor{
		CredentialsFile: credentialsFile,
		ConfigFile:      configFile,
	}

	processor.SetSelectedAssumedProfileAsDefault("profile-2")

	defaultSection := processor.ConfigFile.Section("default")
	assert.Equal(t, "profile-2-role-arn", defaultSection.Key("role_arn").Value())
	assert.Equal(t, "profile-2-source-profile", defaultSection.Key("source_profile").Value())
	assert.Empty(t, defaultSection.Key("region").Value())
}

func TestSetSelectedAssumedProfileAsDefault_SetConfigFileDefaultProfileRegionIfOneIsAvailable(t *testing.T) {
	credentialsFile := ini.Empty()

	configFile := ini.Empty()
	AddConfigSection(configFile, "default")
	AddConfigSection(configFile, "profile-1")
	profile2Section := AddConfigSection(configFile, "profile-2")
	profile2Section.Key("region").SetValue("us-west-2")

	processor := AWSSharedCredentialsProcessor{
		CredentialsFile: credentialsFile,
		ConfigFile:      configFile,
	}

	processor.SetSelectedAssumedProfileAsDefault("profile-2")

	defaultSection := processor.ConfigFile.Section("default")
	assert.Equal(t, "us-west-2", defaultSection.Key("region").Value())
}

func TestSetSelectedAssumedProfileAsDefault_ClearDefaultRegionIfSelectedAssumedProfileHasNoRegion(t *testing.T) {
	credentialsFile := ini.Empty()

	configFile := ini.Empty()
	defaultSection := AddConfigSection(configFile, "default")
	defaultSection.Key("region").SetValue("ap-southeast-2")
	AddConfigSection(configFile, "profile-1")
	AddConfigSection(configFile, "profile-2")

	processor := AWSSharedCredentialsProcessor{
		CredentialsFile: credentialsFile,
		ConfigFile:      configFile,
	}

	assert.Equal(t, "ap-southeast-2", processor.ConfigFile.Section("default").Key("region").Value())

	processor.SetSelectedAssumedProfileAsDefault("profile-1")

	assert.Empty(t, processor.ConfigFile.Section("default").Key("region").Value())
}
