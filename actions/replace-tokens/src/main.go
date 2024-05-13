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
	yamlFile, err := os.Open(os.Getenv("INPUT_YAML_FILE"))
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
		return envValue
	})

	fmt.Printf("Modified content:\n##########%s##########\n", replacedContent)

	if os.Getenv("OUTPUT_YAML_FILE") == "" {
		err = os.WriteFile(os.Getenv("INPUT_YAML_FILE"), []byte(replacedContent), 0644)
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
}
