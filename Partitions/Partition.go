package Partitions

import (
	"fmt"
)

// Definición de la estructura de una partición
type Partition struct {
	PartStatus      [1]byte  // 1 byte
	PartType        [1]byte  // 1 byte (P: Primaria, E: Extendida)
	PartFit         [1]byte  // 1 byte (B: Best, F: First, W: Worst)
	PartStart       int32    // 4 bytes
	PartSize        int32    // 4 bytes
	PartName        [16]byte // 16 bytes
	PartCorrelative int32    // 4 bytes
	PartID          [4]byte  // 4 bytes
}

// Definición de la estructura del MBR
type MBR struct {
	MbrTamanio       int32        // 4 bytes
	MbrFechaCreacion [19]byte     // 19 bytes (fecha en formato string)
	MbrDskSignature  int32        // 4 bytes (número único)
	DskFit           [1]byte      // 1 byte (B: Best, F: First, W: Worst)
	MbrPartitions    [4]Partition // 4 particiones de 36 bytes cada una
}

// Definición de la estructura de un EBR
type EBR struct {
	PartStatus byte     // 1 byte
	PartFit    byte     // 1 byte
	PartStart  int32    // 4 bytes
	PartSize   int32    // 4 bytes
	PartNext   int32    // 4 bytes
	PartName   [16]byte // 16 bytes
	PartMount  byte     // 1 byte
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

// Función para imprimir una partición
func ImprimirPartition(data Partition) {
	fmt.Printf(
		"Name: %s, type: %s, start: %d, size: %d, status: %s, id: %s\n",
		string(data.PartName[:]), string(data.PartType[:]), data.PartStart, data.PartSize,
		string(data.PartStatus[:]), string(data.PartID[:]))
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

//Estructuras relacionadas a EXT2

type Superblock struct {
	S_filesystem_type   int32    // Guarda el número que identifica el sistema de archivos utilizado
	S_inodes_count      int32    // Guarda el número total de inodos
	S_blocks_count      int32    // Guarda el número total de bloques
	S_free_blocks_count int32    // Contiene el número de bloques libres
	S_free_inodes_count int32    // Contiene el número de inodos libres
	S_mtime             [17]byte // Última fecha en el que el sistema fue montado
	S_umtime            [17]byte // Última fecha en que el sistema fue desmontado
	S_mnt_count         int32    // Indica cuantas veces se ha montado el sistema
	S_magic             int32    // Valor que identifica al sistema de archivos, tendrá el valor 0xEF53
	S_inode_size        int32    // Tamaño del inodo
	S_block_size        int32    // Tamaño del bloque
	S_fist_ino          int32    // Primer inodo libre (dirección del inodo)
	S_first_blo         int32    // Primer bloque libre (dirección del inodo)
	S_bm_inode_start    int32    // Guardará el inicio del bitmap de inodos
	S_bm_block_start    int32    // Guardará el inicio del bitmap de bloques
	S_inode_start       int32    // Guardará el inicio de la tabla de inodos
	S_block_start       int32    // Guardará el inicio de la tabla de bloques
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
	fmt.Printf("S_magic: 0x%X\n", sb.S_magic) // Usamos 0x%X para mostrarlo en formato hexadecimal
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
	B_pointers [16]int32 //Este campo es un apuntador al inodo asociado al archivo o carpeta.
	//Cada apuntador es un int32 (4 bytes), lo que permite direccionar hasta 2^32 bloques
}

func PrintPointerblock(pointerblock Pointerblock) {
	fmt.Println("====== Pointerblock ======")
	for i, pointer := range pointerblock.B_pointers {
		fmt.Printf("Pointer %d: %d\n", i, pointer)
	}
	fmt.Println("=========================")
}
