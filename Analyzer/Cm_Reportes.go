package Analyzer

import (
	"Proyecto/Environment"
	"Proyecto/Partitions"
	"Proyecto/Reportes"
	"Proyecto/Utils"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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
		if err := Reportes.GenerateMBRReport(TempMBR, ebrs, reportPath, file); err != nil {
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
		if err := Reportes.GenerateDiskReport(TempMBR, ebrs, reportPath, file, totalDiskSize); err != nil {
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
