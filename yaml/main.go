package main


import (
	"fmt"
	"yaml/config"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./convert <config.json>")
		return
	}

	filePath := os.Args[1]
	cfg := config.LoadConfig(filePath)

	yamlStr := config.ConvertToYAML(cfg)
	if yamlStr == "" {
		fmt.Println("Failed to convert config.")
		return
	}

	outputFile := "deployment.yaml"
	err := os.WriteFile(outputFile, []byte(yamlStr), 0644)
	if err != nil {
		fmt.Println("Error writing YAML file:", err)
		return
	}

	fmt.Println("Converted YAML saved to", outputFile)
	
}

