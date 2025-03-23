package FSystem

import (
	"Proyecto/Environment"
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

// GetFileContent busca un archivo en el sistema de archivos y devuelve su contenido.
func GetFileContent(filePath string, file *os.File, superblock Partitions.Superblock) (string, error) {
	// Dividir la ruta en partes
	steps := strings.Split(filePath, "/")[1:] // Ignorar el primer elemento vacío

	// Leer el inodo raíz (inodo 0)
	var rootInode Partitions.Inode
	if err := Utils.LeerArchivo(file, &rootInode, int64(superblock.S_inode_start)); err != nil {
		return "", fmt.Errorf("error al leer el inodo raíz: %v", err)
	}

	// Buscar el inodo del archivo
	inodeIndex := SearchInodeByPath(steps, rootInode, file, superblock)
	if inodeIndex == -1 {
		return "", fmt.Errorf("archivo no encontrado: %s", filePath)
	}

	// Leer el inodo del archivo
	var fileInode Partitions.Inode
	if err := Utils.LeerArchivo(file, &fileInode, int64(superblock.S_inode_start+inodeIndex*int32(binary.Size(Partitions.Inode{})))); err != nil {
		return "", fmt.Errorf("error al leer el inodo del archivo: %v", err)
	}

	// Verificar permisos de lectura
	if !HasReadPermission(fileInode) {
		return "", fmt.Errorf("permiso denegado: no tienes permisos de lectura sobre el archivo %s", filePath)
	}

	// Leer el contenido del archivo
	content := GetInodeFileData(fileInode, file, superblock)
	return content, nil
}

// Cat procesa el comando CAT y muestra el contenido de los archivos especificados.
func Cat(files []string, id string) {
	fmt.Println("======INICIO CAT======")
	fmt.Println("Archivos:", files)
	fmt.Println("ID de partición:", id)

	// Buscar la partición montada por ID
	var mountedPartition Environment.MountedPartition
	var partitionFound bool

	for _, partitions := range Environment.GetMountedPartitions() {
		for _, partition := range partitions {
			if partition.MountID == id {
				mountedPartition = partition
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

	// Abrir el archivo del sistema de archivos
	file, err := Utils.AbrirArchivo(mountedPartition.MountPath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Leer el Superblock
	var superblock Partitions.Superblock
	if err := Utils.LeerArchivo(file, &superblock, int64(mountedPartition.MountStart)); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		return
	}

	// Leer y mostrar el contenido de cada archivo
	for _, filePath := range files {
		content, err := GetFileContent(filePath, file, superblock)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("Contenido de %s:\n%s\n", filePath, content)
	}

	fmt.Println("======FIN CAT======")
}

// HasReadPermission verifica si el usuario tiene permisos de lectura sobre el inodo.
func HasReadPermission(inode Partitions.Inode) bool {
	// Aquí debes implementar la lógica para verificar los permisos del usuario.
	// Por ahora, asumimos que el usuario tiene permisos de lectura.
	return true
}
