/*
package Environment

import (

	"Proyecto/Partitions"
	"Proyecto/Utils"
	"bytes"
	"fmt"

)

//var mountPartitions = make(map[string][]MountedPartition)

	func getLastDiskID() string {
		var lastID string
		for diskID := range mountedPartitions {
			lastID = diskID
		}
		return lastID
	}

	func generateDiskID(path string) string {
		return fmt.Sprintf("disk_%s", path)
	}

	func Mount(path string, name string) string {
		file, err := Utils.AbrirArchivo(path)
		if err != nil {
			return fmt.Sprintf("Error: No se pudo abrir el archivo en la ruta: %s", path)
		}
		defer file.Close()

		var TempMBR Partitions.MBR
		if err := Utils.LeerArchivo(file, &TempMBR, 0); err != nil {
			return "Error: No se pudo leer el MBR desde el archivo"
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
			return "Error: Partición no encontrada o no es una partición primaria"
		}

		if partition.PartStatus[0] == '1' {
			return "Error: La partición ya está montada"
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
			return "Error: No se pudo sobrescribir el MBR en el archivo"
		}

		return fmt.Sprintf("Partición montada con ID: %s", partitionID)
	}
*/
package Environment

import (
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"
)

<<<<<<< Updated upstream
func Mount(path string, name string) {
	fmt.Println("╔══════════════════════ INICIO MOUT  ══════════════════════════╗")
=======
var (
	mountedPartitions = make(map[string][]MountedPartition)
	mountMutex        sync.Mutex
)

func Mount(path string, name string) string {
	log.Printf("[DEBUG] Iniciando Mount con path=%s, name=%s", path, name)

	// Adquirir mutex con timeout para evitar deadlocks
	mutexAcquired := make(chan bool, 1)
	go func() {
		mountMutex.Lock()
		mutexAcquired <- true
	}()

	select {
	case <-mutexAcquired:
		defer func() {
			mountMutex.Unlock()
			log.Printf("[DEBUG] Mutex liberado")
		}()
	case <-time.After(5 * time.Second):
		log.Printf("[ERROR] Timeout al adquirir mutex")
		return "error: timeout al adquirir bloqueo del sistema"
	}

	log.Printf("[DEBUG] Abriendo archivo del disco")
>>>>>>> Stashed changes
	file, err := Utils.AbrirArchivo(path)
	if err != nil {
		msg := fmt.Sprintf("no se pudo abrir el disco en %s: %v", path, err)
		log.Printf("[ERROR] %s", msg)
		return msg
	}
	defer file.Close()

	log.Printf("[DEBUG] Leyendo MBR")
	var mbr Partitions.MBR
	if err := Utils.LeerArchivo(file, &mbr, 0); err != nil {
		msg := fmt.Sprintf("no se pudo leer el MBR: %v", err)
		log.Printf("[ERROR] %s", msg)
		return msg
	}

<<<<<<< Updated upstream
	fmt.Printf("  Buscando partición con nombre: '%s'\n", name)

	partitionFound := false
	var partition Partitions.Partition
	var partitionIndex int

	// Convertir el nombre a comparar a un arreglo de bytes de longitud fija
=======
	log.Printf("[DEBUG] Buscando partición '%s'", name)
>>>>>>> Stashed changes
	nameBytes := [16]byte{}
	copy(nameBytes[:], name)

	var foundPart *Partitions.Partition
	var partIndex int

	for i := 0; i < 4; i++ {
		part := &mbr.MbrPartitions[i]
		if bytes.Equal(part.PartName[:], nameBytes[:]) &&
			(part.PartType[0] == 'p' || part.PartType[0] == 'e') {
			foundPart = part
			partIndex = i
			break
		}
	}

	if foundPart == nil {
		msg := fmt.Sprintf("no se encontró la partición '%s' o no es primaria/extendida", name)
		log.Printf("[ERROR] %s", msg)
		return msg
	}

<<<<<<< Updated upstream
	// Verificar si la partición ya está montada
	if partition.PartStatus[0] == '1' {
		fmt.Println("Error: La partición ya está montada")
		return
	}

	fmt.Printf("  Partición encontrada: '%s' en posición %d\n", string(partition.PartName[:]), partitionIndex+1)

	// Generar el ID de la partición utilizando la función `generateDiskID`
	diskID := generateDiskID(path)
	mountedPartitionsInDisk := mountedPartitions[diskID]
	var letter byte
=======
	if foundPart.PartStatus[0] == '1' {
		msg := fmt.Sprintf("la partición '%s' ya está montada", name)
		log.Printf("[ERROR] %s", msg)
		return msg
	}

	log.Printf("[DEBUG] Generando ID para partición")
	partitionID := generatePartitionID(path, partIndex)
	log.Printf("[DEBUG] ID generado: %s", partitionID)
>>>>>>> Stashed changes

	log.Printf("[DEBUG] Actualizando estado de partición")
	foundPart.PartStatus[0] = '1'
	copy(foundPart.PartID[:], partitionID)

	log.Printf("[DEBUG] Escribiendo MBR actualizado")
	if err := Utils.EscribirArchivo(file, mbr, 0); err != nil {
		msg := fmt.Sprintf("no se pudo actualizar el MBR: %v", err)
		log.Printf("[ERROR] %s", msg)
		return msg
	}

<<<<<<< Updated upstream
	// Crear el ID de la partición utilizando el último par de dígitos de un carnet
	carnet := "202203806"
	lastTwoDigits := carnet[len(carnet)-2:]
	partitionID := fmt.Sprintf("%s%d%c", lastTwoDigits, partitionIndex+1, letter)

	partition.PartStatus[0] = '1'
	copy(partition.PartID[:], partitionID)
	TempMBR.MbrPartitions[partitionIndex] = partition

	// Almacenar la posición inicial de la partición
	mountedPartitions[diskID] = append(mountedPartitions[diskID], MountedPartition{
=======
	log.Printf("[DEBUG] Registrando partición montada")
	mountedPartitions[path] = append(mountedPartitions[path], MountedPartition{
>>>>>>> Stashed changes
		MountPath:   path,
		MountName:   name,
		MountID:     partitionID,
		MountStatus: '1',
		MountStart:  partition.PartStart,
	})

<<<<<<< Updated upstream
	// Escribir el MBR actualizado en el archivo
	if err := Utils.EscribirArchivo(file, TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo sobrescribir el MBR en el archivo")
		return
	}

	// Imprimir el mensaje confirmando que la partición ha sido montada, junto con su ID.
	fmt.Printf("  Partición montada con ID: %s\n", partitionID)
	fmt.Println("")
	fmt.Println("  MBR actualizado:")
	Partitions.ImprimirMBR(TempMBR)
	fmt.Println("")
	ImprimirMountedPartitions()

	fmt.Println("╚══════════════════════   FIN MOUNT   ═════════════════════════╝")
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
=======
	successMsg := fmt.Sprintf("partición montada con ID: %s", partitionID)
	log.Printf("[SUCCESS] %s", successMsg)
	return successMsg
}

func generatePartitionID(diskPath string, partitionIndex int) string {
	// Esta función ya no necesita el mutex porque Mount ya lo tiene
	mounts := mountedPartitions[diskPath]
	var letter byte = 'A'
	if len(mounts) > 0 {
		lastID := mounts[len(mounts)-1].MountID
		letter = lastID[len(lastID)-1] + 1
	}

	carnet := "202203806"
	lastTwoDigits := carnet[len(carnet)-2:]
	return fmt.Sprintf("%s%d%c", lastTwoDigits, partitionIndex+1, letter)
>>>>>>> Stashed changes
}
