package Analyzer

import (
	"Proyecto/Environment"
	"Proyecto/FSystem"
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var paramRegex = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

// getEntrada recibe una cadena de texto y devuelve el comando y los parámetros.
func getEntrada(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		comandos := strings.ToLower(parts[0]) // Convertir comando a minúsculas
		parametros := strings.Join(parts[1:], " ")
		return comandos, parametros
	}
	return "", input
}

// extraerParametros extrae los parámetros de una cadena usando una expresión regular.
func extraerParametros(parametros string) map[string]string {
	matches := paramRegex.FindAllStringSubmatch(parametros, -1)
	paramMap := make(map[string]string)

	for _, match := range matches {
		flagName := strings.ToLower(match[1]) // Convertir nombre del parámetro a minúsculas
		flagValue := strings.Trim(match[2], "\"")
		paramMap[flagName] = flagValue
	}

	return paramMap
}

// AnalizadorEntrada es la función principal que recibe los comandos del usuario.
func AnalizadorEntrada() {
	for {
		fmt.Println("\033[1;36m╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                      \033[1;33mINGRESE UN COMANDO\033[1;36m                              ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝\033[0m")
		fmt.Print("Ingresar un Comando: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		comandos, parametros := getEntrada(input)
		fmt.Println("Comando: ", comandos, " - ", "Parametro: ", parametros)

		AnalizadorComandos(comandos, parametros)
	}
}

// AnalizadorComandos redirige el comando a la función correspondiente.
func AnalizadorComandos(comandos string, parametros string) {
	switch {
	case strings.Contains(comandos, "mkdisk"):
		fn_mkdisk(parametros)
	case strings.Contains(comandos, "fdisk"):
		fn_fdisk(parametros)
	case strings.Contains(comandos, "mount"):
		fn_mount(parametros)
	case strings.Contains(comandos, "rmdisk"):
		fn_rmdisk(parametros)
	case strings.Contains(comandos, "mkfs"):
		fn_mkfs(parametros)
	case strings.Contains(comandos, "rep"):
		fmt.Println("COMANDO REP")
	default:
		fmt.Println("Error: Comando inválido o no encontrado")
	}
}

// fn_mkdisk procesa el comando mkdisk.
func fn_mkdisk(parametros string) {
	paramMap := extraerParametros(parametros)

	// Valores por defecto
	size, err := strconv.Atoi(paramMap["size"])
	if err != nil || size <= 0 {
		fmt.Println("Error: El valor de 'size' debe ser un número entero mayor que 0")
		return
	}

	fit := strings.ToLower(paramMap["fit"])
	if fit == "" {
		fit = "ff" // Valor por defecto
	} else if fit != "bf" && fit != "ff" && fit != "wf" {
		fmt.Println("Error: El valor de 'fit' debe ser 'bf', 'ff' o 'wf'")
		return
	}

	unit := strings.ToLower(paramMap["unit"])
	if unit == "" {
		unit = "m" // Valor por defecto
	} else if unit != "k" && unit != "m" {
		fmt.Println("Error: El valor de 'unit' debe ser 'k' o 'm'")
		return
	}

	path := strings.ToLower(paramMap["path"]) // Convertir ruta a minúsculas
	if path == "" {
		fmt.Println("Error: El parámetro 'path' es obligatorio")
		return
	}

	// Llamar a la función Mkdisk con los parámetros procesados
	Environment.Mkdisk(size, fit, unit, path)
}

// fn_fdisk procesa el comando fdisk.
func fn_fdisk(parametros string) {
	paramMap := extraerParametros(parametros)

	// Valores por defecto
	unit := strings.ToLower(paramMap["unit"])
	if unit == "" {
		unit = "k" // Valor por defecto: kilobytes
	} else if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: La unidad debe ser 'b', 'k' o 'm'")
		return
	}

	fit := strings.ToLower(paramMap["fit"])
	if fit == "" {
		fit = "wf" // Valor por defecto: Worst Fit
	} else if fit != "bf" && fit != "ff" && fit != "wf" {
		fmt.Println("Error: El ajuste debe ser 'bf', 'ff' o 'wf'")
		return
	}

	partType := strings.ToLower(paramMap["type"])
	if partType == "" {
		partType = "p" // Valor por defecto: Primaria
	} else if partType != "p" && partType != "e" && partType != "l" {
		fmt.Println("Error: El tipo de partición debe ser 'p', 'e' o 'l'")
		return
	}

	// Validar parámetros obligatorios
	size, err := strconv.Atoi(paramMap["size"])
	if err != nil || size <= 0 {
		fmt.Println("Error: El valor de 'size' debe ser un número entero mayor que 0")
		return
	}

	name := strings.ToLower(paramMap["name"]) // Convertir nombre a minúsculas
	if name == "" {
		fmt.Println("Error: El nombre de la partición es obligatorio")
		return
	}

	path := strings.ToLower(paramMap["path"]) // Convertir ruta a minúsculas
	if path == "" {
		fmt.Println("Error: El parámetro 'path' es obligatorio")
		return
	}

	// Convertir el tamaño a bytes
	sizeInBytes := size
	switch unit {
	case "k":
		sizeInBytes *= 1024
	case "m":
		sizeInBytes *= 1024 * 1024
	}

	// Llamar a la función FDISK con los parámetros procesados
	Environment.Fdisk(sizeInBytes, path, name, unit, partType, fit)
}

// fn_mount procesa el comando mount.
func fn_mount(parametros string) {
	paramMap := extraerParametros(parametros)

	path := strings.ToLower(paramMap["path"]) // Convertir ruta a minúsculas
	name := strings.ToLower(paramMap["name"]) // Convertir nombre a minúsculas

	if path == "" || name == "" {
		fmt.Println("Error: Path y Name son obligatorios")
		return
	}

	// Llamar a la función Mount con los parámetros procesados
	Environment.Mount(path, name)
}

// fn_rmdisk procesa el comando rmdisk.
func fn_rmdisk(parametros string) {
	paramMap := extraerParametros(parametros)

	path := strings.ToLower(paramMap["path"]) // Convertir ruta a minúsculas
	if path == "" {
		fmt.Println("Error: El parámetro 'path' es obligatorio")
		return
	}

	// Llamar a la función Rmdisk del paquete Environment
	Environment.Rmdisk(path)
}

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
			fmt.Println("Error: Flag not found")
		}
	}

	// Verifica que se hayan establecido todas las flags necesarias
	if *id == "" {
		fmt.Println("Error: id es un parámetro obligatorio.")
		return
	}

	if *type_ == "" {
		fmt.Println("Error: type es un parámetro obligatorio.")
		return
	}

	// Llamar a la función
	FSystem.Mkfs(*id, *type_, *fs_)
}
