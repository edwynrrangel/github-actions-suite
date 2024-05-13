package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
)

func main() {
	if os.Getenv("INPUT_YAML_FILE") == "" {
		panic("no yaml file provided")
	}
	basePath := "/github/workspace/"
	pathInputFile := fmt.Sprintf("%s%s", basePath, os.Getenv("INPUT_YAML_FILE"))
	fmt.Printf("YAML file: %s\n", pathInputFile)

	yamlFile, err := os.Open(pathInputFile)
	if err != nil {
		panic(err)
	}
	defer yamlFile.Close()

	data, err := io.ReadAll(yamlFile)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`#{(.*?)}#`)
	replacedContent := re.ReplaceAllStringFunc(string(data), func(token string) string {
		// Extrae el nombre de la variable de entorno del token
		varName := token[2 : len(token)-2] // Elimina los delimitadores #{ y }#
		envValue := os.Getenv(varName)
		if envValue == "" {
			return token
		}
		return envValue
	})

	pathOutputFile := pathInputFile
	if os.Getenv("OUTPUT_YAML_FILE") != "" {
		pathOutputFile = fmt.Sprintf("%s%s", basePath, os.Getenv("OUTPUT_YAML_FILE"))
	}

	err = os.WriteFile(pathOutputFile, []byte(replacedContent), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Replaced tokens successfully\n##########\n%s\n##########\n", replacedContent)
}
