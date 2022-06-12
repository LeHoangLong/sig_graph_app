package common

import (
	"fmt"
	"os"

	"go.uber.org/dig"
	"gopkg.in/yaml.v3"
)

func ProvideConfig(
	iContainer *dig.Container,
) error {
	return iContainer.Provide(func() (Config, error) {
		configPath, found := os.LookupEnv("CONFIG_FILE")
		if !found {
			return Config{}, NotFound
		}

		configYaml, err := os.ReadFile(configPath)
		if err != nil {
			return Config{}, fmt.Errorf("Could not read config.yaml file %s", err.Error())
		}

		var config Config
		err = yaml.Unmarshal(configYaml, &config)
		if err != nil {
			return Config{}, fmt.Errorf("Could not parse yaml file %s", err.Error())
		}

		return config, nil
	})
}
