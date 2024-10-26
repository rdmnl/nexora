package config

import (
	"log"
	"os"

	"github.com/rdmnl/nexora/shared"

	"gopkg.in/yaml.v3"
)

func LoadConfig(configPath string) *shared.Config {
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("Could not read %s: %v. Using only detected nodes.", configPath, err)
		return &shared.Config{}
	}

	var config shared.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Printf("Error parsing %s: %v. Using only detected nodes.", configPath, err)
	}
	return &config
}
