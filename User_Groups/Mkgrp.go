package User_Groups

import (
	"Proyecto/Environment"
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func MKGRP(name string) {
	fmt.Println("╔══════════════════════ INICIO MKGRP  ═════════════════════════╗")

	fmt.Println("  Nombre del grupo:", name)

	// Verificar si el usuario actual es root
	if !IsRootLoggedIn() {
		fmt.Println("Error: Solo el usuario root puede crear grupos.")
		return
	}

	// Obtener la partición montada actualmente
	mountedPartition := Environment.GetCurrentMountedPartition()
	if mountedPartition == nil {
		fmt.Println("Error: No hay ninguna partición montada.")
		return
	}

	// Abrir el archivo del sistema de archivos
	file, err := Utils.AbrirArchivo(mountedPartition.MountPath)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Leer el Superblock de la partición
	var tempSuperblock Partitions.Superblock
	if err := Utils.LeerArchivo(file, &tempSuperblock, int64(mountedPartition.MountStart)); err != nil {
		fmt.Println("Error al leer el Superblock:", err)
		return
	}

	// Buscar el archivo "users.txt" en el sistema de archivos
	indexInode := InitSearch("/users.txt", file, tempSuperblock)
	if indexInode == -1 {
		fmt.Println("Error: No se encontró el archivo users.txt.")
		return
	}

	// Leer el inodo del archivo "users.txt"
	var crrInode Partitions.Inode
	if err := Utils.LeerArchivo(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Partitions.Inode{})))); err != nil {
		fmt.Println("Error al leer el inodo de users.txt:", err)
		return
	}

	// Obtener el contenido actual del archivo "users.txt"
	data := GetInodeFileData(crrInode, file, tempSuperblock)

	// Verificar si el grupo ya existe
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if strings.Contains(line, ",G,") { // Buscar ",G," en lugar de "G,"
			parts := strings.Split(line, ",")
			if len(parts) >= 3 && strings.TrimSpace(parts[2]) == name {
				fmt.Println("Error: El grupo ya existe.")
				return
			}
		}
	}

	// Crear el nuevo grupo
	newGroupID := getNextGroupID(lines)
	newGroupLine := fmt.Sprintf("%d,G,%s\n", newGroupID, name)

	// Asignar un nuevo bloque para el grupo
	newBlockIndex := findFreeBlock(file, tempSuperblock)
	if newBlockIndex == -1 {
		fmt.Println("❌ Error: No hay bloques libres disponibles.")
		return
	}

	// Escribir el nuevo grupo en el bloque
	var fileBlock Partitions.Fileblock
	copy(fileBlock.B_content[:], newGroupLine)
	blockOffset := int64(tempSuperblock.S_block_start + newBlockIndex*int32(binary.Size(Partitions.Fileblock{})))
	if err := Utils.EscribirArchivo(file, fileBlock, blockOffset); err != nil {
		fmt.Println("❌ Error al escribir el bloque de archivo:", err)
		return
	}

	// Asignar el nuevo bloque al inodo
	for i := 0; i < len(crrInode.I_block); i++ {
		if crrInode.I_block[i] == -1 {
			crrInode.I_block[i] = int32(newBlockIndex)
			break
		}
	}

	// Actualizar el tamaño del inodo
	crrInode.I_size += int32(len(newGroupLine))

	// Escribir el inodo actualizado en el archivo
	inodeOffset := int64(tempSuperblock.S_inode_start + indexInode*int32(binary.Size(Partitions.Inode{})))
	if err := Utils.EscribirArchivo(file, crrInode, inodeOffset); err != nil {
		fmt.Println("❌ Error al actualizar el inodo:", err)
		return
	}

	fmt.Println("✅ Grupo creado con éxito.")
	fmt.Println("╚══════════════════════   FIN MKGRP   ═════════════════════════╝")
}

// Función auxiliar para obtener el siguiente ID de grupo
func getNextGroupID(lines []string) int {
	maxID := 0
	for _, line := range lines {
		if strings.Contains(line, ",G,") { // Buscar ",G," en lugar de "G,"
			parts := strings.Split(line, ",")
			if len(parts) >= 1 {
				id, err := strconv.Atoi(parts[0])
				if err == nil && id > maxID {
					maxID = id
				}
			}
		}
	}
	return maxID + 1
}
