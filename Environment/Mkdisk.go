/*package Environment

import (
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func Mkdisk(size int, fit string, unit string, path string) {
	fmt.Println("╔═════════════════════ INICIO MKDISK  ═════════════════════════╗")
	fmt.Printf("  Size: %d\n  Fit: %s\n  Unit: %s\n  Path: %s\n", size, fit, unit, path)

	// Valores por defecto para parámetros opcionales
	if fit == "" {
		fit = "ff"
	}
	if unit == "" {
		unit = "m"
	}

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

	// Crear directorios si no existen
	if err := os.MkdirAll(path[:strings.LastIndex(path, "/")], os.ModePerm); err != nil {
		fmt.Println("Error al crear directorios:", err.Error())
		return
	}

	// Crear archivo
	if err := Utils.CrearArchivo(path); err != nil {
		fmt.Println("Error al crear archivo:", err.Error())
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
		fmt.Println("Error al abrir archivo:", err.Error())
		return
	}
	defer file.Close()

	// Escribir ceros en el archivo
	zeroBlock := make([]byte, sizeInBytes)
	if _, err := file.Write(zeroBlock); err != nil {
		fmt.Println("Error al escribir en el archivo:", err.Error())
		return
	}

	// Crear MBR
	var nuevoMBR Partitions.MBR
	nuevoMBR.MbrTamanio = int32(sizeInBytes)
	nuevoMBR.MbrDskSignature = rand.Int31()
	copy(nuevoMBR.DskFit[:], fit)

	formattedDate := time.Now().Format("15:04 - 02/01/2006")
	copy(nuevoMBR.MbrFechaCreacion[:], formattedDate)

	// Escribir el MBR en el archivo
	if err := Utils.EscribirArchivo(file, nuevoMBR, 0); err != nil {
		fmt.Println("Error al escribir el MBR:", err.Error())
		return
	}

	// Leer el MBR para verificar que se escribió correctamente
	var tempMBR Partitions.MBR
	if err := Utils.LeerArchivo(file, &tempMBR, 0); err != nil {
		fmt.Println("Error al leer el MBR:", err.Error())
		return
	}

	// Imprimir el MBR leído
	Partitions.ImprimirMBR(tempMBR)
	fmt.Println("╚═════════════════════   FIN MKDISK   ═════════════════════════╝")
}
*/

// Environment/environment.go
package Environment

import (
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func Mkdisk(size int, fit string, unit string, path string) {
	ClearLogs() // Limpiar el buffer antes de empezar

	Log("======INICIO MKDISK======")
	Log(fmt.Sprintf("Size: %d\nFit: %s\nUnit: %s\nPath: %s", size, fit, unit, path))

	// Valores por defecto para parámetros opcionales
	if fit == "" {
		fit = "ff"
	}
	if unit == "" {
		unit = "m"
	}

	// Validaciones
	if fit != "bf" && fit != "wf" && fit != "ff" {
		Log("Error: Fit debe ser 'bf', 'wf' o 'ff'")
		return
	}
	if size <= 0 {
		Log("Error: Size debe ser mayor a 0")
		return
	}
	if unit != "k" && unit != "m" {
		Log("Error: Las unidades válidas son 'k' o 'm'")
		return
	}

	// Crear directorios si no existen
	if err := os.MkdirAll(path[:strings.LastIndex(path, "/")], os.ModePerm); err != nil {
		Log("Error al crear directorios: " + err.Error())
		return
	}

	// Crear archivo
	if err := Utils.CrearArchivo(path); err != nil {
		Log("Error al crear archivo: " + err.Error())
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
		Log("Error al abrir archivo: " + err.Error())
		return
	}
	defer file.Close()

	// Escribir ceros en el archivo (creando el espacio del disco)
	zeroBlock := make([]byte, sizeInBytes)
	if _, err := file.Write(zeroBlock); err != nil {
		Log("Error al escribir en el archivo: " + err.Error())
		return
	}

	// Crear MBR
	var nuevoMBR Partitions.MBR
	nuevoMBR.MbrTamanio = int32(sizeInBytes)
	nuevoMBR.MbrDskSignature = rand.Int31()
	copy(nuevoMBR.DskFit[:], fit)

	formattedDate := time.Now().Format("02/01/2006 15:04")
	copy(nuevoMBR.MbrFechaCreacion[:], formattedDate)

	// Escribir el MBR en el archivo
	if err := Utils.EscribirArchivo(file, nuevoMBR, 0); err != nil {
		Log("Error al escribir el MBR: " + err.Error())
		return
	}

	// Leer el MBR para verificar que se escribió correctamente
	var tempMBR Partitions.MBR
	if err := Utils.LeerArchivo(file, &tempMBR, 0); err != nil {
		Log("Error al leer el MBR: " + err.Error())
		return
	}

	// Imprimir el MBR leído
	Partitions.ImprimirMBR(tempMBR)
	Log("Disco creado exitosamente")
	Log("======FIN MKDISK======")
}
