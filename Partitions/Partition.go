package Partitions

import (
	"fmt"
)

// Definición de la estructura de una partición
type Partition struct {
	PartStatus      [1]byte  // 1 byte
	PartType        [1]byte  // 1 byte (P: Primaria, E: Extendida)
	PartFit         [1]byte  // 1 byte (B: Best, F: First, W: Worst)
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

// Función para imprimir un MBR
func ImprimirMBR(data MBR) {
	fmt.Printf("MbrFechaCreacion: %s, DskFit: %s, Tamaño: %d, MbrDskSignature: %d\n",
		string(data.MbrFechaCreacion[:]),
		string(data.DskFit[:]),
		data.MbrTamanio,
		data.MbrDskSignature)
	for i := 0; i < 4; i++ {
		ImprimirPartition(data.MbrPartitions[i])
	}
}

// Función para imprimir una partición
func ImprimirPartition(data Partition) {
	fmt.Printf(
		"Name: %s, type: %s, start: %d, size: %d, status: %s, id: %s\n",
		string(data.PartName[:]), string(data.PartType[:]), data.PartStart, data.PartSize,
		string(data.PartStatus[:]), string(data.PartID[:]))
}

// Función para imprimir un EBR
func ImrprimirEBR(data EBR) {
	fmt.Printf("Name: %s, fit: %c, start: %d, size: %d, next: %d, mount: %c\n",
		string(data.PartName[:]),
		data.PartFit,
		data.PartStart,
		data.PartSize,
		data.PartNext,
		data.PartMount)
}
