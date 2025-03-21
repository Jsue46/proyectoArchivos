package Environment

import (
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"bytes"
	"fmt"
	"strings"
)

// Función para montar una Partición
func Mount(path string, name string) {
	file, err := Utils.AbrirArchivo(path)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo en la ruta:", path)
		return
	}
	defer file.Close()

	var TempMBR Partitions.MBR
	if err := Utils.LeerArchivo(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR desde el archivo")
		return
	}

	fmt.Printf("Buscando partición con nombre: '%s'\n", name)

	partitionFound := false
	var partition Partitions.Partition
	var partitionIndex int

	nameBytes := [16]byte{}
	copy(nameBytes[:], []byte(name))

	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].PartType[0] == 'p' && bytes.Equal(TempMBR.MbrPartitions[i].PartName[:], nameBytes[:]) {
			partition = TempMBR.MbrPartitions[i]
			partitionIndex = i
			partitionFound = true
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: Partición no encontrada o no es una partición primaria")
		return
	}

	if partition.PartStatus[0] == '1' {
		fmt.Println("Error: La partición ya está montada")
		return
	}

	diskID := generateDiskID(path)
	mountedPartitionsInDisk := mountedPartitions[diskID]
	var letter byte

	if len(mountedPartitionsInDisk) == 0 {
		if len(mountedPartitions) == 0 {
			letter = 'A'
		} else {
			lastDiskID := getLastDiskID()
			lastLetter := mountedPartitions[lastDiskID][0].MountID[len(mountedPartitions[lastDiskID][0].MountID)-1]
			letter = lastLetter + 1
		}
	} else {
		letter = mountedPartitionsInDisk[0].MountID[len(mountedPartitionsInDisk[0].MountID)-1]
	}

	carnet := "202203806"
	lastTwoDigits := carnet[len(carnet)-2:]
	partitionID := fmt.Sprintf("%s%d%c", lastTwoDigits, partitionIndex+1, letter)

	partition.PartStatus[0] = '1'
	copy(partition.PartID[:], partitionID)

	TempMBR.MbrPartitions[partitionIndex] = partition

	mountedPartitions[diskID] = append(mountedPartitions[diskID], MountedPartition{
		MountPath:   path,
		MountName:   name,
		MountID:     partitionID,
		MountStatus: '1',
	})

	if err := Utils.EscribirArchivo(file, TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo sobrescribir el MBR en el archivo")
		return
	}

	fmt.Printf("Partición montada con ID: %s\n", partitionID)
}

func getLastDiskID() string {
	var lastDiskID string
	for diskID := range mountedPartitions {
		lastDiskID = diskID
	}
	return lastDiskID
}

func generateDiskID(path string) string {
	return strings.ToLower(path)
}
