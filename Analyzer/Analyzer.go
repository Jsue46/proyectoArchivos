package Analyzer

import (
	"Proyecto/Environment"
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// input := "mkdisk -size=3000 -unit=K -fit=BF -path=/home/daaniieel/Disco/disk1.mia"
// input := "fdisk -size=500 -unit=K -path=/home/daaniieel/Disco/disk1.bin -type=P -fit=b -name=Particion1"
// input := "rmdisk -path="/home/daaniieel/Disco/disk1.mia""

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

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
	if strings.Contains(command, "mkdisk") {
		fn_mkdisk(params)
	} else if strings.Contains(command, "fdisk") {
		fn_fdisk(params)
	} else if strings.Contains(command, "mount") {
		fn_mount(params)
	} else if strings.Contains(command, "rmdisk") {
		fn_rmdisk(params)
	} else if strings.Contains(command, "rep") {
		fmt.Print("COMANDO REP")
	} else {
		fmt.Println("Error: Comando inválido o no encontrado")
	}
}

func fn_mkdisk(params string) {
	var size int
	var fit, unit, path string
	// Valores por defecto
	fit = "ff"
	unit = "m"
	// Expresión regular para extraer los parámetros
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(params, -1)

	// Procesar los parámetros
	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.ToLower(match[2])
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

// Función para crear una partición
func fn_fdisk(input string) {
	var size int
	var unit, path, partType, fit, name string

	// Expresión regular para extraer los parámetros
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	// Procesar los parámetros
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])
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

// Función para montar una partición
func fn_mount(params string) {
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	path := fs.String("path", "", "Ruta")
	name := fs.String("name", "", "Nombre de la partición")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])
		flagValue = strings.Trim(flagValue, "\"")
		fs.Set(flagName, flagValue)
	}

	if *path == "" || *name == "" {
		fmt.Println("Error: Path y Name son obligatorios")
		return
	}
	lowercaseName := strings.ToLower(*name)
	// Llamar a la función Mount con los parámetros procesados
	Environment.Mount(*path, lowercaseName)
}

// Función para eliminar un disco
func fn_rmdisk(params string) {
	// Definir flag
	fs := flag.NewFlagSet("rmdisk", flag.ExitOnError)
	path := fs.String("path", "", "Ruta del disco a eliminar")
	// Expresión regular para extraer los parámetros
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(params, -1)

	// Procesar los parámetros
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.Trim(match[2], "\"")

		switch flagName {
		case "path":
			*path = flagValue
		default:
			fmt.Println("Error: Parámetro no reconocido:", flagName)
			return
		}
	}
	// Validaciones
	if *path == "" {
		fmt.Println("Error: El parámetro 'path' es obligatorio")
		return
	}

	// Llamar a la función Rmdisk del paquete Environment
	Environment.Rmdisk(*path)
}
