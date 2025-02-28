package utils

import "strings"

// FormatYAML adds indentation to make the YAML more readable
func FormatYAML(yamlStr string) string {
	return strings.ReplaceAll(yamlStr, "\n", "\n  ")
}
