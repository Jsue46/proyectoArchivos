package User_Groups

import (
	"Proyecto/Environment"
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CHGRP(user string, grp string) {
	fmt.Println("╔══════════════════════ INICIO CHGRP  ═════════════════════════╗")
	fmt.Println("  Usuario:", user)
	fmt.Println("  Nuevo grupo:", grp)

	if !IsRootLoggedIn() {
		fmt.Println("❌ Error: Solo el usuario root puede cambiar grupos de usuarios.")
		return
	}

	if len(user) > 10 || len(grp) > 10 {
		fmt.Println("❌ Error: user y grp deben tener máximo 10 caracteres.")
		return
	}

	mountedPartition := Environment.GetCurrentMountedPartition()
	if mountedPartition == nil {
		fmt.Println("❌ Error: No hay ninguna partición montada.")
		return
	}

	file, err := Utils.AbrirArchivo(mountedPartition.MountPath)
	if err != nil {
		fmt.Println("❌ Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	var tempSuperblock Partitions.Superblock
	if err := Utils.LeerArchivo(file, &tempSuperblock, int64(mountedPartition.MountStart)); err != nil {
		fmt.Println("❌ Error al leer el Superblock:", err)
		return
	}

	indexInode := InitSearch("/users.txt", file, tempSuperblock)
	if indexInode == -1 {
		fmt.Println("❌ Error: No se encontró el archivo users.txt.")
		return
	}

	var crrInode Partitions.Inode
	if err := Utils.LeerArchivo(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Partitions.Inode{})))); err != nil {
		fmt.Println("❌ Error al leer el inodo de users.txt:", err)
		return
	}

	blockSize := binary.Size(Partitions.Fileblock{})
	userFound := false
	var userLine string
	var userBlockIndex int32 = -1
	var targetGroupBlockIndex int32 = -1
	newGroupID := -1
	var userWasOnlyContent bool = false

	// Primera pasada: Buscar usuario y validar nuevo grupo
	for i, block := range crrInode.I_block {
		if block == -1 {
			continue
		}

		var fileBlock Partitions.Fileblock
		blockOffset := int64(tempSuperblock.S_block_start + block*int32(blockSize))

		if err := Utils.LeerArchivo(file, &fileBlock, blockOffset); err != nil {
			fmt.Println("❌ Error al leer bloque de archivo:", err)
			return
		}

		content := string(fileBlock.B_content[:])
		lines := strings.Split(content, "\n")
		nonEmptyLines := 0

		// Buscar el usuario y contar líneas no vacías
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			nonEmptyLines++

			parts := strings.Split(line, ",")
			if len(parts) >= 5 && strings.TrimSpace(parts[1]) == "U" && strings.TrimSpace(parts[3]) == user {
				if id, err := strconv.Atoi(strings.TrimSpace(parts[0])); err != nil || id <= 0 {
					fmt.Println("❌ Error: El usuario está eliminado.")
					return
				}
				userLine = line
				userBlockIndex = int32(i)
				userFound = true

				// Verificar si este bloque solo contiene al usuario
				if nonEmptyLines == 1 {
					userWasOnlyContent = true
				}
			}

			// Buscar el nuevo grupo
			if len(parts) >= 3 && strings.TrimSpace(parts[1]) == "G" && strings.TrimSpace(parts[2]) == grp {
				if id, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil && id > 0 {
					newGroupID = id
					targetGroupBlockIndex = int32(i)
				}
			}
		}
	}

	if !userFound {
		fmt.Println("❌ Error: El usuario no existe en el sistema.")
		return
	}

	if newGroupID == -1 {
		fmt.Println("❌ Error: El grupo no existe o está eliminado.")
		return
	}

	// Eliminar usuario del bloque actual
	var userBlock Partitions.Fileblock
	userBlockOffset := int64(tempSuperblock.S_block_start + crrInode.I_block[userBlockIndex]*int32(blockSize))
	if err := Utils.LeerArchivo(file, &userBlock, userBlockOffset); err != nil {
		fmt.Println("❌ Error al leer bloque de usuario:", err)
		return
	}

	userContent := string(userBlock.B_content[:])
	userLines := strings.Split(userContent, "\n")
	var newUserLines []string
	for _, line := range userLines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) >= 5 && strings.TrimSpace(parts[1]) == "U" && strings.TrimSpace(parts[3]) == user {
			continue // Eliminar esta línea
		}
		newUserLines = append(newUserLines, line)
	}

	// Si el bloque solo contenía al usuario, eliminarlo completamente
	if userWasOnlyContent {
		fmt.Println("🗑️ Eliminando bloque que solo contenía al usuario...")

		// Liberar el bloque
		if err := freeBlock(file, tempSuperblock, crrInode.I_block[userBlockIndex]); err != nil {
			fmt.Println("❌ Error al liberar bloque:", err)
			return
		}

		// Desplazar los bloques posteriores
		for i := int(userBlockIndex); i < len(crrInode.I_block)-1; i++ {
			crrInode.I_block[i] = crrInode.I_block[i+1]
		}
		crrInode.I_block[len(crrInode.I_block)-1] = -1

		// Actualizar el inodo
		inodeOffset := int64(tempSuperblock.S_inode_start + indexInode*int32(binary.Size(Partitions.Inode{})))
		if err := Utils.EscribirArchivo(file, crrInode, inodeOffset); err != nil {
			fmt.Println("❌ Error al actualizar inodo:", err)
			return
		}
	} else {
		// Actualizar bloque original sin el usuario
		userBlock.B_content = [64]byte{}
		copy(userBlock.B_content[:], strings.Join(newUserLines, "\n"))
		if err := Utils.EscribirArchivo(file, userBlock, userBlockOffset); err != nil {
			fmt.Println("❌ Error al actualizar bloque original:", err)
			return
		}
	}

	// Actualizar el ID de grupo en la línea del usuario
	parts := strings.Split(userLine, ",")
	parts[0] = strconv.Itoa(newGroupID)
	parts[2] = grp
	updatedUserLine := strings.Join(parts, ",")

	// Estrategia de colocación inteligente
	var targetBlock Partitions.Fileblock
	targetBlockOffset := int64(tempSuperblock.S_block_start + crrInode.I_block[targetGroupBlockIndex]*int32(blockSize))
	if err := Utils.LeerArchivo(file, &targetBlock, targetBlockOffset); err != nil {
		fmt.Println("❌ Error al leer bloque destino:", err)
		return
	}

	targetContent := string(targetBlock.B_content[:])
	targetLines := strings.Split(targetContent, "\n")

	if len(targetContent)+len(updatedUserLine)+1 <= 64 {
		// Agregar al bloque del grupo
		targetLines = append(targetLines, updatedUserLine)
		newTargetContent := strings.Join(targetLines, "\n")

		targetBlock.B_content = [64]byte{}
		copy(targetBlock.B_content[:], newTargetContent)
		if err := Utils.EscribirArchivo(file, targetBlock, targetBlockOffset); err != nil {
			fmt.Println("❌ Error al actualizar bloque destino:", err)
			return
		}
		fmt.Println("✅ Usuario agregado al bloque del grupo destino.")
	} else {
		// Crear nuevo bloque contiguo al del grupo
		fmt.Println("⚠️ No hay espacio en bloque actual, creando nuevo bloque...")

		// Buscar un bloque libre
		newBlockIndex := findFreeBlock(file, tempSuperblock)
		if newBlockIndex == -1 {
			fmt.Println("❌ Error: No hay bloques libres disponibles.")
			return
		}

		// Buscar la posición para insertar el nuevo bloque (después del bloque del grupo)
		insertPosition := -1
		for i := 0; i < len(crrInode.I_block); i++ {
			if crrInode.I_block[i] == crrInode.I_block[targetGroupBlockIndex] {
				insertPosition = i
				break
			}
		}

		if insertPosition == -1 {
			fmt.Println("❌ Error: No se pudo determinar posición para nuevo bloque.")
			return
		}

		// Desplazar bloques posteriores para hacer espacio
		for i := len(crrInode.I_block) - 1; i > insertPosition+1; i-- {
			crrInode.I_block[i] = crrInode.I_block[i-1]
		}

		// Asignar nuevo bloque
		crrInode.I_block[insertPosition+1] = newBlockIndex

		// Actualizar el inodo
		inodeOffset := int64(tempSuperblock.S_inode_start + indexInode*int32(binary.Size(Partitions.Inode{})))
		if err := Utils.EscribirArchivo(file, crrInode, inodeOffset); err != nil {
			fmt.Println("❌ Error al actualizar inodo:", err)
			return
		}

		// Crear nuevo bloque con el usuario
		var newBlock Partitions.Fileblock
		newBlock.B_content = [64]byte{}
		copy(newBlock.B_content[:], updatedUserLine)

		newBlockOffset := int64(tempSuperblock.S_block_start + newBlockIndex*int32(blockSize))
		if err := Utils.EscribirArchivo(file, newBlock, newBlockOffset); err != nil {
			fmt.Println("❌ Error al escribir nuevo bloque:", err)
			return
		}

		fmt.Println("✅ Nuevo bloque creado para usuario cambiado de grupo.")
	}

	fmt.Println("✅ Usuario movido al nuevo grupo correctamente.")
	fmt.Println("╚════════════════════    FIN CHGRP  ═══════════════════════════╝")
}

func freeBlock(file *os.File, sb Partitions.Superblock, blockIndex int32) error {
	// Leer el bitmap de bloques
	bitmapSize := sb.S_blocks_count
	bitmap := make([]byte, bitmapSize)
	if _, err := file.Seek(int64(sb.S_bm_block_start), 0); err != nil {
		return err
	}
	if err := binary.Read(file, binary.LittleEndian, bitmap); err != nil {
		return err
	}

	// Marcar el bloque como libre
	if blockIndex >= 0 && blockIndex < int32(bitmapSize) {
		bitmap[blockIndex] = 0
	} else {
		return fmt.Errorf("índice de bloque inválido")
	}

	// Escribir el bitmap actualizado
	if _, err := file.Seek(int64(sb.S_bm_block_start), 0); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, bitmap); err != nil {
		return err
	}

	return nil
}
