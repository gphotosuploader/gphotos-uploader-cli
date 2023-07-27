package config_test

import (
	"bytes"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/configuration"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestConfigDumpCommand(t *testing.T) {
	actual := new(bytes.Buffer)
	configuration.Settings = viper.New()
	configuration.SetDefaults(configuration.Settings)
	rootCommand := cli.New()
	rootCommand.SetOut(actual)
	rootCommand.SetArgs([]string{"config", "dump"})

	err := rootCommand.Execute()
	assert.NoError(t, err)

	// Map to store the parsed YAML data
	var data map[string]interface{}

	// Unmarshal the YAML string into the data map
	err = yaml.Unmarshal(actual.Bytes(), &data)
	assert.NoError(t, err)

	// Test [auth] configuration.
	assert.NotNil(t, data["auth"].(map[string]interface{})["account"])
	assert.NotNil(t, data["auth"].(map[string]interface{})["client_id"])
	assert.NotNil(t, data["auth"].(map[string]interface{})["client_secret"])
	assert.NotNil(t, data["auth"].(map[string]interface{})["secrets_type"])

	// Test [directories] configuration.
	assert.NotNil(t, data["directories"].(map[string]interface{})["data"])

	// Test [logging] configuration.
	assert.NotNil(t, data["logging"].(map[string]interface{})["level"])

	// Test [output] configuration.
	assert.NotNil(t, data["output"].(map[string]interface{})["no_color"])
}
