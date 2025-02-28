package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Values represents the structure of the values.yaml file for Helm
type Values struct {
	ReplicaCount        int                      `yaml:"replicaCount"`
	Image               map[string]string        `yaml:"image"`
	Resources           map[string]interface{}   `yaml:"resources"`
	Env                 []map[string]interface{} `yaml:"env"`
	EnvFrom             []map[string]interface{} `yaml:"envFrom"`
	Volumes             []interface{}            `yaml:"volumes"`
	VolumeMounts        []interface{}            `yaml:"volumeMounts"`
	ImagePullSecrets    []interface{}            `yaml:"imagePullSecrets"`
	NameOverride        string                   `yaml:"nameOverride"`
	FullnameOverride    string                   `yaml:"fullnameOverride"`
	GenericAnnotations  map[string]interface{}   `yaml:"genericAnnotaions"`
	DeploymentAnnotations map[string]interface{} `yaml:"deploymentAnnotations"`
	DeploymentLabels    map[string]interface{}   `yaml:"deplymentLabels"`
	PodSecurityContext  map[string]interface{}   `yaml:"podSecurityContext"`
	SecurityContext     map[string]interface{}   `yaml:"securityContext"`
	Route               map[string]interface{}   `yaml:"route"`
	Service             map[string]interface{}   `yaml:"service"`
	ServiceAccount      map[string]interface{}   `yaml:"serviceAccount"`
	Ingress             map[string]interface{}   `yaml:"ingress"`
	LivenessProbe       map[string]interface{}   `yaml:"livenessProbe"`
	ReadinessProbe      map[string]interface{}   `yaml:"readinessProbe"`
	Autoscaling         map[string]interface{}   `yaml:"autoscaling"`
	NodeSelector        map[string]interface{}   `yaml:"nodeSelector"`
	Tolerations         []interface{}            `yaml:"tolerations"`
	Affinity            map[string]interface{}   `yaml:"affinity"`
}

// initDefaultValues initializes default values
func initDefaultValues() Values {
	values := Values{
		ReplicaCount:     0,
		Image:            map[string]string{"tag": "latest"},
		Resources:        map[string]interface{}{},
		Env:              []map[string]interface{}{},
		EnvFrom:          []map[string]interface{}{},
		Volumes:          []interface{}{},
		VolumeMounts:     []interface{}{},
		ImagePullSecrets: []interface{}{},
		NameOverride:     "",
		FullnameOverride: "",
		GenericAnnotations: map[string]interface{}{
			"enabled": true,
		},
		DeploymentAnnotations: map[string]interface{}{
			"collectord.io/index":  map[string]interface{}{},
			"collectord.io/output": map[string]interface{}{},
		},
		DeploymentLabels:   map[string]interface{}{},
		PodSecurityContext: map[string]interface{}{},
		SecurityContext:    map[string]interface{}{},
		Route: map[string]interface{}{
			"annotations": map[string]interface{}{
				"haproxy.router.openshift.io/balance": "leastconn",
			},
			"tls": map[string]interface{}{
				"termination": "edge",
			},
			"to": map[string]interface{}{
				"kind": "Service",
			},
		},
		Service: map[string]interface{}{
			"type": "ClusterIP",
			"name": "80-tcp",
			"port": 80,
		},
		ServiceAccount: map[string]interface{}{
			"create":     true,
			"automount":  true,
			"annotations": map[string]interface{}{},
			"name":       "",
		},
		Ingress: map[string]interface{}{
			"enabled":     false,
			"className":   "",
			"annotations": map[string]interface{}{},
			"hosts":       []interface{}{},
			"tls":         []interface{}{},
		},
		LivenessProbe:  map[string]interface{}{},
		ReadinessProbe: map[string]interface{}{},
		Autoscaling: map[string]interface{}{
			"enabled":                        false,
			"minReplicas":                    1,
			"maxReplicas":                    100,
			"targetCPUUtilizationPercentage": 80,
		},
		NodeSelector: map[string]interface{}{},
		Tolerations:  []interface{}{},
		Affinity:     map[string]interface{}{},
	}
	return values
}

// deploymentLogic extracts values from the deployment YAML
func deploymentLogic(deploymentYaml string, values *Values) error {
	data, err := os.ReadFile(deploymentYaml)
	if err != nil {
		return fmt.Errorf("error reading deployment file: %w", err)
	}

	var deploymentData map[string]interface{}
	if err := yaml.Unmarshal(data, &deploymentData); err != nil {
		return fmt.Errorf("error unmarshaling deployment data: %w", err)
	}

	// Extract replicaCount
	if spec, ok := deploymentData["spec"].(map[string]interface{}); ok {
		if replicas, ok := spec["replicas"].(int); ok {
			values.ReplicaCount = replicas
		}
	}

	// Extract container info
	if spec, ok := deploymentData["spec"].(map[string]interface{}); ok {
		if template, ok := spec["template"].(map[string]interface{}); ok {
			if templateSpec, ok := template["spec"].(map[string]interface{}); ok {
				if containers, ok := templateSpec["containers"].([]interface{}); ok && len(containers) > 0 {
					if container, ok := containers[0].(map[string]interface{}); ok {
						// Extract image tag
						if image, ok := container["image"].(string); ok {
							parts := strings.Split(image, ":")
							if len(parts) > 1 {
								values.Image["tag"] = parts[1]
							}
						}

						// Extract resources
						if resources, ok := container["resources"].(map[string]interface{}); ok {
							values.Resources = resources
						}

						// Extract env variables
						if envVars, ok := container["env"].([]interface{}); ok {
							for _, e := range envVars {
								if env, ok := e.(map[string]interface{}); ok {
									entry := make(map[string]interface{})
									
									if name, ok := env["name"].(string); ok {
										entry["name"] = name
									}
									
									if value, ok := env["value"]; ok {
										entry["value"] = value
									} else if valueFrom, ok := env["valueFrom"].(map[string]interface{}); ok {
										entry["valueFrom"] = valueFrom
									}
									
									values.Env = append(values.Env, entry)
								}
							}
						}

						// Extract envFrom
						if envFromVars, ok := container["envFrom"].([]interface{}); ok {
							for _, ef := range envFromVars {
								if envFrom, ok := ef.(map[string]interface{}); ok {
									entry := make(map[string]interface{})
									
									if configMapRef, ok := envFrom["configMapRef"].(map[string]interface{}); ok {
										entry["configMapRef"] = map[string]interface{}{
											"name": configMapRef["name"],
										}
									}
									
									if secretRef, ok := envFrom["secretRef"].(map[string]interface{}); ok {
										entry["secretRef"] = map[string]interface{}{
											"name": secretRef["name"],
										}
									}
									
									values.EnvFrom = append(values.EnvFrom, entry)
								}
							}
						}

						// Extract securityContext
						if securityContext, ok := container["securityContext"].(map[string]interface{}); ok {
							values.SecurityContext = securityContext
						}

						// Extract livenessProbe
						if livenessProbe, ok := container["livenessProbe"].(map[string]interface{}); ok {
							values.LivenessProbe = livenessProbe
						}

						// Extract readinessProbe
						if readinessProbe, ok := container["readinessProbe"].(map[string]interface{}); ok {
							values.ReadinessProbe = readinessProbe
						}

						// Extract volumeMounts
						if volumeMounts, ok := container["volumeMounts"].([]interface{}); ok {
							values.VolumeMounts = volumeMounts
						}
					}
				}

				// Extract volumes
				if volumes, ok := templateSpec["volumes"].([]interface{}); ok {
					values.Volumes = volumes
				}

				// Extract podSecurityContext
				if podSecurityContext, ok := templateSpec["securityContext"].(map[string]interface{}); ok {
					values.PodSecurityContext = podSecurityContext
				}

				// Extract nodeSelector
				if nodeSelector, ok := templateSpec["nodeSelector"].(map[string]interface{}); ok {
					values.NodeSelector = nodeSelector
				}

				// Extract tolerations
				if tolerations, ok := templateSpec["tolerations"].([]interface{}); ok {
					values.Tolerations = tolerations
				}

				// Extract affinity
				if affinity, ok := templateSpec["affinity"].(map[string]interface{}); ok {
					values.Affinity = affinity
				}
			}
		}
	}

	// Extract metadata
	if metadata, ok := deploymentData["metadata"].(map[string]interface{}); ok {
		if annotations, ok := metadata["annotations"].(map[string]interface{}); ok {
			if index, ok := annotations["collectord.io/index"]; ok {
				values.DeploymentAnnotations["collectord.io/index"] = index
			}
			if output, ok := annotations["collectord.io/output"]; ok {
				values.DeploymentAnnotations["collectord.io/output"] = output
			}
		}
		if labels, ok := metadata["labels"].(map[string]interface{}); ok {
			values.DeploymentLabels = labels
		}
	}

	return nil
}

// processServiceYaml extracts values from the service YAML
func processServiceYaml(serviceYaml string, values *Values) error {
	data, err := os.ReadFile(serviceYaml)
	if err != nil {
		return fmt.Errorf("error reading service file: %w", err)
	}

	var serviceData map[string]interface{}
	if err := yaml.Unmarshal(data, &serviceData); err != nil {
		return fmt.Errorf("error unmarshaling service data: %w", err)
	}

	if spec, ok := serviceData["spec"].(map[string]interface{}); ok {
		if serviceType, ok := spec["type"].(string); ok {
			values.Service["type"] = serviceType
		}
		
		if ports, ok := spec["ports"].([]interface{}); ok && len(ports) > 0 {
			if port, ok := ports[0].(map[string]interface{}); ok {
				if portNumber, ok := port["port"].(int); ok {
					values.Service["port"] = portNumber
				}
			}
		}
	}

	return nil
}

// processRouteYaml extracts values from the route YAML
func processRouteYaml(routeYaml string, values *Values) error {
	data, err := os.ReadFile(routeYaml)
	if err != nil {
		return fmt.Errorf("error reading route file: %w", err)
	}

	var routeData map[string]interface{}
	if err := yaml.Unmarshal(data, &routeData); err != nil {
		return fmt.Errorf("error unmarshaling route data: %w", err)
	}

	values.Ingress["enabled"] = true
	
	if spec, ok := routeData["spec"].(map[string]interface{}); ok {
		host := ""
		if h, ok := spec["host"].(string); ok {
			host = h
		}
		
		path := ""
		if to, ok := spec["to"].(map[string]interface{}); ok {
			if name, ok := to["name"].(string); ok {
				path = name
			}
		}
		
		values.Ingress["hosts"] = []interface{}{
			map[string]interface{}{
				"host": host,
				"paths": []interface{}{
					map[string]interface{}{
						"path": path,
					},
				},
			},
		}
	}

	return nil
}

// generateValuesYaml generates the values.yaml file
func generateValuesYaml(deploymentYaml, serviceYaml, routeYaml string) error {
	// Initialize default values
	values := initDefaultValues()

	// Process deployment YAML
	if deploymentYaml != "" {
		if err := deploymentLogic(deploymentYaml, &values); err != nil {
			return err
		}
	}

	// Process service YAML
	if serviceYaml != "" {
		if err := processServiceYaml(serviceYaml, &values); err != nil {
			return err
		}
	}

	// Process route YAML
	if routeYaml != "" {
		if err := processRouteYaml(routeYaml, &values); err != nil {
			return err
		}
	}

	// Write values to file
	additionalComments := "### THIS VALUE FILE MADE BY SRE ALPHA \n"
	
	file, err := os.Create("values.yaml")
	if err != nil {
		return fmt.Errorf("error creating values.yaml file: %w", err)
	}
	defer file.Close()
	
	if _, err := file.WriteString(additionalComments); err != nil {
		return fmt.Errorf("error writing comments to values.yaml: %w", err)
	}
	
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	if err := encoder.Encode(values); err != nil {
		return fmt.Errorf("error encoding values to YAML: %w", err)
	}
	
	fmt.Println("values.yaml generated successfully.")
	return nil
}

func printHelp() {
	helpText := `
Hi welcome !
Your using pico binary file that convert deployment, service and route into Helm values file

USAGE:
The right syntax is:
pico [DEPLOYMENT_FILE_PATH] [SERVICE_FILE_PATH] [ROUTE_FILE_PATH]
You must add at least deployment file for using the binary, Its better if the deployment is neat
If you dont have kubectl-neat binary shame on you !
BTW you can ask SRE alpha team for the neat binary :)

CAUTION: 
This will make a values.yaml file in your run folder, It can run-over files with the name values.yaml so be careful !!

LIFE-HACK:
Be good be opensource :)
`
	fmt.Println(helpText)
}

func main() {
	// Define flags
	helpFlag := flag.Bool("help", false, "Display help information")
	flag.BoolVar(helpFlag, "h", false, "Display help information (shorthand)")
	
	// Parse flags
	flag.Parse()
	
	// Check if help flag is set
	if *helpFlag {
		printHelp()
		return
	}
	
	// Get positional arguments
	args := flag.Args()
	
	// Check if at least one argument is provided
	if len(args) == 0 {
		fmt.Println("You must add at least one yaml file to extract values from !!!")
		fmt.Println("EX:")
		fmt.Println("pico ./my_deployment.yaml ./my_service.yaml ./my_route.yaml")
		fmt.Println("For more info use:")
		fmt.Println("pico -h or pico --help")
		return
	}
	
	deploymentYaml := ""
	serviceYaml := ""
	routeYaml := ""
	
	if len(args) > 0 {
		deploymentYaml = args[0]
	}
	
	if len(args) > 1 {
		serviceYaml = args[1]
	}
	
	if len(args) > 2 {
		routeYaml = args[2]
	}
	
	fmt.Printf("files are %s and %s and %s\n", deploymentYaml, serviceYaml, routeYaml)
	
	err := generateValuesYaml(deploymentYaml, serviceYaml, routeYaml)
	if err != nil {
		fmt.Printf("Probably your file isn't valid!\nError: %s\n", err)
		os.Exit(1)
	}
}