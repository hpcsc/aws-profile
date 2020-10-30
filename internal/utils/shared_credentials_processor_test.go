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

func TestGetProfilesFromCredentialsFile(t *testing.T) {
	t.Run("return non default profiles", func(t *testing.T) {
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

	})
}

func TestGetAssumedProfilesFromConfigFile(t *testing.T) {
	t.Run("return non default profiles with role arn and source profile attributes", func(t *testing.T) {
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
		assert.Equal(t, "profile-2-role-arn", result[0].RoleArn)
		assert.Equal(t, "profile-2-source-profile", result[0].SourceProfile)
	})

	t.Run("return role arn and mfa serial with profiles if available", func(t *testing.T) {
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

	})

	t.Run("return region with profiles if available", func(t *testing.T) {
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

	})

}

func TestFindRegionInConfigFile(t *testing.T) {
	t.Run("return empty if profile not found", func(t *testing.T) {
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

	})

	t.Run("return empty if selected profile has no region", func(t *testing.T) {
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

	})

	t.Run("return region if selected profile has region", func(t *testing.T) {
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

	})

}

func TestSetSelectedProfileAsDefault(t *testing.T) {
	t.Run("set default profile access key id and secret access key", func(t *testing.T) {
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

	})

	t.Run("reset default profile in config file", func(t *testing.T) {
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

	})

	t.Run("set default profile region in config file if available", func(t *testing.T) {
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
	})

	t.Run("clear default region if selected profile has no region", func(t *testing.T) {
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

	})
}

func TestSetSelectedAssumedProfileAsDefault(t *testing.T) {
	t.Run("set role arn and source profile for default profile in config file", func(t *testing.T) {
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

	})

	t.Run("set region for default profile in config file", func(t *testing.T) {
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

	})

	t.Run("clear default region if selected assumed profile has no region", func(t *testing.T) {
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

	})

	t.Run("set mfa_serial for default profile in config file", func(t *testing.T) {
		credentialsFile := ini.Empty()

		configFile := ini.Empty()
		AddConfigSection(configFile, "default")
		AddConfigSection(configFile, "profile-1")
		profile2Section := AddConfigSection(configFile, "profile-2")
		profile2Section.Key("mfa_serial").SetValue("profile-2-mfa-serial")

		processor := AWSSharedCredentialsProcessor{
			CredentialsFile: credentialsFile,
			ConfigFile:      configFile,
		}

		processor.SetSelectedAssumedProfileAsDefault("profile-2")

		defaultSection := processor.ConfigFile.Section("default")
		assert.Equal(t, "profile-2-mfa-serial", defaultSection.Key("mfa_serial").Value())

	})

	t.Run("clear default mfa_serial if selected assumed profile has no mfa_serial", func(t *testing.T) {
		credentialsFile := ini.Empty()

		configFile := ini.Empty()
		defaultSection := AddConfigSection(configFile, "default")
		defaultSection.Key("mfa_serial").SetValue("initial-mfa-serial")
		AddConfigSection(configFile, "profile-1")
		AddConfigSection(configFile, "profile-2")

		processor := AWSSharedCredentialsProcessor{
			CredentialsFile: credentialsFile,
			ConfigFile:      configFile,
		}

		assert.Equal(t, "initial-mfa-serial", processor.ConfigFile.Section("default").Key("mfa_serial").Value())

		processor.SetSelectedAssumedProfileAsDefault("profile-1")

		assert.Empty(t, processor.ConfigFile.Section("default").Key("mfa_serial").Value())

	})
}
