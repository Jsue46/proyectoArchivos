package Analyzer

import (
	"Proyecto/Environment"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

<<<<<<< HEAD
// input := "mkdisk -size=3000 -unit=K -fit=BF -path=/home/daaniieel/Disco/disk1.bin"
// input := "fdisk -size=500 -unit=K -path=/home/daaniieel/Disco/disk1.bin -type=P -fit=b -name=Particion1"
var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
=======
//input := "mkdisk -size=3000 -unit=K -fit=BF -path=/home/daaniieel/Disco/disk1.bin"
>>>>>>> origin/main

// getEntrada recibe una cadena de texto y devuelve el comando y los parámetros.
func getEntrada(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		params := strings.Join(parts[1:], " ")
		return command, params
	}
	return "", input
}

// Analizador es la función principal que recibe los comandos del usuario.
func AnalizadorEntrada() {

	for {
		var input string
		fmt.Println("\033[1;36m╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                      \033[1;33mINGRESE UN COMANDO\033[1;36m                              ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝\033[0m")
		fmt.Println("Ingresar un Comando: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input = scanner.Text()

		command, params := getEntrada(input)

		fmt.Println("Comando: ", command, " - ", "Parametro: ", params)

		AnalizadorComandos(command, params)

	}
}

func AnalizadorComandos(command string, params string) {
<<<<<<< HEAD
	if strings.Contains(command, "mkdisk") {
		fn_mkdisk(params)
	} else if strings.Contains(command, "fdisk") {
		fn_fdisk(params)
	} else if strings.Contains(command, "rep") {
		fmt.Print("COMANDO REP")
	} else {
		fmt.Println("Error: Comando inválido o no encontrado")
	}
=======

	if strings.Contains(command, "mkdisk") {
		fn_mkdisk(params)
	} else if strings.Contains(command, "rep") {
		fmt.Print("COMANDO REP")
	} else {
		fmt.Println("Error: Commando invalido o no encontrado")
	}

>>>>>>> origin/main
}

// fn_mkdisk procesa el comando mkdisk con sus parámetros.
func fn_mkdisk(params string) {
	// Variables para almacenar los parámetros
	var size int
	var fit, unit, path string

	// Expresión regular para extraer los parámetros
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(params, -1)

	// Procesar los parámetros
	for _, match := range matches {
		flagName := match[1]                   // Nombre del flag (size, fit, unit, path)
		flagValue := strings.ToLower(match[2]) // Valor del flag (en minúsculas)
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size":
			// Convertir el valor a entero
			sizeVal, err := strconv.Atoi(flagValue)
			if err != nil {
				fmt.Println("Error: El valor de 'size' debe ser un número entero")
				return
			}
			size = sizeVal
		case "fit":
			fit = flagValue
		case "unit":
			unit = flagValue
		case "path":
			path = flagValue
		default:
			fmt.Println("Error: Parámetro no reconocido:", flagName)
			return
		}
	}

	// Validaciones
	if size <= 0 {
		fmt.Println("Error: El valor de 'size' debe ser mayor que 0")
		return
	}

	if fit != "bf" && fit != "ff" && fit != "wf" {
		fmt.Println("Error: El valor de 'fit' debe ser 'bf', 'ff' o 'wf'")
		return
	}

	if unit != "k" && unit != "m" {
		fmt.Println("Error: El valor de 'unit' debe ser 'k' o 'm'")
		return
	}

	if path == "" {
		fmt.Println("Error: El parámetro 'path' es obligatorio")
		return
	}

	// Llamar a la función Mkdisk con los parámetros procesados
	Environment.Mkdisk(size, fit, unit, path)
}
<<<<<<< HEAD

// fn_fdisk procesa el comando fdisk con sus parámetros.
func fn_fdisk(input string) {
	// Variables para almacenar los parámetros
	var size int
	var unit, path, partType, fit, name string

	// Expresión regular para extraer los parámetros
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	// Procesar los parámetros
	for _, match := range matches {
		flagName := match[1]                   // Nombre del flag (size, unit, path, type, fit, name)
		flagValue := strings.ToLower(match[2]) // Valor del flag (en minúsculas)
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size":
			// Convertir el valor a entero
			sizeVal, err := strconv.Atoi(flagValue)
			if err != nil {
				fmt.Println("Error: El valor de 'size' debe ser un número entero")
				return
			}
			size = sizeVal
		case "unit":
			unit = flagValue
		case "path":
			path = flagValue
		case "type":
			partType = flagValue
		case "fit":
			fit = flagValue
		case "name":
			name = flagValue
		default:
			fmt.Println("Error: Parámetro no reconocido:", flagName)
			return
		}
	}

	// Validaciones
	if size <= 0 {
		fmt.Println("Error: El valor de 'size' debe ser mayor que 0")
		return
	}

	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: La unidad debe ser 'b', 'k' o 'm'")
		return
	}

	if partType != "p" && partType != "e" && partType != "l" {
		fmt.Println("Error: El tipo de partición debe ser 'p', 'e' o 'l'")
		return
	}

	if fit != "bf" && fit != "ff" && fit != "wf" {
		fmt.Println("Error: El ajuste debe ser 'bf', 'ff' o 'wf'")
		return
	}

	if name == "" {
		fmt.Println("Error: El nombre de la partición es obligatorio")
		return
	}

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
=======
>>>>>>> origin/main
