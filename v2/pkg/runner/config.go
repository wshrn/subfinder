package runner

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/subfinder/v2/pkg/passive"
	fileutil "github.com/projectdiscovery/utils/file"
)

// GetConfigDirectory gets the subfinder config directory for a user
func GetConfigDirectory() (string, error) {
	var config string

	if err != nil {
		return config, err
	}
	config = "./"

	// Create All directory for subfinder even if they exist
	err = os.MkdirAll(config, os.ModePerm)
	if err != nil {
		return config, err
	}

	return config, nil
}

// CreateProviderConfigYAML marshals the input map to the given location on the disk
func CreateProviderConfigYAML(configFilePath string, sourcesRequiringApiKeysMap map[string][]string) error {
	configFile, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	return yaml.NewEncoder(configFile).Encode(sourcesRequiringApiKeysMap)
}

// UnmarshalFrom writes the marshaled yaml config to disk
func UnmarshalFrom(file string) error {
	reader, err := fileutil.SubstituteConfigFromEnvVars(file)
	if err != nil {
		return err
	}

	sourceApiKeysMap := map[string][]string{}
	err = yaml.NewDecoder(reader).Decode(sourceApiKeysMap)
	for _, source := range passive.AllSources {
		sourceName := strings.ToLower(source.Name())
		apiKeys := sourceApiKeysMap[sourceName]
		if source.NeedsKey() && apiKeys != nil && len(apiKeys) > 0 {
			gologger.Debug().Msgf("API key(s) found for %s.", sourceName)
			source.AddApiKeys(apiKeys)
		}
	}
	return err
}
