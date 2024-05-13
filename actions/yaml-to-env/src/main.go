package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

func main() {
	if os.Getenv("INPUT_YAML_FILE") == "" {
		panic("no yaml file provided")
	}
	yamlFile, err := os.Open(os.Getenv("INPUT_YAML_FILE"))
	if err != nil {
		panic(err)
	}
	defer yamlFile.Close()

	data, err := io.ReadAll(yamlFile)
	if err != nil {
		panic(err)
	}

	var envs map[string]interface{}
	err = yaml.Unmarshal(data, &envs)
	if err != nil {
		panic(err)
	}

	mapEnvs := mapEnvironmentVariables("", envs)

	envFile := os.Getenv("GITHUB_ENV")
	if envFile == "" {
		panic("GITHUB_ENV not set")
	}

	file, err := os.OpenFile(envFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Printf("Environment variables:\n")
	for key, value := range mapEnvs {
		fmt.Printf("%s=%s\n", key, value)
		_, err := file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			panic(err)
		}
	}
}

func mapEnvironmentVariables(prefix string, envs interface{}) map[string]string {
	mapEnvs := make(map[string]string)
	switch envs := envs.(type) {
	case map[string]interface{}:
		for key, value := range envs {
			fullKey := constructFullKey(prefix, key)
			for k, v := range mapEnvironmentVariables(fullKey, value) {
				mapEnvs[k] = v
			}
		}
	case []interface{}:
		for i, item := range envs {
			indexedKey := fmt.Sprintf("%s_%d", prefix, i)
			for k, v := range mapEnvironmentVariables(indexedKey, item) {
				mapEnvs[k] = v
			}
		}
	case map[interface{}]interface{}:
		normalizedMap := convertMap(envs)
		for k, v := range mapEnvironmentVariables(prefix, normalizedMap) {
			mapEnvs[k] = v
		}
	default:
		if prefix != "" {
			mapEnvs[prefix] = fmt.Sprintf("%v", envs)
		}
	}

	return mapEnvs
}

func constructFullKey(prefix, key string) string {
	if prefix == "" {
		return camelToSnake(key)
	}
	return fmt.Sprintf("%s_%s", prefix, camelToSnake(key))
}

func convertMap(orig map[interface{}]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{})
	for k, v := range orig {
		if keyStr, ok := k.(string); ok {
			normalized[keyStr] = v
		}
	}
	return normalized
}

func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result.WriteByte('_')
		}
		result.WriteRune(unicode.ToUpper(r))
	}
	return result.String()
}
