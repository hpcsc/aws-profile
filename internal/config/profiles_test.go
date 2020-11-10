package config

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func StubProfiles(startIndex int, endIndex int) []Profile {
	var profiles []Profile

	for i := startIndex; i <= endIndex; i++ {
		profiles = append(profiles, Profile{
			ProfileName:        fmt.Sprintf("profile-%d", i),
			DisplayProfileName: fmt.Sprintf("profile-%d display name", i),
		})
	}

	return profiles
}

func TestFindProfileInCredentialsFile(t *testing.T) {
	t.Run("return nil if profile not found", func(t *testing.T) {
		profiles := Profiles{
			CredentialsProfiles:   StubProfiles(1, 2),
			ConfigAssumedProfiles: []Profile{},
		}

		result := profiles.FindProfileInCredentialsFile("profile-3")

		require.Nil(t, result)

	})

	t.Run("return profile if found", func(t *testing.T) {
		profiles := Profiles{
			CredentialsProfiles:   StubProfiles(1, 2),
			ConfigAssumedProfiles: []Profile{},
		}

		result := profiles.FindProfileInCredentialsFile("profile-2")

		require.NotNil(t, result)
		require.Equal(t, result.ProfileName, "profile-2")

	})
}

func TestFindProfileInConfigFile(t *testing.T) {
	t.Run("return nil if profile not found", func(t *testing.T) {
		profiles := Profiles{
			CredentialsProfiles:   []Profile{},
			ConfigAssumedProfiles: StubProfiles(1, 2),
		}

		result := profiles.FindProfileInConfigFile("profile-3")

		require.Nil(t, result)

	})

	t.Run("return profile if found", func(t *testing.T) {
		profiles := Profiles{
			CredentialsProfiles:   []Profile{},
			ConfigAssumedProfiles: StubProfiles(1, 2),
		}

		result := profiles.FindProfileInConfigFile("profile-2")

		require.NotNil(t, result)
		require.Equal(t, result.ProfileName, "profile-2")

	})
}

func TestGetAllDisplayProfileNames(t *testing.T) {
	t.Run("return profile names from both credentials and config files", func(t *testing.T) {
		profiles := Profiles{
			CredentialsProfiles:   StubProfiles(1, 2),
			ConfigAssumedProfiles: StubProfiles(3, 3),
		}

		result := profiles.GetAllDisplayProfileNames()

		expected := []string{
			"profile-1 display name",
			"profile-2 display name",
			"profile-3 display name",
		}
		require.ElementsMatch(t, expected, result)

	})
}

func TestFilter(t *testing.T) {
	t.Run("return filtered profiles from both credentials and config files", func(t *testing.T) {
		profiles := Profiles{
			CredentialsProfiles: []Profile{
				{
					ProfileName:        "credentials profile 1",
					DisplayProfileName: "credentials profile 1",
				},
				{
					ProfileName:        "credentials profile 2 - match",
					DisplayProfileName: "credentials profile 2",
				},
			},
			ConfigAssumedProfiles: []Profile{
				{
					ProfileName:        "config profile 1",
					DisplayProfileName: "config profile 1",
				},
				{
					ProfileName:        "config profile 2 - match",
					DisplayProfileName: "config profile 2",
				},
			},
		}

		result := profiles.Filter("match")

		require.Equal(t, 2, len(result))
		require.Equal(t, "credentials profile 2 - match", result[0].ProfileName)
		require.Equal(t, "config profile 2 - match", result[1].ProfileName)

	})

	t.Run("return empty if none matches filter", func(t *testing.T) {
		profiles := Profiles{
			CredentialsProfiles: []Profile{
				{
					ProfileName:        "credentials profile 1",
					DisplayProfileName: "credentials profile 1",
				},
				{
					ProfileName:        "credentials profile 2",
					DisplayProfileName: "credentials profile 2",
				},
			},
			ConfigAssumedProfiles: []Profile{
				{
					ProfileName:        "config profile 1",
					DisplayProfileName: "config profile 1",
				},
				{
					ProfileName:        "config profile 2",
					DisplayProfileName: "config profile 2",
				},
			},
		}

		result := profiles.Filter("match")

		require.Equal(t, 0, len(result))

	})
}
