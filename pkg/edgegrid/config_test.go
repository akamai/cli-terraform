package edgegrid

import (
	"flag"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/edgegrid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestGetEdgegridConfig(t *testing.T) {
	tests := map[string]struct {
		configFile     string
		configSection  string
		flagAccountKey string
		configEnvs     map[string]string
		expectedConfig edgegrid.Config
		withError      bool
	}{
		"valid config from file": {
			configFile:    "./testdata/.edgerc",
			configSection: "test_section",
			expectedConfig: edgegrid.Config{
				Host:         "akaa-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX.luna.akamaiapis.net/",
				ClientToken:  "akab-XXXXXXXXXXXXXXXX-XXXXXXXXXXXXXXXX",
				ClientSecret: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				AccessToken:  "akab-XXXXXXXXXXXXXXXX-XXXXXXXXXXXXXXXX",
			},
		},
		"valid config with accountkey": {
			configFile:     "./testdata/.edgerc",
			configSection:  "test_section",
			flagAccountKey: "XXXXXXXX",
			expectedConfig: edgegrid.Config{
				Host:         "akaa-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX.luna.akamaiapis.net/",
				ClientToken:  "akab-XXXXXXXXXXXXXXXX-XXXXXXXXXXXXXXXX",
				ClientSecret: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				AccessToken:  "akab-XXXXXXXXXXXXXXXX-XXXXXXXXXXXXXXXX",
				AccountKey:   "XXXXXXXX",
			},
		},
		"valid config from file, envs passed": {
			configFile:    "./testdata/.edgerc",
			configSection: "test_section",
			configEnvs: map[string]string{
				"AKAMAI_TEST_SECTION_HOST":          "env-host",
				"AKAMAI_TEST_SECTION_CLIENT_TOKEN":  "env-client-token",
				"AKAMAI_TEST_SECTION_CLIENT_SECRET": "env-client-secret",
				"AKAMAI_TEST_SECTION_ACCESS_TOKEN":  "env-access-token",
			},
			expectedConfig: edgegrid.Config{
				Host:         "env-host",
				ClientToken:  "env-client-token",
				ClientSecret: "env-client-secret",
				AccessToken:  "env-access-token",
			},
		},
		"valid config from file, envs passed, with accountkey flag": {
			configFile:     "./testdata/.edgerc",
			configSection:  "test_section",
			flagAccountKey: "flag-account-key",
			configEnvs: map[string]string{
				"AKAMAI_TEST_SECTION_HOST":          "env-host",
				"AKAMAI_TEST_SECTION_CLIENT_TOKEN":  "env-client-token",
				"AKAMAI_TEST_SECTION_CLIENT_SECRET": "env-client-secret",
				"AKAMAI_TEST_SECTION_ACCESS_TOKEN":  "env-access-token",
				"AKAMAI_TEST_SECTION_ACCOUNT_KEY":   "env-account-key",
			},
			expectedConfig: edgegrid.Config{
				Host:         "env-host",
				ClientToken:  "env-client-token",
				ClientSecret: "env-client-secret",
				AccessToken:  "env-access-token",
				AccountKey:   "flag-account-key",
			},
		},
		"valid config from file, envs passed, no accountkey flag": {
			configFile:     "./testdata/.edgerc",
			configSection:  "test_section",
			flagAccountKey: "",
			configEnvs: map[string]string{
				"AKAMAI_TEST_SECTION_HOST":          "env-host",
				"AKAMAI_TEST_SECTION_CLIENT_TOKEN":  "env-client-token",
				"AKAMAI_TEST_SECTION_CLIENT_SECRET": "env-client-secret",
				"AKAMAI_TEST_SECTION_ACCESS_TOKEN":  "env-access-token",
				"AKAMAI_TEST_SECTION_ACCOUNT_KEY":   "env-account-key",
			},
			expectedConfig: edgegrid.Config{
				Host:         "env-host",
				ClientToken:  "env-client-token",
				ClientSecret: "env-client-secret",
				AccessToken:  "env-access-token",
				AccountKey:   "env-account-key",
			},
		},
		"invalid edgerc file": {
			configFile:    "./testdata/edgerc-invalid",
			configSection: "test_section",
			withError:     true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			app := cli.NewApp()
			set := flag.NewFlagSet("test", 0)
			set.String("edgerc", "~/.egderc", "")
			set.String("section", "default", "")
			flags := []string{"--edgerc", test.configFile, "--section", test.configSection}
			if test.flagAccountKey != "" {
				set.String("accountkey", "", "")
				flags = append(flags, "--accountkey", test.flagAccountKey)
			}
			err := set.Parse(flags)
			assert.NoError(t, err)
			cliCtx := cli.NewContext(app, set, nil)
			for k, v := range test.configEnvs {
				require.NoError(t, os.Setenv(k, v))
			}
			defer func() {
				for k := range test.configEnvs {
					require.NoError(t, os.Unsetenv(k))
				}
			}()
			res, err := GetEdgegridConfig(cliCtx)
			if test.withError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, res)
			assert.Equal(t, test.expectedConfig.Host, res.Host)
			assert.Equal(t, test.expectedConfig.ClientToken, res.ClientToken)
			assert.Equal(t, test.expectedConfig.ClientSecret, res.ClientSecret)
			assert.Equal(t, test.expectedConfig.AccessToken, res.AccessToken)
			assert.Equal(t, test.expectedConfig.AccountKey, res.AccountKey)
		})
	}
}

func TestGetEdgercPath(t *testing.T) {
	tests := map[string]struct {
		edgercPath string
		expected   string
	}{
		"edgerc flag provided": {
			edgercPath: "/some/path",
			expected:   "/some/path",
		},
		"edgerc flag not provided, return default": {
			expected: "~/.edgerc",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			app := cli.NewApp()
			set := flag.NewFlagSet("test", 0)
			set.String("edgerc", test.edgercPath, "")
			cliCtx := cli.NewContext(app, set, nil)
			res := GetEdgercPath(cliCtx)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestGetEdgercSection(t *testing.T) {
	tests := map[string]struct {
		section  string
		expected string
	}{
		"section flag provided": {
			section:  "/some/path",
			expected: "/some/path",
		},
		"section flag not provided, return default": {
			expected: "default",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			app := cli.NewApp()
			set := flag.NewFlagSet("test", 0)
			set.String("section", test.section, "")
			cliCtx := cli.NewContext(app, set, nil)
			res := GetEdgercSection(cliCtx)
			assert.Equal(t, test.expected, res)
		})
	}
}
