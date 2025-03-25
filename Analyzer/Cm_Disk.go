package Analyzer

import (
	"Proyecto/Environment"
	"fmt"
	"strconv"
	"strings"
)

// fn_mkdisk procesa el comando mkdisk.
func fn_mkdisk(parametros string) {
	paramMap := extraerParametros(parametros)

	// Validar y procesar parámetros
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

	path := strings.ToLower(paramMap["path"])
	if path == "" {
		fmt.Println("Error: El parámetro 'path' es obligatorio")
		return
	}

	// Llamar a la función Mkdisk con los parámetros procesados
	Environment.Mkdisk(size, fit, unit, path)
}

// fn_rmdisk procesa el comando rmdisk.
func fn_rmdisk(parametros string) {
	paramMap := extraerParametros(parametros)

	path := strings.ToLower(paramMap["path"])
	if path == "" {
		fmt.Println("Error: El parámetro 'path' es obligatorio")
		return
	}

	// Llamar a la función Rmdisk del paquete Environment
	Environment.Rmdisk(path)
}

// fn_fdisk procesa el comando fdisk.
func fn_fdisk(parametros string) {
	paramMap := extraerParametros(parametros)

	// Validar y procesar parámetros
	unit := strings.ToLower(paramMap["unit"])
	if unit == "" {
		unit = "k" // Valor por defecto
	} else if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: La unidad debe ser 'b', 'k' o 'm'")
		return
	}

	fit := strings.ToLower(paramMap["fit"])
	if fit == "" {
		fit = "wf" // Valor por defecto
	} else if fit != "bf" && fit != "ff" && fit != "wf" {
		fmt.Println("Error: El ajuste debe ser 'bf', 'ff' o 'wf'")
		return
	}

	partType := strings.ToLower(paramMap["type"])
	if partType == "" {
		partType = "p" // Valor por defecto
	} else if partType != "p" && partType != "e" && partType != "l" {
		fmt.Println("Error: El tipo de partición debe ser 'p', 'e' o 'l'")
		return
	}

	size, err := strconv.Atoi(paramMap["size"])
	if err != nil || size <= 0 {
		fmt.Println("Error: El valor de 'size' debe ser un número entero mayor que 0")
		return
	}

	name := strings.ToLower(paramMap["name"])
	if name == "" {
		fmt.Println("Error: El nombre de la partición es obligatorio")
		return
	}

	path := strings.ToLower(paramMap["path"])
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

	path := strings.ToLower(paramMap["path"])
	name := strings.ToLower(paramMap["name"])

	if path == "" || name == "" {
		fmt.Println("Error: Path y Name son obligatorios")
		return
	}

	// Llamar a la función Mount con los parámetros procesados
	Environment.Mount(path, name)
}

// fn_mounted procesa el comando mounted.
func fn_mounted(_ string) {
	mountedPartitions := Environment.GetMountedPartitions()

	// Verificar si hay particiones montadas
	if len(mountedPartitions) == 0 {
		fmt.Println("No hay particiones montadas.")
		return
	}

	// Mostrar los IDs de las particiones montadas
	var ids []string
	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			ids = append(ids, partition.MountID)
		}
	}
	fmt.Println("╔═════════════════ PARTICIONES MONTADAS ═══════════════════════╗")
	fmt.Println(" ", strings.Join(ids, ", "))
	fmt.Println("╚═════════════════════ FIN PARTCIONES ═════════════════════════╝")
}
