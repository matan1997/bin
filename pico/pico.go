package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	// "text/template"
)

const (
    // Define a constant for a maximum number of retries
    MaxRetries = 3

    // Define a constant for a default timeout in seconds
    DefaultTimeout = 30
)

type Values struct {
	ReplicaCount int               `yaml:"replicaCount"`
	Image        map[string]string `yaml:"image"`
	Resources    map[string]interface{} `yaml:"resources"`
	Env          []map[string]interface{} `yaml:"env"`
	Service      struct {
		Type string `yaml:"type"`
		Port int    `yaml:"port"`
	} `yaml:"service"`
}

var rootCmd = &cobra.Command{
	Use:   "pico",
	Example: "pico <deployment_file_path>",
	Short: "Convert Kubernetes YAML to Helm values.yaml",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deploymentFile := args[0]
		serviceFile := ""
		if len(args) > 1 {
			serviceFile = args[1]
		}
		convertYAML(deploymentFile, serviceFile)
	},
}

func convertYAML(deploymentFile, serviceFile string) {
	values := Values{
		ReplicaCount: 1,
		Image:        map[string]string{"tag": "latest"},
		Resources:    map[string]interface{}{},
		Env:          []map[string]interface{}{},
	} 

	if deploymentFile != "" {
		data, err := os.ReadFile(deploymentFile)
		if err != nil {
			fmt.Println("Error reading deployment file:", err)
			return
		}
		var deployment map[string]interface{}
		yaml.Unmarshal(data, &deployment)
		
		if spec, exists := deployment["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{}); exists {
			if replicas, exists := spec["replicas"].(int); exists {
				fmt.Println("replicas", replicas)
				values.ReplicaCount = replicas
			}
			if envs, exists := spec["containers"].([]interface{})[0].(map[string]interface{})["env"].([]interface{}); exists {
				for _, e := range envs {
					if envMap, ok := e.(map[string]interface{}); ok {
						values.Env = append(values.Env, envMap)
						fmt.Println(envMap)
					}
				}
			}else{
				fmt.Println("shit")
			}
		}else{
			fmt.Println("nooo")
		}

		
		
			
			// if envs, exists := spec["env"].([]map[string]string); exists {
			// 	fmt.Println(envs)
			// 	values.Env = envs
			// }

	

		// fmt.Println(err)
	}

	if serviceFile != "" {
		data, err := os.ReadFile(serviceFile)
		
		if err != nil {
			fmt.Println("Error reading service file:", err)
			return
		}
		
		var service map[string]interface{}
		yaml.Unmarshal(data, &service)
		if spec, ok := service["spec"].(map[string]interface{}); ok {
			values.Service.Type = spec["type"].(string)
			values.Service.Port = int(spec["ports"].([]interface{})[0].(map[string]interface{})["port"].(int))
		}
	}

	
	
	/*
		test doc
	*/
	// var MarkdownDocTemplate = `% hoooo {{ .Image }} {{ .ReplicaCount }}
	// #dssd
	// % jfjfjf
	// %{{ .Env }}`
	// tpl, err := template.New("doc").Parse(MarkdownDocTemplate)
	// if err != nil {
	// 	panic(err) // Handle the error properly
	// }
	// tpl.Execute(os.Stdout, values)

	
	
	output, err := yaml.Marshal(values)
	if err != nil {
		fmt.Println("Error generating values.yaml:", err)
		return
	}

	outputs := []byte("# matan made\n" + string(output))
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	os.WriteFile("values.yaml", outputs, 0644)
	fmt.Println("values.yaml generated successfully ❤️")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
