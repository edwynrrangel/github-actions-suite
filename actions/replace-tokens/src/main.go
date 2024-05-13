package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
)

func main() {
	pathInputFile := os.Getenv("INPUT_YAML_FILE")
	if pathInputFile == "" {
		panic("no yaml file provided")
	}
	fmt.Printf("YAML file: %s\n", pathInputFile)
	yamlFile, err := os.Open(pathInputFile)
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(yamlFile)
	if err != nil {
		panic(err)
	}
	yamlFile.Close()

	re := regexp.MustCompile(`#{(.*?)}#`)
	replacedContent := re.ReplaceAllStringFunc(string(data), func(token string) string {
		// Extrae el nombre de la variable de entorno del token
		varName := token[2 : len(token)-2] // Elimina los delimitadores #{ y }#
		envValue := os.Getenv(varName)
		return envValue
	})

	if os.Getenv("OUTPUT_YAML_FILE") == "" {
		err = os.WriteFile(pathInputFile, []byte(replacedContent), 0644)
		if err != nil {
			panic(err)
		}
	}

	if os.Getenv("OUTPUT_YAML_FILE") != "" {
		err = os.WriteFile(os.Getenv("OUTPUT_YAML_FILE"), []byte(replacedContent), 0644)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Replace tokens successfully\n%s##########\n", replacedContent)
}
