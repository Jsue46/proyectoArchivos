package FSystem

import (
	"Proyecto/Environment"
	"Proyecto/Partitions"
	"Proyecto/User_Groups"
	"Proyecto/Utils"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

func Mkfs(id string, type_ string, fs_ string) {
	fmt.Println("╔═════════════════════ INICIO MKFS   ══════════════════════════╗")
	fmt.Println("  Id:", id)
	fmt.Println("  Type:", type_)
	fmt.Println("  Fs:", fs_)

	// Obtener la fecha actual en formato día/mes/año
	currentDate := time.Now().Format("02/01/2006")

	// Buscar la partición montada por ID
	var mountedPartition Environment.MountedPartition
	var partitionFound bool

	for _, partitions := range Environment.GetMountedPartitions() {
		for _, partition := range partitions {
			if partition.MountID == id {
				mountedPartition = partition
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Println("Particion no encontrada")
		return
	}

	// Verificar si la partición está montada
	if mountedPartition.MountStatus != '1' {
		fmt.Println("La particion aun no esta montada")
		return
	}

	// Abrir el archivo binario de la partición
	file, err := Utils.AbrirArchivo(mountedPartition.MountPath)
	if err != nil {
		return
	}

	// Leer el MBR (Master Boot Record) del archivo binario
	var TempMBR Partitions.MBR
	if err := Utils.LeerArchivo(file, &TempMBR, 0); err != nil {
		return
	}

	// Buscar la partición dentro del MBR que coincida con el ID proporcionado
	var index int = -1
	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].PartSize != 0 {
			if strings.Contains(string(TempMBR.MbrPartitions[i].PartID[:]), id) {
				index = i
				break
			}
		}
	}

	if index == -1 {
		fmt.Println("Particion no encontrada (2)")
		return
	}

	// Calcular el número de inodos basado en el tamaño de la partición
	numerador := int32(TempMBR.MbrPartitions[index].PartSize - int32(binary.Size(Partitions.Superblock{})))
	denominador_base := int32(4 + int32(binary.Size(Partitions.Inode{})) + 3*int32(binary.Size(Partitions.Fileblock{})))
	var temp int32 = 0
	if fs_ == "2fs" {
		temp = 0
	} else {
		fmt.Print("Error por el momento solo está disponible 2FS.")
	}
	denominador := denominador_base + temp
	n := int32(numerador / denominador)

	fmt.Println("INODOS:", n)

	// Crear el Superblock con los campos calculados
	var newSuperblock Partitions.Superblock
	newSuperblock.S_filesystem_type = 2 // EXT2
	newSuperblock.S_inodes_count = n
	newSuperblock.S_blocks_count = 3 * n
	newSuperblock.S_free_blocks_count = 3*n - 2
	newSuperblock.S_free_inodes_count = n - 2
	copy(newSuperblock.S_mtime[:], currentDate)
	copy(newSuperblock.S_umtime[:], currentDate)
	newSuperblock.S_mnt_count = 1
	newSuperblock.S_magic = 0xEF53
	newSuperblock.S_inode_size = int32(binary.Size(Partitions.Inode{}))
	newSuperblock.S_block_size = int32(binary.Size(Partitions.Fileblock{}))

	// Calcular las posiciones de inicio de los bloques en el disco
	newSuperblock.S_bm_inode_start = TempMBR.MbrPartitions[index].PartStart + int32(binary.Size(Partitions.Superblock{}))
	newSuperblock.S_bm_block_start = newSuperblock.S_bm_inode_start + n
	newSuperblock.S_inode_start = newSuperblock.S_bm_block_start + 3*n
	newSuperblock.S_block_start = newSuperblock.S_inode_start + n*newSuperblock.S_inode_size

	// Crear el sistema de archivos EXT2
	if fs_ == "2fs" {
		create_ext2(n, TempMBR.MbrPartitions[index], newSuperblock, currentDate, file)
	} else {
		fmt.Println("EXT3 no está soportado.")
	}

	defer file.Close()

	fmt.Println("╚══════════════════════   FIN MKFS   ══════════════════════════╝")
}

// Función para leer y mostrar el contenido de users.txt
func PrintUsersTxt(file *os.File, tempSuperblock Partitions.Superblock) {
	fmt.Println("======Contenido de users.txt======")

	// Buscar el archivo "users.txt" en el sistema de archivos
	indexInode := User_Groups.InitSearch("/users.txt", file, tempSuperblock)
	if indexInode == -1 {
		fmt.Println("Error: No se encontró el archivo users.txt.")
		return
	}

	// Leer el inodo del archivo "users.txt"
	var crrInode Partitions.Inode
	if err := Utils.LeerArchivo(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Partitions.Inode{})))); err != nil {
		fmt.Println("Error al leer el inodo de users.txt:", err)
		return
	}

	// Obtener el contenido del archivo "users.txt"
	data := GetInodeFileData(crrInode, file, tempSuperblock)

	// Mostrar el contenido de users.txt
	fmt.Println(data)
	fmt.Println("======Fin del contenido de users.txt======")
}
