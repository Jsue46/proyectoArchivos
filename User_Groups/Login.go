package User_Groups

import (
	"Proyecto/Environment"
	"Proyecto/Partitions"
	"Proyecto/Utils"
	"encoding/binary"
	"fmt"
	"strings"
)

// Variable para almacenar el usuario actual
var currentUser string

// Función para verificar si el usuario actual es root
func IsRootLoggedIn() bool {
	return currentUser == "root"
}

// Función para establecer el usuario actual
func SetCurrentUser(user string) {
	currentUser = user
}

// Funcion para Iniciar Sesión
func Login(user string, pass string, id string) {
	fmt.Println("╔═════════════════════ INICIAR SESION ═════════════════════════╗")
	fmt.Printf("  USUARIO: %s\n", user)
	fmt.Printf("  CONTRASEÑA: %s\n", pass)
	fmt.Printf("  ID: %s\n", id)

	mountedPartitions := Environment.GetMountedPartitions()
	var filepath string
	var partitionFound bool
	var login bool = false

	// Buscar la partición montada
	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.MountID == id {
				if partition.LoggedIn {
					fmt.Println("╔═════════════════════════════════════════╗")
					fmt.Println("║       ⚠️ Error: Ya hay Usuario Activo   ║")
					fmt.Println("╚═════════════════════════════════════════╝")
					return
				}
				filepath = partition.MountPath
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Println("╔═════════════════════════════════════════╗")
		fmt.Println("║ ⚠️ Error: Partición No Fue Encontrada   ║")
		fmt.Println("╚═════════════════════════════════════════╝")
		return
	}

	// Abrir el archivo del sistema de archivos binario
	file, err := Utils.AbrirArchivo(filepath)
	if err != nil {
		fmt.Println("╔═════════════════════════════════════════╗")
		fmt.Printf("║ ⚠️ Error: No se pudo Abrir el Archivo: %s║\n", err)
		fmt.Println("╚═════════════════════════════════════════╝")
		return
	}
	defer file.Close()

	// Leer el MBR
	var TempMBR Partitions.MBR
	if err := Utils.LeerArchivo(file, &TempMBR, 0); err != nil {
		fmt.Println("╔═════════════════════════════════════════╗")
		fmt.Printf("║ ⚠️ Error: No se pudo Leer el MBR: %s    ║\n", err)
		fmt.Println("╚═════════════════════════════════════════╝")
		return
	}

	// Buscar la partición dentro del MBR
	var index int = -1
	for i := 0; i < 4; i++ {
		if TempMBR.MbrPartitions[i].PartSize != 0 {
			if strings.Contains(string(TempMBR.MbrPartitions[i].PartID[:]), id) {
				if TempMBR.MbrPartitions[i].PartStatus[0] == '1' {
					index = i
				} else {
					fmt.Println("╔═════════════════════════════════════════╗")
					fmt.Println("║ ⚠️ Error: La Partición No Está Montada  ║")
					fmt.Println("╚═════════════════════════════════════════╝")
					return
				}
				break
			}
		}
	}

	if index == -1 {
		fmt.Println("╔═════════════════════════════════════════╗")
		fmt.Println("║ ⚠️ Error: No se Encontró La Partición    ║")
		fmt.Println("╚═════════════════════════════════════════╝")
		return
	}

	// Leer el Superblock
	var tempSuperblock Partitions.Superblock
	if err := Utils.LeerArchivo(file, &tempSuperblock, int64(TempMBR.MbrPartitions[index].PartStart)); err != nil {
		fmt.Println("╔═════════════════════════════════════════╗")
		fmt.Printf("║ ⚠️ Error: No se Pudo Leer el SuperBlock: %s ║\n", err)
		fmt.Println("╚═════════════════════════════════════════╝")
		return
	}

	// Buscar el archivo users.txt
	indexInode := InitSearch("/users.txt", file, tempSuperblock)
	if indexInode == -1 {
		fmt.Println("╔═══════════════════════════════════════════════╗")
		fmt.Println("║ ⚠️ Error: Archivo '/users.txt' No Encontrado  ║")
		fmt.Println("╚═══════════════════════════════════════════════╝")
		return
	}

	// Leer el Inodo del archivo "users.txt"
	var crrInode Partitions.Inode
	if err := Utils.LeerArchivo(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Partitions.Inode{})))); err != nil {
		fmt.Println("╔═════════════════════════════════════════╗")
		fmt.Println("║ ⚠️ Error: No se Pudo Leer el Inodo      ║")
		fmt.Println("╚═════════════════════════════════════════╝")
		return
	}

	// Obtener el contenido del archivo users.txt
	data := GetInodeFileData(crrInode, file, tempSuperblock)
	lines := strings.Split(data, "\n")

	// Verificar credenciales
	for _, line := range lines {
		words := strings.Split(line, ",")
		if len(words) == 5 {
			if words[3] == user && words[4] == pass {
				login = true
				break
			}
		}
	}

	// Si el login es exitoso, marcar la partición como logueada
	if login {
		fmt.Printf(" ═══════════  Usuario '%s' Logueado con Exito  ═════════════════", user)
		Environment.ParticionConInicioSesion(id)
		SetCurrentUser(user)
	} else {
		fmt.Println("╔════════════════════════════════════════════╗")
		fmt.Println("║ ⚠️ Error: Usuario o Contraseña Incorrectos ║")
		fmt.Println("╚════════════════════════════════════════════╝")
	}

	fmt.Println("╚═════════════════   FIN INICIO SESION  ═══════════════════════╝")
}
