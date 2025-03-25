package Analyzer

import (
	"Proyecto/FSystem"
	"flag"
	"fmt"
	"strings"
)

// fn_mkfs procesa el comando mkfs.
func fn_mkfs(input string) {
	fs := flag.NewFlagSet("mkfs", flag.ExitOnError)
	id := fs.String("id", "", "Id")
	type_ := fs.String("type", "", "Tipo")
	fs_ := fs.String("fs", "2fs", "Fs")

	// Parsear la cadena de entrada, no os.Args
	matches := paramRegex.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id", "type", "fs":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("\033[1;31mError: Flag not found\033[0m")
		}
	}

	// Verifica que se hayan establecido todas las flags necesarias
	if *id == "" {
		fmt.Println("\033[1;31mError: id es un parámetro obligatorio.\033[0m")
		return
	}

	if *type_ == "" {
		fmt.Println("\033[1;31mError: type es un parámetro obligatorio.\033[0m")
		return
	}

	// Llamar a la función
	FSystem.Mkfs(*id, *type_, *fs_)
}

// fn_cat procesa el comando CAT.
func fn_cat(parametros string) {
	paramMap := extraerParametros(parametros)

	// Obtener la lista de archivos
	var files []string
	for key, value := range paramMap {
		if strings.HasPrefix(key, "file") {
			files = append(files, value)
		}
	}

	// Obtener el ID de la partición
	id := paramMap["id"]
	if id == "" {
		fmt.Println("Error: El parámetro 'id' es obligatorio para el comando CAT")
		return
	}

	// Llamar a la función Cat
	FSystem.Cat(files, id)
}
