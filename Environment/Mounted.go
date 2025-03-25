package Environment

import "fmt"

// Mapa para almacenar las particiones montadas
var mountedPartitions = make(map[string][]MountedPartition)

// Estructura para almacenar las particiones montadas
type MountedPartition struct {
	MountPath   string
	MountName   string
	MountID     string
	MountStatus byte
	LoggedIn    bool
	MountStart  int32
}

func ImprimirMountedPartitions() {
	fmt.Println("  Particiones Montadas:")

	// Si no hay particiones montadas, muestra un mensaje y termina la función.
	if len(mountedPartitions) == 0 {
		fmt.Println("No hay Particiones Montadas.")
		return
	}

	// Itera sobre cada disco montado y sus particiones.
	for diskID, partitions := range mountedPartitions {
		fmt.Printf("Disco ID: %s\n", diskID)
		for _, partition := range partitions {
			// Determina si la partición está logueada o no.
			loginStatus := "No"
			if partition.LoggedIn {
				loginStatus = "Sí"
			}
			// Imprime los detalles de la partición.
			fmt.Printf("  - Partición Name: %s, ID: %s, Path: %s, Status: %c, LoggedIn: %s\n",
				partition.MountName, partition.MountID, partition.MountPath, partition.MountStatus, loginStatus)
		}
	}
	fmt.Println("")
}

// Función para obtener las particiones montadas
func GetMountedPartitions() map[string][]MountedPartition {
	return mountedPartitions
}

// Función para marcar una partición con inicio de sesión
func ParticionConInicioSesion(id string) {
	for diskID, partitions := range mountedPartitions {
		for i, partition := range partitions {
			// Si la partición coincide con el ID buscado, se inicia sesión.
			if partition.MountID == id {
				mountedPartitions[diskID][i].LoggedIn = true
				fmt.Printf("\t Partición con ID %s encontrada.\n", id)
				fmt.Println("\t Inicio de sesión exitoso.")
				return
			}
		}
	}
	// Si no se encuentra la partición, se muestra un mensaje de error.
	fmt.Printf("No se encontró la partición con ID %s para marcarla con Inicio de Sesión.\n", id)
}

// Función para marcar una partición como sin sesión iniciada
func ParticionSinInicioSesion(id string) {
	mountedPartitions := GetMountedPartitions()

	// Buscar la partición por su ID y marcar como sin sesión iniciada
	for _, partitions := range mountedPartitions {
		for i, partition := range partitions {
			if partition.MountID == id {
				partitions[i].LoggedIn = false
				fmt.Printf("Partición con ID %s marcada como sin inicio de sesión.\n", id)
				return
			}
		}
	}

	// Si no se encuentra la partición, mostrar un mensaje
	fmt.Printf("No se encontró la partición con ID %s para marcarla como sin inicio de sesión.\n", id)
}

// Función para obtener la partición montada actualmente
func GetCurrentMountedPartition() *MountedPartition {
	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.LoggedIn {
				return &partition
			}
		}
	}
	return nil
}
