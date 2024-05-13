package main

import (
	"fmt"
	"io"
	"os"
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

	// Obtener la ruta del archivo de entorno
	envFile := os.Getenv("GITHUB_ENV")
	if envFile == "" {
		panic("GITHUB_ENV not set")
	}

	// Abrir el archivo en modo de solo lectura
	envfile, err := os.OpenFile(envFile, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer envfile.Close()

	// Leer todo el contenido del archivo
	envs, err := io.ReadAll(envfile)
	if err != nil {
		panic(err)
	}

	// Imprimir el contenido para verificación
	fmt.Printf("Environment file content: %s\n", envs)
}
