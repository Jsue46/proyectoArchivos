package FSystem

import (
	"Proyecto/Environment"
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"encoding/binary"
	"fmt"
	"strings"
)

func Mkfs(id string, type_ string, fs_ string) {
	fmt.Println("======INICIO MKFS======")
	fmt.Println("Id:", id)
	fmt.Println("Type:", type_)
	fmt.Println("Fs:", fs_)

	// Buscar la partición montada por ID
	mountedPartition, partitionFound := findMountedPartition(id)
	if !partitionFound {
		fmt.Println("Particion no encontrada")
		return
	}

	// Verificar si la partición está montada (estado '1')
	if mountedPartition.MountStatus != '1' {
		fmt.Println("La particion aun no esta montada")
		return
	}

	// Abrir el archivo binario de la partición
	file, err := Utils.AbrirArchivo(mountedPartition.MountPath)
	if err != nil {
		return
	}
	defer file.Close()

	// Leer el MBR (Master Boot Record) del archivo binario
	var TempMBR Partitions.MBR
	if err := Utils.LeerArchivo(file, &TempMBR, 0); err != nil {
		return
	}

	// Buscar la partición dentro del MBR que coincida con el ID proporcionado
	index := findPartitionIndex(TempMBR, id)
	if index == -1 {
		fmt.Println("Particion no encontrada (2)")
		return
	}

	// Calcular el número de inodos basado en el tamaño de la partición
	n := calculateInodes(TempMBR.MbrPartitions[index], fs_)

	// Crear el Superblock con los campos calculados
	newSuperblock := createSuperblock(TempMBR.MbrPartitions[index], n)

	// Si el sistema de archivos es 2FS, se crea un sistema de archivos EXT2
	if fs_ == "2fs" {
		create_ext2(n, TempMBR.MbrPartitions[index], newSuperblock, "25/02/2025", file)
	} else {
		fmt.Println("EXT3 no está soportado.")
	}

	fmt.Println("======FIN MKFS======")
}

// Funciones auxiliares
func findMountedPartition(id string) (Environment.MountedPartition, bool) {
	for _, partitions := range Environment.GetMountedPartitions() {
		for _, partition := range partitions {
			if partition.MountID == id {
				return partition, true
			}
		}
	}
	return Environment.MountedPartition{}, false
}

func findPartitionIndex(mbr Partitions.MBR, id string) int {
	for i := 0; i < 4; i++ {
		if mbr.MbrPartitions[i].PartSize != 0 && strings.Contains(string(mbr.MbrPartitions[i].PartID[:]), id) {
			return i
		}
	}
	return -1
}

func calculateInodes(partition Partitions.Partition, fs string) int32 {
	numerador := int32(partition.PartSize - int32(binary.Size(Partitions.Superblock{})))
	denominador_base := int32(4 + int32(binary.Size(Partitions.Inode{})) + 3*int32(binary.Size(Partitions.Fileblock{})))
	var temp int32 = 0
	if fs == "2fs" {
		temp = 0
	} else {
		fmt.Print("Error por el momento solo está disponible 2FS.")
	}
	return numerador / (denominador_base + temp)
}

func createSuperblock(partition Partitions.Partition, n int32) Partitions.Superblock {
	var newSuperblock Partitions.Superblock
	newSuperblock.S_filesystem_type = 2 // EXT2
	newSuperblock.S_inodes_count = n
	newSuperblock.S_blocks_count = 3 * n
	newSuperblock.S_free_blocks_count = 3*n - 2
	newSuperblock.S_free_inodes_count = n - 2
	copy(newSuperblock.S_mtime[:], "25/02/2025")
	copy(newSuperblock.S_umtime[:], "25/02/2025")
	newSuperblock.S_mnt_count = 1
	newSuperblock.S_magic = 0xEF53
	newSuperblock.S_inode_size = int32(binary.Size(Partitions.Inode{}))
	newSuperblock.S_block_size = int32(binary.Size(Partitions.Fileblock{}))

	newSuperblock.S_bm_inode_start = partition.PartStart + int32(binary.Size(Partitions.Superblock{}))
	newSuperblock.S_bm_block_start = newSuperblock.S_bm_inode_start + n
	newSuperblock.S_inode_start = newSuperblock.S_bm_block_start + 3*n
	newSuperblock.S_block_start = newSuperblock.S_inode_start + n*newSuperblock.S_inode_size

	return newSuperblock
}
