package config

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/ini.v1"
	"testing"
)

func TestSetSelectedProfileAsDefault(t *testing.T) {
	t.Run("set default profile access key id and secret access key", func(t *testing.T) {
		credentialsFile := ini.Empty()
		AddCredentialsSection(credentialsFile, "default")
		AddCredentialsSection(credentialsFile, "profile-1")
		AddCredentialsSection(credentialsFile, "profile-2")

		configFile := ini.Empty()

		SetSelectedProfileAsDefault("profile-2", credentialsFile, configFile)

		defaultSection := credentialsFile.Section("default")
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

		defaultSection := configFile.Section("default")
		assert.NotEmpty(t, defaultSection.Key("role_arn").Value())
		assert.NotEmpty(t, defaultSection.Key("source_profile").Value())

		SetSelectedProfileAsDefault("profile-2", credentialsFile, configFile)

		defaultSection = configFile.Section("default")
		assert.Empty(t, defaultSection.Key("role_arn").Value())
		assert.Empty(t, defaultSection.Key("source_profile").Value())
	})

	t.Run("set default profile session token if available", func(t *testing.T) {
		credentialsFile := ini.Empty()
		AddCredentialsSection(credentialsFile, "default")
		AddCredentialsSection(credentialsFile, "profile-1")
		profile2Section := AddCredentialsSection(credentialsFile, "profile-2")
		profile2Section.Key("aws_session_token").SetValue("profile-2-session-token")

		configFile := ini.Empty()

		SetSelectedProfileAsDefault("profile-2", credentialsFile, configFile)

		defaultSection := credentialsFile.Section("default")
		assert.Equal(t, "profile-2-id", defaultSection.Key("aws_access_key_id").Value())
		assert.Equal(t, "profile-2-secret", defaultSection.Key("aws_secret_access_key").Value())
		assert.Equal(t, "profile-2-session-token", defaultSection.Key("aws_session_token").Value())
	})

	t.Run("clear default profile session token if session token not set", func(t *testing.T) {
		credentialsFile := ini.Empty()
		defaultSection := AddCredentialsSection(credentialsFile, "default")
		defaultSection.Key("aws_session_token").SetValue("default-session-token")
		AddCredentialsSection(credentialsFile, "profile-1")
		AddCredentialsSection(credentialsFile, "profile-2")

		configFile := ini.Empty()

		SetSelectedProfileAsDefault("profile-2", credentialsFile, configFile)

		assert.Equal(t, "profile-2-id", defaultSection.Key("aws_access_key_id").Value())
		assert.Equal(t, "profile-2-secret", defaultSection.Key("aws_secret_access_key").Value())
		assert.False(t, defaultSection.HasKey("aws_session_token"))
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

		defaultSection := configFile.Section("default")
		assert.Empty(t, defaultSection.Key("region").Value())

		SetSelectedProfileAsDefault("profile-1", credentialsFile, configFile)

		defaultSection = configFile.Section("default")
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

		assert.Equal(t, "ap-southeast-2", configFile.Section("default").Key("region").Value())

		SetSelectedProfileAsDefault("profile-1", credentialsFile, configFile)

		assert.Empty(t, configFile.Section("default").Key("region").Value())
	})

	t.Run("set default profile mfa serial in config file if available", func(t *testing.T) {
		credentialsFile := ini.Empty()
		AddCredentialsSection(credentialsFile, "default")
		AddCredentialsSection(credentialsFile, "profile-1")
		AddCredentialsSection(credentialsFile, "profile-2")

		configFile := ini.Empty()
		AddConfigSection(configFile, "default")
		profile1Section := AddConfigSection(configFile, "profile profile-1")
		profile1Section.Key("mfa_serial").SetValue("my-mfa-serial")

		defaultSection := configFile.Section("default")
		assert.Empty(t, defaultSection.Key("mfa_serial").Value())

		SetSelectedProfileAsDefault("profile-1", credentialsFile, configFile)

		defaultSection = configFile.Section("default")
		assert.Equal(t, "my-mfa-serial", defaultSection.Key("mfa_serial").Value())
	})

	t.Run("clear default mfa serial if selected profile has no mfa serial", func(t *testing.T) {
		credentialsFile := ini.Empty()
		AddCredentialsSection(credentialsFile, "default")
		AddCredentialsSection(credentialsFile, "profile-1")
		AddCredentialsSection(credentialsFile, "profile-2")

		configFile := ini.Empty()
		defaultSection := AddConfigSection(configFile, "default")
		defaultSection.Key("mfa_serial").SetValue("my-mfa-serial")
		AddConfigSection(configFile, "profile-3")

		assert.Equal(t, "my-mfa-serial", configFile.Section("default").Key("mfa_serial").Value())

		SetSelectedProfileAsDefault("profile-1", credentialsFile, configFile)

		assert.Empty(t, configFile.Section("default").Key("mfa_serial").Value())
	})
}

func TestSetSelectedAssumedProfileAsDefault(t *testing.T) {
	t.Run("set role arn and source profile for default profile in config file", func(t *testing.T) {
		configFile := ini.Empty()
		AddConfigSection(configFile, "default")
		AddConfigSection(configFile, "profile-1")
		AddConfigSection(configFile, "profile-2")

		SetSelectedAssumedProfileAsDefault("profile-2", configFile)

		defaultSection := configFile.Section("default")
		assert.Equal(t, "profile-2-role-arn", defaultSection.Key("role_arn").Value())
		assert.Equal(t, "profile-2-source-profile", defaultSection.Key("source_profile").Value())
		assert.Empty(t, defaultSection.Key("region").Value())
	})

	t.Run("set region for default profile in config file", func(t *testing.T) {
		configFile := ini.Empty()
		AddConfigSection(configFile, "default")
		AddConfigSection(configFile, "profile-1")
		profile2Section := AddConfigSection(configFile, "profile-2")
		profile2Section.Key("region").SetValue("us-west-2")

		SetSelectedAssumedProfileAsDefault("profile-2", configFile)

		defaultSection := configFile.Section("default")
		assert.Equal(t, "us-west-2", defaultSection.Key("region").Value())
	})

	t.Run("clear default region if selected assumed profile has no region", func(t *testing.T) {
		configFile := ini.Empty()
		defaultSection := AddConfigSection(configFile, "default")
		defaultSection.Key("region").SetValue("ap-southeast-2")
		AddConfigSection(configFile, "profile-1")
		AddConfigSection(configFile, "profile-2")

		assert.Equal(t, "ap-southeast-2", configFile.Section("default").Key("region").Value())

		SetSelectedAssumedProfileAsDefault("profile-1", configFile)

		assert.Empty(t, configFile.Section("default").Key("region").Value())
	})

	t.Run("set mfa_serial for default profile in config file", func(t *testing.T) {
		configFile := ini.Empty()
		AddConfigSection(configFile, "default")
		AddConfigSection(configFile, "profile-1")
		profile2Section := AddConfigSection(configFile, "profile-2")
		profile2Section.Key("mfa_serial").SetValue("profile-2-mfa-serial")

		SetSelectedAssumedProfileAsDefault("profile-2", configFile)

		defaultSection := configFile.Section("default")
		assert.Equal(t, "profile-2-mfa-serial", defaultSection.Key("mfa_serial").Value())
	})

	t.Run("clear default mfa_serial if selected assumed profile has no mfa_serial", func(t *testing.T) {
		configFile := ini.Empty()
		defaultSection := AddConfigSection(configFile, "default")
		defaultSection.Key("mfa_serial").SetValue("initial-mfa-serial")
		AddConfigSection(configFile, "profile-1")
		AddConfigSection(configFile, "profile-2")

		assert.Equal(t, "initial-mfa-serial", configFile.Section("default").Key("mfa_serial").Value())

		SetSelectedAssumedProfileAsDefault("profile-1", configFile)

		assert.Empty(t, configFile.Section("default").Key("mfa_serial").Value())
	})
}
