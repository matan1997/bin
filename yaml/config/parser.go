package config

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig reads a JSON file and returns a DeploymentConfig struct
func LoadConfig(filePath string) DeploymentConfig {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return DeploymentConfig{}
	}

	var config DeploymentConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	return config
}

// ConvertToYAML converts DeploymentConfig to YAML
func ConvertToYAML(config DeploymentConfig) string {
	data, err := yaml.Marshal(config)
	if err != nil {
		fmt.Println("Error converting to YAML:", err)
		return ""
	}
	return string(data)
}
