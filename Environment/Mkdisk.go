package Environment

import (
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"fmt"
<<<<<<< HEAD
	"math/rand"
	"time"
)

func Mkdisk(size int, fit string, unit string, path string) {
	fmt.Println("======INICIO MKDISK======")
=======
)

func Mkdisk(size int, fit string, unit string, path string) {
	fmt.Println("\033[1;36m╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                      \033[1;33mINICION  MKDISK\033[1;36m                              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝\033[0m")
>>>>>>> origin/main
	fmt.Printf("Size: %d\nFit: %s\nUnit: %s\nPath: %s\n", size, fit, unit, path)

	// Validaciones
	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Println("Error: Fit debe ser 'bf', 'wf' o 'ff'")
		return
	}
	if size <= 0 {
		fmt.Println("Error: Size debe ser mayor a 0")
		return
	}
	if unit != "k" && unit != "m" {
		fmt.Println("Error: Las unidades válidas son 'k' o 'm'")
		return
	}

	// Crear archivo
	if err := Utils.CrearArchivo(path); err != nil {
		fmt.Println("Error al crear archivo:", err)
		return
	}

	// Convertir tamaño a bytes
	sizeInBytes := size * 1024
	if unit == "m" {
		sizeInBytes *= 1024
	}

	// Abrir archivo
	file, err := Utils.AbrirArchivo(path)
	if err != nil {
		fmt.Println("Error al abrir archivo:", err)
		return
	}
	defer file.Close() // Asegura el cierre del archivo al salir de la función

	// Escribir ceros en un solo bloque en lugar de un bucle
	zeroBlock := make([]byte, sizeInBytes) // Crea un slice de bytes lleno de ceros
	if _, err := file.Write(zeroBlock); err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
		return
	}

<<<<<<< HEAD
	// Crear MBR
	var nuevoMBR Partitions.MBR
	nuevoMBR.MbrTamanio = int32(sizeInBytes)
	nuevoMBR.MbrTamanio = rand.Int31()
	copy(nuevoMBR.DskFit[:], fit)

	formattedDate := time.Now().Format("02/01/2006 15:04")
	copy(nuevoMBR.MbrFechaCreacion[:], formattedDate)

	// Escribir el MBR en el archivo
	if err := Utils.EscribirArchivo(file, nuevoMBR, 0); err != nil {
		fmt.Println("Error al escribir el MBR:", err)
		return
	}

=======
	// Crear y escribir el MBR
	if err := Partitions.CrearMBR(file, int32(sizeInBytes), fit[0]); err != nil {
		fmt.Println("Error al crear el MBR:", err)
		return
	}
>>>>>>> origin/main
	// Leer el MBR para verificar que se escribió correctamente
	var tempMBR Partitions.MBR
	if err := Utils.LeerArchivo(file, &tempMBR, 0); err != nil {
		fmt.Println("Error al leer el MBR:", err)
		return
	}

	// Imprimir el MBR leído
	Partitions.ImprimirMBR(tempMBR)

<<<<<<< HEAD
	fmt.Println("======FIN MKDISK======")
=======
	fmt.Println("\033[1;36m╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                      \033[1;FIN  MKDISK\033[1;36m                              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝\033[0m")
>>>>>>> origin/main
}
