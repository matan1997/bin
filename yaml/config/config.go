package config

type DeploymentConfig struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Replicas int `json:"replicas"`
}
