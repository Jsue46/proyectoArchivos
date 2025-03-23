package FSystem

import (
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

// Función auxiliar para crear el sistema de archivos EXT2
func create_ext2(n int32, partition Partitions.Partition, newSuperblock Partitions.Superblock, date string, file *os.File) {
	fmt.Println("======Start CREATE EXT2======")
	fmt.Println("INODOS:", n)

	// Imprimir el Superblock calculado
	Partitions.PrintSuperblock(newSuperblock)
	fmt.Println("Date:", date)

	// Escribir los bitmaps de inodos y bloques
	for i := int32(0); i < n; i++ {
		if err := Utils.EscribirArchivo(file, byte(0), int64(newSuperblock.S_bm_inode_start+i)); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}

	// Escribir los bitmaps de bloques
	for i := int32(0); i < 3*n; i++ {
		if err := Utils.EscribirArchivo(file, byte(0), int64(newSuperblock.S_bm_block_start+i)); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}

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

	// Leer e imprimir los inodos después de formatear
	fmt.Println("====== Imprimiendo Inodos ======")
	for i := int32(0); i < n; i++ {
		var inode Partitions.Inode
		offset := int64(newSuperblock.S_inode_start + i*int32(binary.Size(Partitions.Inode{})))
		if err := Utils.LeerArchivo(file, &inode, offset); err != nil {
			fmt.Println("Error al leer inodo: ", err)
			return
		}
		Partitions.PrintInode(inode)
	}

	// Leer e imprimir los Folderblocks y Fileblocks
	fmt.Println("====== Imprimiendo Folderblocks y Fileblocks ======")

	// Imprimir Folderblocks
	for i := int32(0); i < 1; i++ {
		var folderblock Partitions.Folderblock
		offset := int64(newSuperblock.S_block_start + i*int32(binary.Size(Partitions.Folderblock{})))
		if err := Utils.LeerArchivo(file, &folderblock, offset); err != nil {
			fmt.Println("Error al leer Folderblock: ", err)
			return
		}
		Partitions.PrintFolderblock(folderblock)
	}

	// Imprimir Fileblocks
	for i := int32(0); i < 1; i++ {
		var fileblock Partitions.Fileblock
		offset := int64(newSuperblock.S_block_start + int32(binary.Size(Partitions.Folderblock{})) + i*int32(binary.Size(Partitions.Fileblock{})))
		if err := Utils.LeerArchivo(file, &fileblock, offset); err != nil {
			fmt.Println("Error al leer Fileblock: ", err)
			return
		}
		Partitions.PrintFileblock(fileblock)
	}

	// Imprimir el Superblock final
	Partitions.PrintSuperblock(newSuperblock)

	fmt.Println("======End CREATE EXT2======")
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

// ----------------------------------------------------------------------------------------------- +Agregado
// Función auxiliar para crear la carpeta raíz y el archivo users.txt
func createRootAndUsersFile(newSuperblock Partitions.Superblock, date string, file *os.File) error {
	var Inode0, Inode1 Partitions.Inode

	// Inicializa los inodos con la fecha proporcionada
	initInode(&Inode0, date)
	initInode(&Inode1, date)

	// Asigna los bloques correspondientes a los inodos
	Inode0.I_block[0] = 0 // El primer inodo representa la carpeta raíz
	Inode1.I_block[0] = 1 // El segundo inodo representa el archivo users.txt

	// Contenido del archivo users.txt (usuarios iniciales)
	data := "1,G,root\n1,U,root,root,123\n" // Define un grupo y un usuario root
	actualSize := int32(len(data))          // Obtiene el tamaño real del contenido
	Inode1.I_size = actualSize              // Asigna el tamaño real al inodo del archivo

	// Crea un bloque de archivo y copia el contenido en él
	var Fileblock1 Partitions.Fileblock
	copy(Fileblock1.B_content[:], data) // Copia segura de datos a Fileblock

	// Crea un bloque de carpeta (raíz) y asigna las entradas iniciales
	var Folderblock0 Partitions.Folderblock
	Folderblock0.B_content[0].B_inodo = 0 // Referencia al propio directorio
	copy(Folderblock0.B_content[0].B_name[:], ".")
	Folderblock0.B_content[1].B_inodo = 0 // Referencia al directorio padre (en este caso, él mismo)
	copy(Folderblock0.B_content[1].B_name[:], "..")
	Folderblock0.B_content[2].B_inodo = 1 // Apunta al archivo users.txt
	copy(Folderblock0.B_content[2].B_name[:], "users.txt")

	// Escribe los inodos y bloques en las posiciones correctas en el archivo del sistema
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

// SearchInodeByPath busca un archivo o directorio en el sistema de archivos EXT2.
func SearchInodeByPath(steps []string, inode Partitions.Inode, file *os.File, superblock Partitions.Superblock) int32 {
	fmt.Println("======Start BUSQUEDA INODO POR PATH======")

	index := int32(0) // Contador de bloques procesados en el inodo actual

	// Extrae el primer elemento del path y elimina espacios en blanco
	SearchedName := strings.Replace(steps[0], " ", "", -1)

	fmt.Println("========== SearchedName:", SearchedName)

	// Iterar sobre los bloques del inodo
	for _, block := range inode.I_block {
		if block != -1 { // Si el bloque es válido (no está vacío)
			if index < 13 { // Manejo de bloques directos (0-12)
				var crrFolderblock Partitions.Folderblock

				// Leer el bloque de carpeta desde el archivo binario
				if err := Utils.LeerArchivo(file, &crrFolderblock, int64(superblock.S_block_start+block*int32(binary.Size(Partitions.Folderblock{})))); err != nil {
					return -1
				}

				// Buscar el archivo/directorio dentro del bloque de carpeta
				for _, folder := range crrFolderblock.B_content {
					fmt.Println("Folder === Name:", string(folder.B_name[:]), "B_inodo", folder.B_inodo)

					// Si el nombre del archivo o directorio coincide
					if strings.Contains(string(folder.B_name[:]), SearchedName) {
						fmt.Println("len(steps)", len(steps), "steps", steps)

						if len(steps) == 1 { // Si llegamos al final de la ruta
							fmt.Println("Folder found======")
							return folder.B_inodo // Retornar índice del inodo encontrado
						} else { // Si aún hay más niveles en la ruta, seguir buscando
							fmt.Println("NextInode======")
							var NextInode Partitions.Inode

							// Leer el siguiente inodo desde el archivo binario
							if err := Utils.LeerArchivo(file, &NextInode, int64(superblock.S_inode_start+folder.B_inodo*int32(binary.Size(Partitions.Inode{})))); err != nil {
								return -1
							}

							// Llamada recursiva para seguir con la búsqueda
							return SearchInodeByPath(steps[1:], NextInode, file, superblock)
						}
					}
				}
			} else {
				fmt.Println("Manejo de bloques indirectos no implementado") // Falta implementar acceso a bloques indirectos
			}
		}
		index++ // Incrementar índice para saber si son bloques directos o indirectos
	}

	fmt.Println("======End BUSQUEDA INODO POR PATH======")
	return -1 // Si no se encontró, retornar -1
}
