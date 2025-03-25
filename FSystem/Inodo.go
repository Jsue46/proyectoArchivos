package FSystem

import (
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"encoding/binary"
	"fmt"
	"os"
)

// Función auxiliar para inicializar un inodo
func initInode(inode *Partitions.Inode, date string) {
	inode.I_uid = 1
	inode.I_gid = 1
	inode.I_size = 0
	copy(inode.I_atime[:], date)
	copy(inode.I_ctime[:], date)
	copy(inode.I_mtime[:], date)
	copy(inode.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		inode.I_block[i] = -1
	}
}

// GetInodeFileData lee el contenido de un archivo a partir de su inodo.
func GetInodeFileData(inode Partitions.Inode, file *os.File, superblock Partitions.Superblock) string {
	fmt.Println(" ═══════════════   CONTENIDO DEL BLOQUE   ═════════════════════ ")

	index := int32(0)
	var content string

	// Iterar sobre los bloques del inodo
	for _, block := range inode.I_block {
		if block != -1 {
			// Manejo de bloques directos (0-12)
			if index < 13 {
				var crrFileblock Partitions.Fileblock

				// Leer el bloque de archivo desde el archivo binario
				if err := Utils.LeerArchivo(file, &crrFileblock, int64(superblock.S_block_start+block*int32(binary.Size(Partitions.Fileblock{})))); err != nil {
					return ""
				}

				// Concatenar el contenido del bloque
				content += string(crrFileblock.B_content[:])
			} else {
				fmt.Print("indirectos") // Falta implementar acceso a bloques indirectos
			}
		}
		index++
	}

	fmt.Println(" ═══════════════  FIN CONTENIDO DEL BLOQUE  ════════════════════ ")
	return content
}

func PrintInode(inode Partitions.Inode) {
	fmt.Println(" ═══════════════════════   INODE   ════════════════════════════ ")
	fmt.Printf(" I_uid: %d\n", inode.I_uid)
	fmt.Printf(" I_gid: %d\n", inode.I_gid)
	fmt.Printf(" I_size: %d\n", inode.I_size)
	fmt.Printf(" I_atime: %s\n", string(inode.I_atime[:]))
	fmt.Printf(" I_ctime: %s\n", string(inode.I_ctime[:]))
	fmt.Printf(" I_mtime: %s\n", string(inode.I_mtime[:]))
	fmt.Printf(" I_type: %s\n", string(inode.I_type[:]))
	fmt.Printf(" I_perm: %s\n", string(inode.I_perm[:]))
	fmt.Printf(" I_block: %v\n", inode.I_block)
	fmt.Println(" ══════════════════════════════════════════════════════════════ ")
}
