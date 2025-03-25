package User_Groups

import (
	"Proyecto/Environment"
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"encoding/binary"
	"fmt"
	"strings"
)

func RMUSR(user string) {
	fmt.Println("╔══════════════════════ INICIO RMUSR  ═════════════════════════╗")
	fmt.Println("  Usuario a eliminar:", user)

	// Verificar si el usuario actual es root
	if !IsRootLoggedIn() {
		fmt.Println("❌ Error: Solo el usuario root puede eliminar usuarios.")
		return
	}

	// Obtener la partición montada actualmente
	mountedPartition := Environment.GetCurrentMountedPartition()
	if mountedPartition == nil {
		fmt.Println("❌ Error: No hay ninguna partición montada.")
		return
	}

	// Abrir el archivo del sistema de archivos
	file, err := Utils.AbrirArchivo(mountedPartition.MountPath)
	if err != nil {
		fmt.Println("❌ Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Leer el Superblock de la partición
	var tempSuperblock Partitions.Superblock
	if err := Utils.LeerArchivo(file, &tempSuperblock, int64(mountedPartition.MountStart)); err != nil {
		fmt.Println("❌ Error al leer el Superblock:", err)
		return
	}

	// Buscar el archivo "users.txt"
	indexInode := InitSearch("/users.txt", file, tempSuperblock)
	if indexInode == -1 {
		fmt.Println("❌ Error: No se encontró el archivo users.txt.")
		return
	}

	// Leer el inodo del archivo "users.txt"
	var crrInode Partitions.Inode
	if err := Utils.LeerArchivo(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Partitions.Inode{})))); err != nil {
		fmt.Println("❌ Error al leer el inodo de users.txt:", err)
		return
	}

	blockSize := binary.Size(Partitions.Fileblock{})
	userFound := false

	// Recorrer todos los bloques del inodo
	for _, block := range crrInode.I_block {
		if block == -1 {
			continue
		}

		var fileBlock Partitions.Fileblock
		blockOffset := int64(tempSuperblock.S_block_start + block*int32(blockSize))

		if err := Utils.LeerArchivo(file, &fileBlock, blockOffset); err != nil {
			fmt.Println("❌ Error al leer bloque de archivo:", err)
			return
		}

		// Convertir el contenido a string y dividir por líneas
		content := string(fileBlock.B_content[:])
		lines := strings.Split(content, "\n")
		modified := false

		// Buscar el usuario en este bloque
		for i, line := range lines {
			if line == "" {
				continue
			}

			parts := strings.Split(line, ",")
			if len(parts) >= 5 && strings.TrimSpace(parts[1]) == "U" && strings.TrimSpace(parts[3]) == user {
				fmt.Println("✅ Usuario encontrado, marcando como eliminado en su bloque.")
				// Marcar el usuario como eliminado (ID = 0)
				lines[i] = "0,U," + strings.TrimSpace(parts[2]) + "," + strings.TrimSpace(parts[3]) + "," + strings.TrimSpace(parts[4])
				modified = true
				userFound = true
				break
			}
		}

		// Si encontramos y modificamos el usuario, actualizar el bloque
		if modified {
			// Unir las líneas y asegurarse de que no exceda el tamaño del bloque
			newContent := strings.Join(lines, "\n")
			if len(newContent) > len(fileBlock.B_content) {
				newContent = newContent[:len(fileBlock.B_content)]
			}

			// Limpiar el bloque y copiar el nuevo contenido
			fileBlock.B_content = [64]byte{}
			copy(fileBlock.B_content[:], newContent)

			// Escribir el bloque modificado en el archivo
			if err := Utils.EscribirArchivo(file, fileBlock, blockOffset); err != nil {
				fmt.Println("❌ Error al escribir el bloque actualizado:", err)
				return
			}

			fmt.Println("✅ Usuario eliminado correctamente.")
			break
		}
	}

	if !userFound {
		fmt.Println("❌ Error: El usuario no existe en el sistema.")
	}

	fmt.Println("╚══════════════════════   FIN RMUSR   ═════════════════════════╝")
}
