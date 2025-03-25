package User_Groups

import (
	"Proyecto/Environment"
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"encoding/binary"
	"fmt"
	"strings"
)

func RMGRP(name string) {
	fmt.Println("╔══════════════════════ INICIO RMGRP  ═════════════════════════╗")
	fmt.Println("  Nombre del grupo a eliminar:", name)

	if !IsRootLoggedIn() {
		fmt.Println("❌ Error: Solo el usuario root puede eliminar grupos.")
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

		lines := strings.Split(string(fileBlock.B_content[:]), "\n")
		for i, line := range lines {
			parts := strings.Split(line, ",")
			if len(parts) >= 3 && strings.TrimSpace(parts[1]) == "G" && strings.TrimSpace(parts[2]) == name {
				fmt.Println("✅ Grupo encontrado, marcando como eliminado en su bloque.")
				lines[i] = "0,G," + parts[2]
				fileBlock.B_content = [64]byte{}
				copy(fileBlock.B_content[:], strings.Join(lines, "\n"))

				if err := Utils.EscribirArchivo(file, fileBlock, blockOffset); err != nil {
					fmt.Println("❌ Error al escribir el bloque actualizado:", err)
					return
				}

				fmt.Println("✅ Grupo eliminado correctamente en su bloque.")
				fmt.Println("╚══════════════════════   FIN RMGRP   ═════════════════════════╝")
				return
			}
		}
	}

	fmt.Println("❌ Error: El grupo no existe en el sistema.")
}
