package FSystem

import (
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"encoding/binary"
	"fmt"
	"os"
)

func create_ext2(n int32, partition Partitions.Partition, newSuperblock Partitions.Superblock, date string, file *os.File) {
	fmt.Println("======Start CREATE EXT2======")
	fmt.Println("INODOS:", n)

	// Escribir los bitmaps de inodos y bloques
	writeBitmaps(n, newSuperblock, file)

	// Inicializar los inodos y bloques con valores predeterminados
	if err := initInodesAndBlocks(n, newSuperblock, file); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Crear la carpeta raíz y el archivo "users.txt"
	if err := createRootAndUsersFile(newSuperblock, date, file); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Escribir el superbloque actualizado en el archivo
	if err := Utils.EscribirArchivo(file, newSuperblock, int64(partition.PartStart)); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Marcar los primeros inodos y bloques como usados
	if err := markUsedInodesAndBlocks(newSuperblock, file); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("======End CREATE EXT2======")
}

func writeBitmaps(n int32, newSuperblock Partitions.Superblock, file *os.File) {
	for i := int32(0); i < n; i++ {
		if err := Utils.EscribirArchivo(file, byte(0), int64(newSuperblock.S_bm_inode_start+i)); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}

	for i := int32(0); i < 3*n; i++ {
		if err := Utils.EscribirArchivo(file, byte(0), int64(newSuperblock.S_bm_block_start+i)); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
}

// Función auxiliar para crear la carpeta raíz y el archivo users.txt
func createRootAndUsersFile(newSuperblock Partitions.Superblock, date string, file *os.File) error {
	var Inode0, Inode1 Partitions.Inode
	initInode(&Inode0, date)
	initInode(&Inode1, date)

	Inode0.I_block[0] = 0
	Inode1.I_block[0] = 1

	// Asignar el tamaño real del contenido
	data := "1,G,root\n1,U,root,root,123\n"
	actualSize := int32(len(data))
	Inode1.I_size = actualSize // Esto ahora refleja el tamaño real del contenido

	var Fileblock1 Partitions.Fileblock
	copy(Fileblock1.B_content[:], data) // Copia segura de datos a Fileblock

	var Folderblock0 Partitions.Folderblock
	Folderblock0.B_content[0].B_inodo = 0
	copy(Folderblock0.B_content[0].B_name[:], ".")
	Folderblock0.B_content[1].B_inodo = 0
	copy(Folderblock0.B_content[1].B_name[:], "..")
	Folderblock0.B_content[2].B_inodo = 1
	copy(Folderblock0.B_content[2].B_name[:], "users.txt")

	// Escribir los inodos y bloques en las posiciones correctas
	if err := Utils.EscribirArchivo(file, Inode0, int64(newSuperblock.S_inode_start)); err != nil {
		return err
	}
	if err := Utils.EscribirArchivo(file, Inode1, int64(newSuperblock.S_inode_start+int32(binary.Size(Partitions.Inode{})))); err != nil {
		return err
	}
	if err := Utils.EscribirArchivo(file, Folderblock0, int64(newSuperblock.S_block_start)); err != nil {
		return err
	}
	if err := Utils.EscribirArchivo(file, Fileblock1, int64(newSuperblock.S_block_start+int32(binary.Size(Partitions.Folderblock{})))); err != nil {
		return err
	}

	return nil
}

// Función auxiliar para inicializar inodos y bloques
func initInodesAndBlocks(n int32, newSuperblock Partitions.Superblock, file *os.File) error {
	var newInode Partitions.Inode
	for i := int32(0); i < 15; i++ {
		newInode.I_block[i] = -1
	}

	for i := int32(0); i < n; i++ {
		if err := Utils.EscribirArchivo(file, newInode, int64(newSuperblock.S_inode_start+i*int32(binary.Size(Partitions.Inode{})))); err != nil {
			return err
		}
	}

	var newFileblock Partitions.Fileblock
	for i := int32(0); i < 3*n; i++ {
		if err := Utils.EscribirArchivo(file, newFileblock, int64(newSuperblock.S_block_start+i*int32(binary.Size(Partitions.Fileblock{})))); err != nil {
			return err
		}
	}

	return nil
}
