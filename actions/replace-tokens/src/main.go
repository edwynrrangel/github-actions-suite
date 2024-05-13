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

	fmt.Printf("data: %s\n", data)

	// Convierte los datos a string para procesamiento
	content := string(data)

	// Expresión regular para encontrar todos los tokens
	re := regexp.MustCompile(`#{(.*?)}#`)

	// Reemplaza cada token encontrado por su valor de variable de entorno correspondiente
	replacedContent := re.ReplaceAllStringFunc(content, func(token string) string {
		// Extrae el nombre de la variable de entorno del token
		varName := token[2 : len(token)-2] // Elimina los delimitadores #{ y }#
		// Obtiene el valor de la variable de entorno
		envValue := os.Getenv(varName)
		return envValue
	})

	// Imprime el contenido final para verificar
	fmt.Printf("Modified content: %s\n", replacedContent)
}
