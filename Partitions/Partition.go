package Partitions

import (
	"fmt"
)

// Definición de la estructura del MBR
type MBR struct {
	MbrTamanio       int32
	MbrFechaCreacion [19]byte
	MbrDskSignature  int32
	DskFit           [1]byte
	MbrPartitions    [4]Partition
}

// Función para imprimir un MBR
func ImprimirMBR(data MBR) {
	fmt.Printf("MbrFechaCreacion: %s, DskFit: %s, Tamaño: %d, MbrDskSignature: %d\n",
		string(data.MbrFechaCreacion[:]),
		string(data.DskFit[:]),
		data.MbrTamanio,
		data.MbrDskSignature)
	for i := 0; i < 4; i++ {
		ImprimirPartition(data.MbrPartitions[i])
	}
}

// Definición de la estructura de una partición
type Partition struct {
	PartStatus      [1]byte
	PartType        [1]byte
	PartFit         [1]byte
	PartStart       int32
	PartSize        int32
	PartName        [16]byte
	PartCorrelative int32
	PartID          [4]byte
}

// Función para imprimir una partición
func ImprimirPartition(data Partition) {
	fmt.Printf(
		"Name: %s, type: %s, start: %d, size: %d, status: %s, id: %s\n",
		string(data.PartName[:]), string(data.PartType[:]), data.PartStart, data.PartSize,
		string(data.PartStatus[:]), string(data.PartID[:]))
}

// Definición de la estructura de un EBR
type EBR struct {
	PartMount byte
	PartFit   byte
	PartStart int32
	PartSize  int32
	PartNext  int32
	PartName  [16]byte
}

// Función para imprimir un EBR
func ImrprimirEBR(data EBR) {
	fmt.Printf("Name: %s, fit: %c, start: %d, size: %d, next: %d, mount: %c\n",
		string(data.PartName[:]),
		data.PartFit,
		data.PartStart,
		data.PartSize,
		data.PartNext,
		data.PartMount)
}

// Definición de la estructura del superbloque
type Superblock struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_blocks_count int32
	S_free_inodes_count int32
	S_mtime             [17]byte
	S_umtime            [17]byte
	S_mnt_count         int32
	S_magic             int32
	S_inode_size        int32
	S_block_size        int32
	S_fist_ino          int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
}

func PrintSuperblock(sb Superblock) {
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

type Inode struct {
	I_uid   int32
	I_gid   int32
	I_size  int32
	I_atime [17]byte
	I_ctime [17]byte
	I_mtime [17]byte
	I_block [15]int32
	I_type  [1]byte
	I_perm  [3]byte
}

func PrintInode(inode Inode) {
	fmt.Println("====== Inode ======")
	fmt.Printf("I_uid: %d\n", inode.I_uid)
	fmt.Printf("I_gid: %d\n", inode.I_gid)
	fmt.Printf("I_size: %d\n", inode.I_size)
	fmt.Printf("I_atime: %s\n", string(inode.I_atime[:]))
	fmt.Printf("I_ctime: %s\n", string(inode.I_ctime[:]))
	fmt.Printf("I_mtime: %s\n", string(inode.I_mtime[:]))
	fmt.Printf("I_type: %s\n", string(inode.I_type[:]))
	fmt.Printf("I_perm: %s\n", string(inode.I_perm[:]))
	fmt.Printf("I_block: %v\n", inode.I_block)
	fmt.Println("===================")
}

type Folderblock struct {
	B_content [4]Content
}

func PrintFolderblock(folderblock Folderblock) {
	fmt.Println("====== Folderblock ======")
	for i, content := range folderblock.B_content {
		fmt.Printf("Content %d: Name: %s, Inodo: %d\n", i, string(content.B_name[:]), content.B_inodo)
	}
	fmt.Println("=========================")
}

type Content struct {
	B_name  [12]byte //Este campo almacena el nombre de un archivo o carpeta
	B_inodo int32    //Se utiliza un arreglo de bytes ([12]byte) para asegurar que el nombre tenga un tamaño fijo y sea fácil de manejar
}

type Fileblock struct {
	B_content [64]byte
}

func PrintFileblock(fileblock Fileblock) {
	fmt.Println("====== Fileblock ======")
	fmt.Printf("B_content: %s\n", string(fileblock.B_content[:]))
	fmt.Println("=======================")
}

type Pointerblock struct {
	B_pointers [16]int32
}

func PrintPointerblock(pointerblock Pointerblock) {
	fmt.Println("====== Pointerblock ======")
	for i, pointer := range pointerblock.B_pointers {
		fmt.Printf("Pointer %d: %d\n", i, pointer)
	}
	fmt.Println("=========================")
}
