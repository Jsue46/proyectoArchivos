package Partitions

import (
	"fmt"
)

// Definición de la estructura de una partición
type Partition struct {
	PartStatus      byte     // 1 byte
	PartType        byte     // 1 byte (P: Primaria, E: Extendida)
	PartFit         byte     // 1 byte (B: Best, F: First, W: Worst)
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

func ImprimirMBR(data MBR) {
	fmt.Printf("MbrFechaCreacion: %s, DskFit: %s, Tamaño: %d, MbrDskSignature: %d\n",
		string(data.MbrFechaCreacion[:]),
		string(data.DskFit[:]),
		data.MbrTamanio,
		data.MbrDskSignature)
}
