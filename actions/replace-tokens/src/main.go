package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
)

func main() {
	inputFilePath := os.Getenv("INPUT_YAML_FILE")
	if inputFilePath == "" {
		panic("no yaml file provided")
	}
	fmt.Printf("YAML file: %s\n", inputFilePath)
	// Añadir debugging para confirmar la ruta del archivo
	fmt.Println("Intentando abrir el archivo:", inputFilePath)
	yamlFile, err := os.Open(inputFilePath)
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

	outputFilePath := os.Getenv("OUTPUT_YAML_FILE")
	if outputFilePath == "" {
		outputFilePath = inputFilePath
	}

	err = os.WriteFile(outputFilePath, []byte(replacedContent), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Replace tokens successfully\n%s##########\n", replacedContent)
}
