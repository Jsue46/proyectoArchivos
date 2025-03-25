package User_Groups

import (
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func InitSearch(path string, file *os.File, tempSuperblock Partitions.Superblock) int32 {
	fmt.Println(" ══════════════════  BUSQUEDA INICIAL  ════════════════════════ ")
	fmt.Println("  path:", path)

	// Dividir la ruta en partes usando "/" como separador
	TempStepsPath := strings.Split(path, "/")
	StepsPath := TempStepsPath[1:]

	fmt.Println("  StepsPath:", StepsPath, "len(StepsPath):", len(StepsPath))
	for _, step := range StepsPath {
		fmt.Println("  step:", step)
	}

	var Inode0 Partitions.Inode
	// Leer el inodo raíz
	if err := Utils.LeerArchivo(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		return -1
	}

	fmt.Println(" ════════════════════  FIN BUSQUEDA   ═════════════════════════ ")

	// Llamar a la función que busca el inodo del archivo según la ruta
	return SarchInodeByPath(StepsPath, Inode0, file, tempSuperblock)
}

// stack
func pop(s *[]string) string {
	lastIndex := len(*s) - 1
	last := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return last
}

func SarchInodeByPath(StepsPath []string, Inode Partitions.Inode, file *os.File, tempSuperblock Partitions.Superblock) int32 {
	fmt.Println(" ═══════════════  BUSQUEDA INODO POR PATH   ═══════════════════ ")

	index := int32(0)

	// Extrae el primer elemento del path y elimina espacios en blanco
	SearchedName := strings.Replace(pop(&StepsPath), " ", "", -1)

	fmt.Println(" ═══════════════ SearchedName:", SearchedName)

	// Iterar sobre los bloques del inodo
	for _, block := range Inode.I_block {
		if block != -1 {
			if index < 13 {
				var crrFolderBlock Partitions.Folderblock

				// Leer el bloque de carpeta desde el archivo binario
				if err := Utils.LeerArchivo(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Partitions.Folderblock{})))); err != nil {
					return -1
				}

				// Buscar el archivo/directorio dentro del bloque de carpeta
				for _, folder := range crrFolderBlock.B_content {
					fmt.Println(" ═════ Folder Name:", string(folder.B_name[:]), " B_inodo", folder.B_inodo)

					// Si el nombre del archivo o directorio coincide
					if strings.Contains(string(folder.B_name[:]), SearchedName) {
						fmt.Println("\tlen(StepsPath)", len(StepsPath), "StepsPath", StepsPath)

						if len(StepsPath) == 0 {
							fmt.Println(" ═════ Folder found ═════ ")
							return folder.B_inodo
						} else {
							fmt.Println(" ═════ NextInode ═════ ")
							var NextInode Partitions.Inode

							if err := Utils.LeerArchivo(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Partitions.Inode{})))); err != nil {
								return -1
							}

							// Llamada recursiva para seguir con la búsqueda
							return SarchInodeByPath(StepsPath, NextInode, file, tempSuperblock)
						}
					}
				}
			} else {
				fmt.Println("Manejo de bloques indirectos no implementado")
			}
		}
		index++
	}

	fmt.Println(" ═════════════  FIN BUSQUEDA INODO POR PATH   ═════════════════ ")
	return 0
}

func GetInodeFileData(Inode Partitions.Inode, file *os.File, tempSuperblock Partitions.Superblock) string {
	fmt.Println(" ═══════════════  CONTENIDO DEL BLOQUE   ══════════════════════ ")
	index := int32(0)
	var content string

	// Iterate over i_blocks from Inode
	for _, block := range Inode.I_block {
		if block != -1 {
			//Dentro de los directos
			if index < 13 {
				var crrFileBlock Partitions.Fileblock
				if err := Utils.LeerArchivo(file, &crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Partitions.Fileblock{})))); err != nil {
					fmt.Println("Error al leer el bloque de archivo:", err)
					return ""
				}

				// Mostrar el contenido del bloque
				fmt.Println(" ", string(crrFileBlock.B_content[:]))

				// Agregar el contenido del bloque al resultado final
				content += string(crrFileBlock.B_content[:])

			} else {
				fmt.Println("Manejo de bloques indirectos no implementado")
			}
		}
		index++
	}

	fmt.Println(" ═══════════════  FIN CONTENIDO DEL BLOQUE   ══════════════════ ")
	return content
}

func AppendToFileBlock(inode *Partitions.Inode, newData string, file *os.File, superblock Partitions.Superblock) error {
	fmt.Println(" ═══════════════   AGREAGAR AL BLOQUE    ══════════════════════ ")
	// Obtener el contenido actual del archivo
	existingData := GetInodeFileData(*inode, file, superblock)
	fmt.Println("🔹 Contenido actual de users.txt:")
	fmt.Println(existingData)

	// Unir los datos en un solo string
	fullData := existingData + newData
	fmt.Println("🔹 Nuevo contenido tras agregar:", newData)

	// Tamaño de un bloque
	blockSize := binary.Size(Partitions.Fileblock{})

	// Obtener el índice del último bloque usado
	lastBlockIndex := -1
	for i := 0; i < len(inode.I_block); i++ {
		if inode.I_block[i] != -1 {
			lastBlockIndex = i
		} else {
			break
		}
	}
	fmt.Printf("🔹 Último bloque usado en inode: %d\n", lastBlockIndex)

	// Si no hay bloques, asignamos el primero
	if lastBlockIndex == -1 {
		newBlockIndex := findFreeBlock(file, superblock)
		if newBlockIndex == -1 {
			fmt.Println("❌ Error: No hay bloques libres disponibles")
			return fmt.Errorf("error: no hay bloques libres disponibles")
		}
		inode.I_block[0] = int32(newBlockIndex)
		lastBlockIndex = 0
	}

	// Obtener el bloque actual donde se escribe
	blockOffset := int64(superblock.S_block_start + inode.I_block[lastBlockIndex]*int32(blockSize))

	var fileBlock Partitions.Fileblock

	// Leer el bloque actual
	if err := Utils.LeerArchivo(file, &fileBlock, blockOffset); err != nil {
		fmt.Println("❌ Error al leer el bloque de archivo:", err)
		return err
	}

	// Verificar cuánto espacio libre queda en el bloque actual
	existingContent := strings.TrimRight(string(fileBlock.B_content[:]), "\x00")
	remainingSpace := blockSize - len(existingContent)

	// Si hay espacio, escribir en el mismo bloque
	if len(newData) <= remainingSpace {
		copy(fileBlock.B_content[len(existingContent):], newData)
	} else {
		copy(fileBlock.B_content[:], newData)
	}

	// Escribir el bloque actualizado en el archivo
	if err := Utils.EscribirArchivo(file, fileBlock, blockOffset); err != nil {
		fmt.Println("❌ Error al escribir el bloque de archivo:", err)
		return fmt.Errorf("error al escribir el bloque de archivo: %v", err)
	}

	// Actualizar el tamaño del inodo
	inode.I_size = int32(len(fullData))
	inodeOffset := int64(superblock.S_inode_start + inode.I_block[0]*int32(binary.Size(Partitions.Inode{})))

	if err := Utils.EscribirArchivo(file, *inode, inodeOffset); err != nil {
		fmt.Println("❌ Error al actualizar el inodo:", err)
		return fmt.Errorf("error al actualizar el inodo: %v", err)
	}
	fmt.Println(" ═══════════════   FIN AGREAGAR AL BLOQUE   ═══════════════════ ")
	return nil
}

func findFreeBlock(file *os.File, superblock Partitions.Superblock) int32 {
	var blockBitmap []byte = make([]byte, superblock.S_blocks_count)
	// Leer el bitmap de bloques
	if err := Utils.LeerArchivo(file, &blockBitmap, int64(superblock.S_bm_block_start)); err != nil {
		fmt.Println("Error al leer el bitmap de bloques:", err)
		return -1
	}
	// Buscar el primer bloque libre
	for i, b := range blockBitmap {
		if b == 0 {
			// Marcar el bloque como usado
			blockBitmap[i] = 1
			if err := Utils.EscribirArchivo(file, blockBitmap, int64(superblock.S_bm_block_start)); err != nil {
				fmt.Println("Error al actualizar el bitmap de bloques:", err)
				return -1
			}
			return int32(i)
		}
	}

	return -1
}
