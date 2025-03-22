package FSystem

import (
	"Proyecto/Partitions"
	"fmt"
)

func PrintFolderblock(folderblock Partitions.Folderblock) {
	fmt.Println("=== Folderblock ===")
	for i, content := range folderblock.B_content {
		if content.B_inodo != -1 {
			fmt.Printf("Content %d: Inodo=%d, Name=%s\n", i, content.B_inodo, string(content.B_name[:]))
		}
	}
}

func PrintFileblock(fileblock Partitions.Fileblock) {
	fmt.Println("=== Fileblock ===")
	fmt.Println("Content:", string(fileblock.B_content[:]))
}
