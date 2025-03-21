package Partitions

import (
	"fmt"
)

// Definición de la estructura de una partición
type Partition struct {
<<<<<<< HEAD
	PartStatus      [1]byte  // 1 byte
	PartType        [1]byte  // 1 byte (P: Primaria, E: Extendida)
	PartFit         [1]byte  // 1 byte (B: Best, F: First, W: Worst)
=======
	PartStatus      byte     // 1 byte
	PartType        byte     // 1 byte (P: Primaria, E: Extendida)
	PartFit         byte     // 1 byte (B: Best, F: First, W: Worst)
>>>>>>> origin/main
	PartStart       int32    // 4 bytes
	PartSize        int32    // 4 bytes
	PartName        [16]byte // 16 bytes
	PartCorrelative int32    // 4 bytes
	PartID          [4]byte  // 4 bytes
}

// Definición de la estructura del MBR
type MBR struct {
	MbrTamanio       int32        // 4 bytes
	MbrFechaCreacion [19]byte     // 19 bytes (fecha en formato string)
	MbrDskSignature  int32        // 4 bytes (número único)
	DskFit           [1]byte      // 1 byte (B: Best, F: First, W: Worst)
	MbrPartitions    [4]Partition // 4 particiones de 36 bytes cada una
}

<<<<<<< HEAD
// Definición de la estructura de un EBR
type EBR struct {
	PartStatus byte     // 1 byte
	PartFit    byte     // 1 byte
	PartStart  int32    // 4 bytes
	PartSize   int32    // 4 bytes
	PartNext   int32    // 4 bytes
	PartName   [16]byte // 16 bytes
	PartMount  byte     // 1 byte
}

=======
>>>>>>> origin/main
func ImprimirMBR(data MBR) {
	fmt.Printf("MbrFechaCreacion: %s, DskFit: %s, Tamaño: %d, MbrDskSignature: %d\n",
		string(data.MbrFechaCreacion[:]),
		string(data.DskFit[:]),
		data.MbrTamanio,
		data.MbrDskSignature)
<<<<<<< HEAD
	for i := 0; i < 4; i++ {
		PrintPartition(data.MbrPartitions[i])
	}
}
func PrintPartition(data Partition) {
	fmt.Printf(
		//Usa fmt.Sprintf para formatear la información de la partición en un solo string.
		"Name: %s, type: %s, start: %d, size: %d, status: %s, id: %s\n",
		//string(data.Name[:]) convierte el array de bytes [16]byte a una cadena de texto.
		string(data.PartName[:]), string(data.PartType[:]), data.PartStart, data.PartSize,
		//data.Start y data.Size se imprimen como enteros (int32).
		string(data.PartStatus[:]), string(data.PartID[:]))
}

func PrintEBR(data EBR) {
	fmt.Printf("Name: %s, fit: %c, start: %d, size: %d, next: %d, mount: %c\n",
		string(data.PartName[:]),
		data.PartFit,
		data.PartStart,
		data.PartSize,
		data.PartNext,
		data.PartMount)
=======
>>>>>>> origin/main
}
