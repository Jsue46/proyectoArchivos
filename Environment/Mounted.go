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
	fmt.Println("Particiones montadas:")

	// Si no hay particiones montadas, muestra un mensaje y termina la función.
	if len(mountedPartitions) == 0 {
		fmt.Println("No hay particiones montadas.")
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
			fmt.Printf(" - Partición Name: %s, ID: %s, Path: %s, Status: %c, LoggedIn: %s\n",
				partition.MountName, partition.MountID, partition.MountPath, partition.MountStatus, loginStatus)
		}
	}
	fmt.Println("")
}

// Función para obtener las particiones montadas
func GetMountedPartitions() map[string][]MountedPartition {
	return mountedPartitions
}

// MarkPartitionAsLoggedIn busca una partición por su ID y la marca como logueada (LoggedIn = true).
func MarkPartitionAsLoggedIn(id string) {
	// Recorre todas las particiones montadas en los discos.
	for diskID, partitions := range mountedPartitions {
		for i, partition := range partitions {
			// Si la partición coincide con el ID buscado, se marca como logueada.
			if partition.MountID == id {
				mountedPartitions[diskID][i].LoggedIn = true
				fmt.Printf("Partición con ID %s marcada como logueada.\n", id)
				return
			}
		}
	}
	// Si no se encuentra la partición, se muestra un mensaje de error.
	fmt.Printf("No se encontró la partición con ID %s para marcarla como logueada.\n", id)
}

func MarkPartitionAsLoggedOut(id string) {
	mountedPartitions := GetMountedPartitions()

	// Buscar la partición por su ID y marcar como deslogueada
	for _, partitions := range mountedPartitions {
		for i, partition := range partitions {
			if partition.MountID == id {
				partitions[i].LoggedIn = false
				fmt.Printf("Partición con ID %s deslogueada correctamente.\n", id)
				return
			}
		}
	}

	// Si no se encuentra la partición, mostrar un mensaje
	fmt.Printf("No se encontró la partición con ID %s para desloguearla.\n", id)
}
