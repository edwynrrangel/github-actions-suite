package main

import (
	"fmt"
	"io"
	"os"
	"strings"

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
	fmt.Printf("Environment variables: %+v\n", mapEnvs)

	envFile := os.Getenv("GITHUB_ENV")
	if envFile == "" {
		panic("GITHUB_ENV not set")
	}

	file, err := os.OpenFile(envFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString("TEST_VAR=\"test_value\"\n")
	if err != nil {
		panic(err)
	}

	// Write each environment variable to the file
	for key, value := range mapEnvs {
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
			mapEnvironmentVariables(fullKey, value)
		}
	case []interface{}:
		for i, item := range envs {
			indexedKey := fmt.Sprintf("%s_%d", prefix, i)
			mapEnvironmentVariables(indexedKey, item)
		}
	case map[interface{}]interface{}:
		normalizedMap := convertMap(envs)
		mapEnvironmentVariables(prefix, normalizedMap)
	default:
		mapEnvs[prefix] = fmt.Sprintf("\"%v\"", envs)
	}

	return mapEnvs
}

func constructFullKey(prefix, key string) string {
	if prefix == "" {
		return strings.ToUpper(key)
	}
	return strings.ToUpper(fmt.Sprintf("%s_%s", prefix, key))
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
