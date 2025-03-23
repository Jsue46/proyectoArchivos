package Analyzer

import (
	"Proyecto/Environment"
	"Proyecto/FSystem"
	"Proyecto/Partitions"
	"Proyecto/User_Groups"
	"Proyecto/Utils"
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
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
	case strings.Contains(comandos, "fdisk"):
		fn_fdisk(parametros)
	case strings.Contains(comandos, "mount"):
		fn_mount(parametros)
	case strings.Contains(comandos, "rmdisk"):
		fn_rmdisk(parametros)
	case strings.Contains(comandos, "mkfs"):
		fn_mkfs(parametros)
	case strings.Contains(comandos, "cat"):
		fn_cat(parametros)
	case strings.Contains(comandos, "login"):
		fn_login(parametros)
	case strings.Contains(comandos, "logout"):
		fn_logout(parametros)
	case strings.Contains(comandos, "salir"):
		fn_salir(parametros)
	case strings.Contains(comandos, "rep"):
		fn_rep(parametros)
	default:
		fmt.Println("\033[1;31mError: Comando inválido o no encontrado\033[0m")
	}
}

// fn_mkdisk procesa el comando mkdisk.
func fn_mkdisk(parametros string) {
	paramMap := extraerParametros(parametros)

	// Valores por defecto
	size, err := strconv.Atoi(paramMap["size"])
	if err != nil || size <= 0 {
		fmt.Println("\033[1;31mError: El valor de 'size' debe ser un número entero mayor que 0\033[0m")
		return
	}

	fit := strings.ToLower(paramMap["fit"])
	if fit == "" {
		fit = "ff" // Valor por defecto
	} else if fit != "bf" && fit != "ff" && fit != "wf" {
		fmt.Println("\033[1;31mError: El valor de 'fit' debe ser 'bf', 'ff' o 'wf'\033[0m")
		return
	}

	unit := strings.ToLower(paramMap["unit"])
	if unit == "" {
		unit = "m" // Valor por defecto
	} else if unit != "k" && unit != "m" {
		fmt.Println("\033[1;31mError: El valor de 'unit' debe ser 'k' o 'm'\033[0m")
		return
	}

	path := strings.ToLower(paramMap["path"]) // Convertir ruta a minúsculas
	if path == "" {
		fmt.Println("\033[1;31mError: El parámetro 'path' es obligatorio\033[0m")
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
		fmt.Println("\033[1;31mError: La unidad debe ser 'b', 'k' o 'm'\033[0m")
		return
	}

	fit := strings.ToLower(paramMap["fit"])
	if fit == "" {
		fit = "wf" // Valor por defecto: Worst Fit
	} else if fit != "bf" && fit != "ff" && fit != "wf" {
		fmt.Println("\033[1;31mError: El ajuste debe ser 'bf', 'ff' o 'wf'\033[0m")
		return
	}

	partType := strings.ToLower(paramMap["type"])
	if partType == "" {
		partType = "p" // Valor por defecto: Primaria
	} else if partType != "p" && partType != "e" && partType != "l" {
		fmt.Println("\033[1;31mError: El tipo de partición debe ser 'p', 'e' o 'l'\033[0m")
		return
	}

	// Validar parámetros obligatorios
	size, err := strconv.Atoi(paramMap["size"])
	if err != nil || size <= 0 {
		fmt.Println("\033[1;31mError: El valor de 'size' debe ser un número entero mayor que 0\033[0m")
		return
	}

	name := strings.ToLower(paramMap["name"]) // Convertir nombre a minúsculas
	if name == "" {
		fmt.Println("\033[1;31mError: El nombre de la partición es obligatorio\033[0m")
		return
	}

	path := strings.ToLower(paramMap["path"]) // Convertir ruta a minúsculas
	if path == "" {
		fmt.Println("\033[1;31mError: El parámetro 'path' es obligatorio\033[0m")
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
		fmt.Println("\033[1;31mError: Path y Name son obligatorios\033[0m")
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
		fmt.Println("\033[1;31mError: El parámetro 'path' es obligatorio\033[0m")
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

func fn_login(input string) {
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	user := fs.String("user", "", "Usuario")
	pass := fs.String("pass", "", "Contraseña")
	id := fs.String("id", "", "Id")

	// Parsearlas
	fs.Parse(os.Args[1:])
	// Match de flags en el input
	matches := paramRegex.FindAllStringSubmatch(input, -1)
	// Procesar el input
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

func fn_salir(_ string) {
	fmt.Println("\033[1;32mSaliendo del sistema...\033[0m")
	os.Exit(0)
}

// fn_logout procesa el comando logout.
func fn_logout(_ string) {
	fmt.Println("\033[1;36mProcesando comando LOGOUT...\033[0m")
	User_Groups.Logout()
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

func fn_rep(input string) {
	fs := flag.NewFlagSet("rep", flag.ExitOnError)
	name := fs.String("name", "", "Nombre del reporte a generar (mbr, disk, inode, block, bm_inode, bm_block, sb, file, ls)")
	path := fs.String("path", "", "Ruta donde se generará el reporte")
	id := fs.String("id", "", "ID de la partición")
	pathFileLs := fs.String("path_file_ls", "", "Nombre del archivo o carpeta para reportes file o ls") // Parámetro opcional

	// Parsear los parámetros de entrada
	matches := paramRegex.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.Trim(match[2], "\"")

		switch flagName {
		case "name", "path", "id", "path_file_ls":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag no encontrada:", flagName)
		}
	}

	// Verificar los parámetros obligatorios
	if *name == "" || *path == "" || *id == "" {
		fmt.Println("Error: 'name', 'path' y 'id' son parámetros obligatorios.")
		return
	}

	// Verificar si el disco está montado usando DiskManagement
	mounted := false
	var diskPath string
	for _, partitions := range Environment.GetMountedPartitions() {
		for _, partition := range partitions {
			if partition.MountID == *id {
				mounted = true
				diskPath = partition.MountPath
				break
			}
		}
	}

	if !mounted {
		fmt.Println("Error: La partición con ID", *id, "no está montada.")
		return
	}

	// Crear la carpeta si no existe
	reportsDir := filepath.Dir(*path)
	err := os.MkdirAll(reportsDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error al crear la carpeta:", reportsDir)
		return
	}

	// Generar el reporte según el tipo de reporte solicitado
	switch *name {
	case "mbr":
		// Abrir el archivo binario del disco montado
		file, err := Utils.AbrirArchivo(diskPath)
		if err != nil {
			fmt.Println("Error: No se pudo abrir el archivo en la ruta:", diskPath)
			return
		}
		defer file.Close()

		// Leer el objeto MBR desde el archivo binario
		var TempMBR Partitions.MBR
		if err := Utils.LeerArchivo(file, &TempMBR, 0); err != nil {
			fmt.Println("Error: No se pudo leer el MBR desde el archivo")
			return
		}

		// Leer y procesar los EBRs si hay particiones extendidas
		var ebrs []Partitions.EBR
		for i := 0; i < 4; i++ {
			if string(TempMBR.MbrPartitions[i].PartType[:]) == "e" { // Partición extendida
				fmt.Println("Partición extendida encontrada: ", string(TempMBR.MbrPartitions[i].PartName[:]))

				// El primer EBR está al inicio de la partición extendida
				ebrPosition := TempMBR.MbrPartitions[i].PartStart
				ebrCounter := 1

				// Leer todos los EBRs dentro de la partición extendida
				for ebrPosition != -1 {
					fmt.Printf("Leyendo EBR en posición: %d\n", ebrPosition)
					var tempEBR Partitions.EBR
					if err := Utils.LeerArchivo(file, &tempEBR, int64(ebrPosition)); err != nil {
						fmt.Println("Error: No se pudo leer el EBR desde el archivo")
						break
					}

					// Añadir el EBR a la lista
					ebrs = append(ebrs, tempEBR)
					fmt.Printf("EBR %d leído. Start: %d, Size: %d, Next: %d, Name: %s\n", ebrCounter, tempEBR.PartStart, tempEBR.PartSize, tempEBR.PartNext, string(tempEBR.PartName[:]))

					// Depuración: Mostrar el EBR leído
					Partitions.ImrprimirEBR(tempEBR)

					// Mover a la siguiente posición de EBR
					ebrPosition = tempEBR.PartNext
					ebrCounter++

					// Si no hay más EBRs, salir del bucle
					if ebrPosition == -1 {
						fmt.Println("No hay más EBRs en esta partición extendida.")
					}
				}
			}
		}

		// Generar el archivo .dot del MBR con EBRs
		reportPath := *path
		if err := Utils.GenerateMBRReport(TempMBR, ebrs, reportPath, file); err != nil {
			fmt.Println("Error al generar el reporte MBR:", err)
		} else {
			fmt.Println("Reporte MBR generado exitosamente en:", reportPath)

			// Renderizar el archivo .dot a .jpg usando Graphviz
			dotFile := strings.TrimSuffix(reportPath, filepath.Ext(reportPath)) + ".dot"
			outputJpg := reportPath
			cmd := exec.Command("dot", "-Tjpg", dotFile, "-o", outputJpg)
			err = cmd.Run()
			if err != nil {
				fmt.Println("Error al renderizar el archivo .dot a imagen:", err)
			} else {
				fmt.Println("Imagen generada exitosamente en:", outputJpg)
			}
		}

	//CASE PARA EL REPORTE DISK
	case "disk":
		// Abrir el archivo binario del disco montado
		file, err := Utils.AbrirArchivo(diskPath)
		if err != nil {
			fmt.Println("Error: No se pudo abrir el archivo en la ruta:", diskPath)
			return
		}
		defer file.Close()

		// Leer el objeto MBR desde el archivo binario
		var TempMBR Partitions.MBR
		if err := Utils.LeerArchivo(file, &TempMBR, 0); err != nil {
			fmt.Println("Error: No se pudo leer el MBR desde el archivo")
			return
		}

		// Leer y procesar los EBRs si hay particiones extendidas
		var ebrs []Partitions.EBR
		for i := 0; i < 4; i++ {
			if string(TempMBR.MbrPartitions[i].PartType[:]) == "e" { // Partición extendida
				ebrPosition := TempMBR.MbrPartitions[i].PartStart
				for ebrPosition != -1 {
					var tempEBR Partitions.EBR
					if err := Utils.LeerArchivo(file, &tempEBR, int64(ebrPosition)); err != nil {
						break
					}
					ebrs = append(ebrs, tempEBR)
					ebrPosition = tempEBR.PartNext
				}
			}
		}

		// Calcular el tamaño total del disco
		totalDiskSize := TempMBR.MbrTamanio

		// Generar el archivo .dot del DISK
		reportPath := *path
		if err := Utils.GenerateDiskReport(TempMBR, ebrs, reportPath, file, totalDiskSize); err != nil {
			fmt.Println("Error al generar el reporte DISK:", err)
		} else {
			fmt.Println("Reporte DISK generado exitosamente en:", reportPath)

			// Renderizar el archivo .dot a .jpg usando Graphviz
			dotFile := strings.TrimSuffix(reportPath, filepath.Ext(reportPath)) + ".dot"
			outputJpg := reportPath
			cmd := exec.Command("dot", "-Tjpg", dotFile, "-o", outputJpg)
			err = cmd.Run()
			if err != nil {
				fmt.Println("Error al renderizar el archivo .dot a imagen:", err)
			} else {
				fmt.Println("Imagen generada exitosamente en:", outputJpg)
			}
		}

	case "file", "ls":
		// Para los reportes "file" y "ls", pathFileLs es obligatorio
		if *pathFileLs == "" {
			fmt.Println("Error: 'path_file_ls' es obligatorio para los reportes 'file' y 'ls'.")
			return
		}

		// Lógica para generar los reportes de tipo 'file' y 'ls'
		fmt.Println("Generando reporte", *name, "con archivo/carpeta:", *pathFileLs)
		// Aquí iría la lógica adicional para generar estos reportes

	default:
		fmt.Println("Error: Tipo de reporte no válido.")
	}
}
