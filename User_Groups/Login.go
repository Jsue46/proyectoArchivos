package User_Groups

import (
	"Proyecto/Environment"
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func Login(user string, pass string, id string) {
	fmt.Println("======Start LOGIN======")
	fmt.Println("User:", user)
	fmt.Println("Pass:", pass)
	fmt.Println("Id:", id)

	mountedPartitions := Environment.GetMountedPartitions()
	var filepath string
	var partitionFound bool
	var login bool = false

	// Buscar la partición montada
	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.MountID == id {
				if partition.LoggedIn {
					fmt.Println("Error: Ya existe un usuario logueado en esta partición.")
					return
				}
				filepath = partition.MountPath
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: No se encontró la partición montada con el ID proporcionado")
		return
	}

	// Abrir el archivo del sistema de archivos binario
	file, err := Utils.AbrirArchivo(filepath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Leer el MBR
	var TempMBR Partitions.MBR
	if err := Utils.LeerArchivo(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	// Buscar la partición dentro del MBR
	var index int = -1
	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].PartSize != 0 {
			if strings.Contains(string(TempMBR.MbrPartitions[i].PartID[:]), id) {
				if TempMBR.MbrPartitions[i].PartStatus[0] == '1' {
					index = i
				} else {
					fmt.Println("Error: La partición no está montada.")
					return
				}
				break
			}
		}
	}

	if index == -1 {
		fmt.Println("Error: No se encontró la partición dentro del MBR.")
		return
	}

	// Leer el Superblock
	var tempSuperblock Partitions.Superblock
	if err := Utils.LeerArchivo(file, &tempSuperblock, int64(TempMBR.MbrPartitions[index].PartStart)); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		return
	}

	// Buscar el archivo users.txt
	indexInode := InitSearch("/users.txt", file, tempSuperblock)
	if indexInode == -1 {
		fmt.Println("Error: Archivo '/users.txt' no encontrado en la partición")
		return
	}

	// Leer el Inodo del archivo "users.txt"
	var crrInode Partitions.Inode
	if err := Utils.LeerArchivo(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Partitions.Inode{})))); err != nil {
		fmt.Println("Error: No se pudo leer el Inodo del archivo de usuarios")
		return
	}

	// Obtener el contenido del archivo users.txt
	data := GetInodeFileData(crrInode, file, tempSuperblock)
	lines := strings.Split(data, "\n")

	// Verificar credenciales
	for _, line := range lines {
		words := strings.Split(line, ",")
		if len(words) == 5 {
			if words[3] == user && words[4] == pass {
				login = true
				break
			}
		}
	}

	// Si el login es exitoso, marcar la partición como logueada
	if login {
		fmt.Println("Usuario logueado con éxito")
		Environment.MarkPartitionAsLoggedIn(id)
	} else {
		fmt.Println("Error: Usuario o contraseña incorrectos.")
	}

	fmt.Println("======End LOGIN======")
}

func InitSearch(path string, file *os.File, tempSuperblock Partitions.Superblock) int32 {
	fmt.Println("======Start BUSQUEDA INICIAL ======")
	fmt.Println("path:", path)

	// Dividir la ruta en partes usando "/" como separador
	TempStepsPath := strings.Split(path, "/")
	StepsPath := TempStepsPath[1:] // Omitir el primer elemento vacío si la ruta empieza con "/"

	fmt.Println("StepsPath:", StepsPath, "len(StepsPath):", len(StepsPath))
	for _, step := range StepsPath {
		fmt.Println("step:", step)
	}

	var Inode0 Partitions.Inode
	// Leer el inodo raíz (primer inodo del sistema de archivos)
	if err := Utils.LeerArchivo(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		return -1 // Retornar -1 si hubo un error al leer
	}

	fmt.Println("======End BUSQUEDA INICIAL======")

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
	fmt.Println("======Start BUSQUEDA INODO POR PATH======")

	index := int32(0) // Contador de bloques procesados en el inodo actual

	// Extrae el primer elemento del path y elimina espacios en blanco
	SearchedName := strings.Replace(pop(&StepsPath), " ", "", -1)

	fmt.Println("========== SearchedName:", SearchedName)

	// Iterar sobre los bloques del inodo
	for _, block := range Inode.I_block {
		if block != -1 { // Si el bloque es válido (no está vacío)
			if index < 13 { // Manejo de bloques directos (0-12)
				var crrFolderBlock Partitions.Folderblock

				// Leer el bloque de carpeta desde el archivo binario
				if err := Utils.LeerArchivo(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Partitions.Folderblock{})))); err != nil {
					return -1
				}

				// Buscar el archivo/directorio dentro del bloque de carpeta
				for _, folder := range crrFolderBlock.B_content {
					fmt.Println("Folder === Name:", string(folder.B_name[:]), "B_inodo", folder.B_inodo)

					// Si el nombre del archivo o directorio coincide
					if strings.Contains(string(folder.B_name[:]), SearchedName) {
						fmt.Println("len(StepsPath)", len(StepsPath), "StepsPath", StepsPath)

						if len(StepsPath) == 0 { // Si llegamos al final de la ruta
							fmt.Println("Folder found======")
							return folder.B_inodo // Retornar índice del inodo encontrado
						} else { // Si aún hay más niveles en la ruta, seguir buscando
							fmt.Println("NextInode======")
							var NextInode Partitions.Inode

							// Leer el siguiente inodo desde el archivo binario
							if err := Utils.LeerArchivo(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Partitions.Inode{})))); err != nil {
								return -1
							}

							// Llamada recursiva para seguir con la búsqueda
							return SarchInodeByPath(StepsPath, NextInode, file, tempSuperblock)
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
	return 0 // Si no se encontró, retornar 0
}

func GetInodeFileData(Inode Partitions.Inode, file *os.File, tempSuperblock Partitions.Superblock) string {
	fmt.Println("======Start CONTENIDO DEL BLOQUE======")
	index := int32(0)
	// define content as a string
	var content string

	// Iterate over i_blocks from Inode
	for _, block := range Inode.I_block {
		if block != -1 {
			//Dentro de los directos
			if index < 13 {
				var crrFileBlock Partitions.Fileblock
				// Read object from bin file
				if err := Utils.LeerArchivo(file, &crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Partitions.Fileblock{})))); err != nil {
					return ""
				}

				content += string(crrFileBlock.B_content[:])

			} else {
				fmt.Print("indirectos")
			}
		}
		index++
	}

	fmt.Println("======End CONTENIDO DEL BLOQUE======")
	return content
}

// MKUSER
func AppendToFileBlock(inode *Partitions.Inode, newData string, file *os.File, superblock Partitions.Superblock) error {
	// Leer el contenido existente del archivo utilizando la función GetInodeFileData
	existingData := GetInodeFileData(*inode, file, superblock)

	// Concatenar el nuevo contenido
	fullData := existingData + newData

	// Asegurarse de que el contenido no exceda el tamaño del bloque
	if len(fullData) > len(inode.I_block)*binary.Size(Partitions.Fileblock{}) {
		// Si el contenido excede, necesitas manejar bloques adicionales
		return fmt.Errorf("el tamaño del archivo excede la capacidad del bloque actual y no se ha implementado la creación de bloques adicionales")
	}

	// Escribir el contenido actualizado en el bloque existente
	var updatedFileBlock Partitions.Fileblock
	copy(updatedFileBlock.B_content[:], fullData)
	if err := Utils.EscribirArchivo(file, updatedFileBlock, int64(superblock.S_block_start+inode.I_block[0]*int32(binary.Size(Partitions.Fileblock{})))); err != nil {
		return fmt.Errorf("error al escribir el bloque actualizado: %v", err)
	}

	// Actualizar el tamaño del inodo
	inode.I_size = int32(len(fullData))
	if err := Utils.EscribirArchivo(file, *inode, int64(superblock.S_inode_start+inode.I_block[0]*int32(binary.Size(Partitions.Inode{})))); err != nil {
		return fmt.Errorf("error al actualizar el inodo: %v", err)
	}

	return nil
}
