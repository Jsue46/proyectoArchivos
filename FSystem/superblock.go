package FSystem

import (
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"fmt"
	"os"
)

// Función auxiliar para marcar los inodos y bloques usados
func markUsedInodesAndBlocks(newSuperblock Partitions.Superblock, file *os.File) error {
	if err := Utils.EscribirArchivo(file, byte(1), int64(newSuperblock.S_bm_inode_start)); err != nil {
		return err
	}
	if err := Utils.EscribirArchivo(file, byte(1), int64(newSuperblock.S_bm_inode_start+1)); err != nil {
		return err
	}
	if err := Utils.EscribirArchivo(file, byte(1), int64(newSuperblock.S_bm_block_start)); err != nil {
		return err
	}
	if err := Utils.EscribirArchivo(file, byte(1), int64(newSuperblock.S_bm_block_start+1)); err != nil {
		return err
	}
	return nil
}

func PrintSuperblock(sb Partitions.Superblock) {
	fmt.Println("=== Superblock ===")
	fmt.Println("Filesystem Type:", sb.S_filesystem_type)
	fmt.Println("Inodes Count:", sb.S_inodes_count)
	fmt.Println("Blocks Count:", sb.S_blocks_count)
	fmt.Println("Free Blocks Count:", sb.S_free_blocks_count)
	fmt.Println("Free Inodes Count:", sb.S_free_inodes_count)
	fmt.Println("Mount Time:", string(sb.S_mtime[:]))
	fmt.Println("Unmount Time:", string(sb.S_umtime[:]))
	fmt.Println("Mount Count:", sb.S_mnt_count)
	fmt.Println("Magic:", sb.S_magic)
	fmt.Println("Inode Size:", sb.S_inode_size)
	fmt.Println("Block Size:", sb.S_block_size)
	fmt.Println("Bitmap Inode Start:", sb.S_bm_inode_start)
	fmt.Println("Bitmap Block Start:", sb.S_bm_block_start)
	fmt.Println("Inode Start:", sb.S_inode_start)
	fmt.Println("Block Start:", sb.S_block_start)
}
