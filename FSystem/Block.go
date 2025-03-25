package FSystem

import (
	"Proyecto/Partitions"
	"fmt"
)

// Funcion para imprimir el contenido de un bloque de carpetas
func PrintFolderblock(folderblock Partitions.Folderblock) {
	fmt.Println(" ═════════════════════  FOLDERBLOCK ══════════════════════════ ")
	for i, content := range folderblock.B_content {
		fmt.Printf("  Content %d: Name: %s, Inodo: %d\n", i, string(content.B_name[:]), content.B_inodo)
	}
	fmt.Println(" ══════════════════════════════════════════════════════════════ ")
}

// Funcion para imprimir el contenido de un bloque de archivos
func PrintFileblock(fileblock Partitions.Fileblock) {
	fmt.Println(" ══════════════════════   FILEBLOCK  ══════════════════════════ ")
	fmt.Printf("  B_content: %s\n", string(fileblock.B_content[:]))
	fmt.Println(" ══════════════════════════════════════════════════════════════ ")
}
