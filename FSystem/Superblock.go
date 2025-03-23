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

// Función para imprimir el Superblock
func PrintSuperblock(sb Partitions.Superblock) {
	fmt.Println("====== Superblock ======")
	fmt.Printf("S_filesystem_type: %d\n", sb.S_filesystem_type)
	fmt.Printf("S_inodes_count: %d\n", sb.S_inodes_count)
	fmt.Printf("S_blocks_count: %d\n", sb.S_blocks_count)
	fmt.Printf("S_free_blocks_count: %d\n", sb.S_free_blocks_count)
	fmt.Printf("S_free_inodes_count: %d\n", sb.S_free_inodes_count)
	fmt.Printf("S_mtime: %s\n", string(sb.S_mtime[:]))
	fmt.Printf("S_umtime: %s\n", string(sb.S_umtime[:]))
	fmt.Printf("S_mnt_count: %d\n", sb.S_mnt_count)
	fmt.Printf("S_magic: 0x%X\n", sb.S_magic)
	fmt.Printf("S_inode_size: %d\n", sb.S_inode_size)
	fmt.Printf("S_block_size: %d\n", sb.S_block_size)
	fmt.Printf("S_fist_ino: %d\n", sb.S_fist_ino)
	fmt.Printf("S_first_blo: %d\n", sb.S_first_blo)
	fmt.Printf("S_bm_inode_start: %d\n", sb.S_bm_inode_start)
	fmt.Printf("S_bm_block_start: %d\n", sb.S_bm_block_start)
	fmt.Printf("S_inode_start: %d\n", sb.S_inode_start)
	fmt.Printf("S_block_start: %d\n", sb.S_block_start)
	fmt.Println("========================")
}
