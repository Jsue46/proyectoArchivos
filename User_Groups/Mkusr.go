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

func MKUSR(user string, pass string, grp string) {
	fmt.Println("╔══════════════════════ INICIO MKUSR  ═════════════════════════╗")
	fmt.Println("  Usuario:", user)
	fmt.Println("  Contraseña:", pass)
	fmt.Println("  Grupo:", grp)

	// Verificar si el usuario actual es root
	if !IsRootLoggedIn() {
		fmt.Println("❌ Error: Solo el usuario root puede crear usuarios.")
		return
	}

	// Validaciones de longitud
	if len(user) > 10 || len(pass) > 10 || len(grp) > 10 {
		fmt.Println("❌ Error: user, pass y grp deben tener máximo 10 caracteres.")
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

	// Variables para el grupo y validación de usuario
	groupID := -1
	userExists := false
	blockSize := binary.Size(Partitions.Fileblock{})

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

		// Procesar el contenido del bloque
		content := string(fileBlock.B_content[:])
		lines := strings.Split(content, "\n")

		// Buscar el grupo y verificar usuario en este bloque
		for _, line := range lines {
			if line == "" {
				continue
			}

			parts := strings.Split(line, ",")
			if len(parts) >= 3 {
				// Buscar el grupo
				if strings.TrimSpace(parts[1]) == "G" && strings.TrimSpace(parts[2]) == grp {
					if strings.TrimSpace(parts[0]) == "0" {
						fmt.Println("❌ Error: El grupo está eliminado, no se pueden agregar usuarios.")
						return
					}
					// Convertir ID del grupo a número
					if id, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
						groupID = id
					}
				}

				// Verificar si el usuario ya existe
				if len(parts) >= 5 && strings.TrimSpace(parts[1]) == "U" && strings.TrimSpace(parts[3]) == user {
					userExists = true
				}
			}
		}
	}

	// Validaciones después de revisar todos los bloques
	if groupID == -1 {
		fmt.Println("❌ Error: El grupo no existe en el sistema.")
		return
	}

	if userExists {
		fmt.Println("❌ Error: El usuario ya existe en el sistema.")
		return
	}

	// Crear la nueva línea del usuario
	newUserLine := fmt.Sprintf("%d,U,%s,%s,%s\n", groupID, grp, user, pass)

	// Agregar el usuario de manera persistente en users.txt
	if err := AppendToFileBlock(&crrInode, newUserLine, file, tempSuperblock); err != nil {
		fmt.Println("❌ Error al escribir en el archivo users.txt:", err)
		return
	}

	fmt.Println("✅ Usuario creado correctamente.")
	fmt.Println("╚══════════════════════   FIN MKUSR   ═════════════════════════╝")
}
