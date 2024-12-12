package railway

import (
	"errors"
	"os"
)

type Config struct {
	ApiKey        string
	EnvironmentId string
	ServiceIds    []string
}

func GenerateConfig(serviceIds []string) (*Config, error) {
	config := Config{}

	apiKey := os.Getenv("RAILWAY_API_KEY")
	environtmentId := os.Getenv("RAILWAY_ENVIRONMENT_ID")

	if apiKey == "" {
		return nil, errors.New("api key must be present")
	}

	if environtmentId == "" {
		return nil, errors.New("environment id must be present")
	}

	config.ServiceIds = serviceIds

	return &config, nil
}
