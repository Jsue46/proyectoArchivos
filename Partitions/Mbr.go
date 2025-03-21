package Partitions

import (
	"Proyecto/Utils"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// CrearMBR inicializa y escribe un MBR en un archivo binario.
<<<<<<< HEAD
=======
// CrearMBR inicializa y escribe un MBR en un archivo binario.
>>>>>>> origin/main
func CrearMBR(file *os.File, size int32, fit byte) error {
	// Crear MBR
	var nuevoMBR MBR
	nuevoMBR.MbrTamanio = size
	nuevoMBR.MbrDskSignature = rand.Int31() // Número único aleatorio
	nuevoMBR.DskFit = [1]byte{fit}          // Tipo de ajuste

<<<<<<< HEAD
	// Obtener la fecha actual en el formato "día/mes/año hora:min"
	formatoFecha := time.Now().Format("02/01/2006 15:04")
=======
	// Obtener la fecha actual
	formatoFecha := time.Now().Format("2000-01-01 15:04:05")
>>>>>>> origin/main
	copy(nuevoMBR.MbrFechaCreacion[:], formatoFecha)

	// Inicializar el arreglo de particiones
	for i := 0; i < 4; i++ {
		nuevoMBR.MbrPartitions[i] = Partition{
<<<<<<< HEAD
			PartStatus:      [1]byte{0},   // 0: No montada
			PartType:        [1]byte{'P'}, // P: Primaria (por defecto)
			PartFit:         [1]byte{fit}, // Tipo de ajuste
			PartStart:       -1,           // -1: No iniciada
			PartSize:        0,            // Tamaño inicial 0
			PartName:        [16]byte{},   // Nombre vacío
			PartCorrelative: -1,           // -1: No montada
			PartID:          [4]byte{},    // ID vacío
=======
			PartStatus:      0,          // 0: No montada
			PartType:        'P',        // P: Primaria (por defecto)
			PartFit:         fit,        // Tipo de ajuste
			PartStart:       -1,         // -1: No iniciada
			PartSize:        0,          // Tamaño inicial 0
			PartName:        [16]byte{}, // Nombre vacío
			PartCorrelative: -1,         // -1: No montada
			PartID:          [4]byte{},  // ID vacío
>>>>>>> origin/main
		}
	}

	// Escribir el MBR en el archivo
	if err := Utils.EscribirArchivo(file, nuevoMBR, 0); err != nil {
		return fmt.Errorf("error al escribir el MBR: %v", err)
	}

	return nil
}
