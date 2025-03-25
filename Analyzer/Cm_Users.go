package Analyzer

import (
	"Proyecto/User_Groups"
	"flag"
	"fmt"
	"os"
	"strings"
)

//═════════════════════   Administración de Usuarios y Grupo ═════════════════════

// Parametro para Iniciar Sesion
func fn_login(input string) {
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	user := fs.String("user", "", "Usuario")
	pass := fs.String("pass", "", "Contraseña")
	id := fs.String("id", "", "Id")

	fs.Parse(os.Args[1:])
	matches := paramRegex.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "user", "pass", "id":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("\033[1;31mError: Flag not found\033[0m")
		}
	}

	User_Groups.Login(*user, *pass, *id)
}

// Parametro para Cerrar Sesion
func fn_logout(_ string) {
	User_Groups.Logout()
}

// Parametro para Crear un Grupo
func fn_mkgrp(parametros string) {
	params := extraerParametros(parametros)
	if name, ok := params["name"]; ok {

		User_Groups.MKGRP(name)
	} else {
		fmt.Println("\033[1;31mError: Falta el parámetro -name\033[0m")
	}
}

// Parametro para Eliminar un Grupo
func fn_rmgrp(parametros string) {
	params := extraerParametros(parametros)
	if name, ok := params["name"]; ok {

		User_Groups.RMGRP(name)
	} else {
		fmt.Println("\033[1;31mError: Falta el parámetro -name\033[0m")
	}
}

// Parametro para Crear un Usuario
func fn_mkusr(parametros string) {
	paramMap := extraerParametros(parametros)

	// Validar que existan los parámetros necesarios
	user, userOK := paramMap["user"]
	pass, passOK := paramMap["pass"]
	grp, grpOK := paramMap["grp"]

	if !userOK || !passOK || !grpOK {
		fmt.Println("\033[1;31mError: Faltan parámetros obligatorios (-user, -pass, -grp)\033[0m")
		return
	}

	User_Groups.MKUSR(user, pass, grp)
}

// Parametro para Eliminar un Usuario
func fn_rmusr(parametros string) {
	paramMap := extraerParametros(parametros)

	// Validar que exista el parámetro obligatorio
	user, userOK := paramMap["user"]

	if !userOK {
		fmt.Println("\033[1;31mError: Falta el parámetro obligatorio (-user)\033[0m")
		return
	}

	User_Groups.RMUSR(user)
}

// Funcion para cambiar el grupo de un usuario
func fn_chgrp(parametros string) {
	paramMap := extraerParametros(parametros)

	// Validar que existan los parámetros necesarios
	user, userOK := paramMap["user"]
	grp, grpOK := paramMap["grp"]

	if !userOK || !grpOK {
		fmt.Println("\033[1;31mError: Faltan parámetros obligatorios (-user, -pass, -grp)\033[0m")
		return
	}

	User_Groups.CHGRP(user, grp)
}
