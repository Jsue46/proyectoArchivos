package FSystem

import (
	"Proyecto/Partitions"
	"fmt"
)

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

func PrintInode(inode Partitions.Inode) {
	fmt.Println("=== Inode ===")
	fmt.Println("UID:", inode.I_uid)
	fmt.Println("GID:", inode.I_gid)
	fmt.Println("Size:", inode.I_size)
	fmt.Println("Atime:", string(inode.I_atime[:]))
	fmt.Println("Ctime:", string(inode.I_ctime[:]))
	fmt.Println("Mtime:", string(inode.I_mtime[:]))
	fmt.Println("Perm:", string(inode.I_perm[:]))
	fmt.Println("Blocks:", inode.I_block)
}
