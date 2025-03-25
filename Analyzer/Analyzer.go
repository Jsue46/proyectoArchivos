package Analyzer

import (
<<<<<<< Updated upstream
	"bufio"
=======
	"Proyecto/Environment"
>>>>>>> Stashed changes
	"fmt"
	"regexp"
	"strings"
)

<<<<<<< Updated upstream
var paramRegex = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
=======
func AnalizadorEntrada(input string) string {
	// Función para determinar si un carácter es un delimitador
	isDelimiter := func(r rune) bool {
		return r == ';' || r == '\n'
	}

	// Dividir la entrada en comandos usando los delimitadores ";" y "\n"
	comandos := strings.FieldsFunc(input, isDelimiter)

	// Variable para almacenar los resultados de cada comando
	var resultados []string

	// Procesar cada comando
	for _, comando := range comandos {
		// Eliminar espacios en blanco al inicio y final del comando
		comando = strings.TrimSpace(comando)

		// Si el comando no está vacío, procesarlo
		if comando != "" {
			// Obtener el comando y los parámetros
			command, params := getEntrada(comando)

			// Procesar el comando y almacenar el resultado
			resultado := AnalizadorComandos(command, params)
			resultados = append(resultados, resultado)
		}
	}

	// Unir todos los resultados en un solo string separados por saltos de línea
	return strings.Join(resultados, "\n")
}
>>>>>>> Stashed changes

// getEntrada recibe una cadena de texto y devuelve el comando y los parámetros.
func getEntrada(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		comandos := strings.ToLower(parts[0])
		parametros := strings.Join(parts[1:], " ")
		return comandos, parametros
	}
	return "", input
}

<<<<<<< Updated upstream
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
		fmt.Print("\033[1;32mIngresar un Comando: \033[0m")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		comandos, parametros := getEntrada(input)
		fmt.Printf("\033[1;34mComando: \033[1;35m%s\033[0m - \033[1;34mParametro: \033[1;35m%s\033[0m\n", comandos, parametros)

		AnalizadorComandos(comandos, parametros)
	}
}

// AnalizadorComandos redirige el comando a la función correspondiente.
func AnalizadorComandos(comandos string, parametros string) {
	switch {
	case strings.Contains(comandos, "mkdisk"):
		fn_mkdisk(parametros)
	case strings.Contains(comandos, "rmdisk"):
		fn_rmdisk(parametros)
	case strings.Contains(comandos, "fdisk"):
		fn_fdisk(parametros)
	case strings.Contains(comandos, "mounted"):
		fn_mounted(parametros)
	case strings.Contains(comandos, "mount"):
		fn_mount(parametros)
	case strings.Contains(comandos, "mkfs"):
		fn_mkfs(parametros)
	case strings.Contains(comandos, "cat"): // No Implementado
		fn_cat(parametros)
	case strings.Contains(comandos, "login"):
		fn_login(parametros)
	case strings.Contains(comandos, "logout"):
		fn_logout(parametros)
	case strings.Contains(comandos, "mkgrp"):
		fn_mkgrp(parametros)
	case strings.Contains(comandos, "rmgrp"):
		fn_rmgrp(parametros)
	case strings.Contains(comandos, "mkusr"):
		fn_mkusr(parametros)
	case strings.Contains(comandos, "rmusr"):
		fn_rmusr(parametros)
	case strings.Contains(comandos, "chgrp"):
		fn_chgrp(parametros)
	case strings.Contains(comandos, "mkfile"): //No Implementado
		fn_mkfile(parametros)
	case strings.Contains(comandos, "mkdir"): //No Implementado
		fn_mkdir(parametros)
	case strings.Contains(comandos, "rep"): //Aun Faltan Agregar.
		fn_rep(parametros)
	case strings.Contains(comandos, "salir"):
		fn_salir(parametros)
	default:
		fmt.Println("\033[1;31mError: Comando inválido o no encontrado\033[0m")
	}
}

func fn_salir(_ string) {
	fmt.Println("\033[1;32mSaliendo del sistema...\033[0m")
	os.Exit(0)
}

// Archvivos Y permisos
func fn_mkfile(_ string) {
	fmt.Println("\033[1;32mCreando en el Sistema...\033[0m")
}

func fn_mkdir(_ string) {
	fmt.Println("\033[1;32mCreando en el Sistema..\033[0m")
=======
// AnalizadorComandos procesa el comando y devuelve una cadena con el resultado.
func AnalizadorComandos(command string, params string) string {
	command = strings.ToLower(command)
	switch command {
	case "mkdisk":
		return fn_mkdisk(params)
	case "rmdisk":
		return fn_rmdisk(params)
	case "fdisk":
		return fn_fdisk(params)
	case "mount":
		return fn_mount(params)
	case "mounted":
		return fn_mounted(params)
	default:
		return "Error: Comando inválido"
	}
}

// fn_mkdisk procesa el comando mkdisk.
func fn_mkdisk(params string) string {
	Environment.ClearLogs()
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
				return "Error: El valor de 'size' debe ser un número entero"
			}
			size = sizeVal
		case "fit":
			fit = flagValue
		case "unit":
			unit = flagValue
		case "path":
			path = flagValue
		default:
			return "Error: Parámetro no reconocido: " + flagName
		}
	}

	// Validaciones
	if size <= 0 {
		return "Error: El valor de 'size' debe ser mayor que 0"
	}

	if fit != "bf" && fit != "ff" && fit != "wf" {
		return "Error: El valor de 'fit' debe ser 'bf', 'ff' o 'wf'"
	}

	if unit != "k" && unit != "m" {
		return "Error: El valor de 'unit' debe ser 'k' o 'm'"
	}

	if path == "" {
		return "Error: El parámetro 'path' es obligatorio"
	}

	// Llamar a la función Mkdisk
	Environment.Mkdisk(size, fit, unit, path)

	// Devolver los logs capturados
	return Environment.GetLogs()
}

// fn_rmdisk procesa el comando rmdisk.
func fn_rmdisk(params string) string {
	// Expresión regular para extraer los parámetros
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(params, -1)

	// Variable para almacenar la ruta
	var path string

	// Procesar los parámetros
	for _, match := range matches {
		flagName := match[1]                   // Nombre del flag (path)
		flagValue := strings.ToLower(match[2]) // Valor del flag (en minúsculas)
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			path = flagValue
		default:
			return "Error: Parámetro no reconocido: " + flagName
		}
	}

	// Validaciones
	if path == "" {
		return "Error: El parámetro 'path' es obligatorio"
	}

	// Llamar a la función Rmdisk
	Environment.Rmdisk(path)

	// Devolver un mensaje personalizado
	return "El disco '" + path + "' fue eliminado correctamente."
}

// fn_fdisk procesa el comando fdisk.
func fn_fdisk(params string) string {
	Environment.ClearLogs()
	// Variables para almacenar los parámetros
	var size int
	var path, name string
	var unit, type_, fit string

	// Valores predeterminados para parámetros opcionales
	unit = "k"  // Valor predeterminado para unit
	type_ = "p" // Valor predeterminado para type
	fit = "wf"  // Valor predeterminado para fit

	// Expresión regular para extraer los parámetros
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(params, -1)

	// Procesar los parámetros
	for _, match := range matches {
		flagName := match[1]                   // Nombre del flag (size, path, name, unit, type, fit)
		flagValue := strings.ToLower(match[2]) // Valor del flag (en minúsculas)
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size":
			// Convertir el valor a entero
			sizeVal, err := strconv.Atoi(flagValue)
			if err != nil {
				return "Error: El valor de 'size' debe ser un número entero"
			}
			size = sizeVal
		case "path":
			path = flagValue
		case "name":
			name = flagValue
		case "unit":
			unit = flagValue
		case "type":
			type_ = flagValue
		case "fit":
			fit = flagValue
		default:
			return "Error: Parámetro no reconocido: " + flagName
		}
	}

	// Validaciones obligatorias
	if size <= 0 {
		return "Error: El valor de 'size' debe ser mayor que 0"
	}

	if name == "" {
		return "Error: El parámetro 'name' es obligatorio"
	}

	if path == "" {
		return "Error: El parámetro 'path' es obligatorio"
	}

	// Validaciones condicionales (solo si el parámetro está presente)
	if unit != "b" && unit != "k" && unit != "m" {
		return "Error: El valor de 'unit' debe ser 'b', 'k' o 'm'"
	}

	if type_ != "p" && type_ != "e" && type_ != "l" {
		return "Error: El valor de 'type' debe ser 'p', 'e' o 'l'"
	}

	if fit != "bf" && fit != "ff" && fit != "wf" {
		return "Error: El valor de 'fit' debe ser 'bf', 'ff' o 'wf'"
	}

	// Llamar a la función Fdisk
	Environment.Fdisk(size, path, name, unit, type_, fit)

	// Devolver los logs capturados
	response := fmt.Sprintf(`
======INICIO FDISK======
Size: %d
Path: %s
Name: %s
Unit: %s
Type: %s
Fit: %s
%s
======FIN FDISK======`,
		size, path, name, unit, type_, fit,
		Environment.GetLogs()) // Incorpora los logs técnicos

	return response
}
func fn_mount(params string) string {
	// Parsear parámetros
	path, name, err := parseMountParams(params)
	if err != nil {
		return err.Error()
	}

	// Ejecutar comando MOUNT
	return Environment.Mount(path, name)
}

func fn_mounted(params string) string {
	// El comando mounted no acepta parámetros
	if strings.TrimSpace(params) != "" {
		return "Error: El comando 'mounted' no acepta parámetros"
	}

	// Ejecutar comando MOUNTED
	return Environment.ListMountedPartitions()
}

// parseMountParams extrae path y name de los parámetros
func parseMountParams(params string) (string, string, error) {
	var path, name string
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		paramName := strings.ToLower(match[1])
		paramValue := strings.Trim(match[2], "\"")

		switch paramName {
		case "path":
			path = paramValue
		case "name":
			name = paramValue
		default:
			return "", "", fmt.Errorf("parámetro no reconocido '%s'", paramName)
		}
	}

	if path == "" {
		return "", "", fmt.Errorf("falta parámetro obligatorio 'path'")
	}
	if name == "" {
		return "", "", fmt.Errorf("falta parámetro obligatorio 'name'")
	}

	return path, name, nil
>>>>>>> Stashed changes
}
