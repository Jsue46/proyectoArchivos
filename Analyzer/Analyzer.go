package Analyzer

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var paramRegex = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

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
}
