package User_Groups

import (
	"Proyecto/Environment"
	"fmt"
)

// Función para cerrar sesión
func Logout() {
	fmt.Println("╔═════════════════════ CERRAR  SESION ═════════════════════════╗")
	// Obtener las particiones montadas
	mountedPartitions := Environment.GetMountedPartitions()
	// Variable para verificar si hay una sesión activa
	var sessionFound bool = false

	// Buscar la partición que tiene una sesión activa
	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.LoggedIn {
				// Marcar la partición como deslogueada
				Environment.ParticionSinInicioSesion(partition.MountID)
				sessionFound = true
				break
			}
		}
		if sessionFound {
			break
		}
	}

	// Si no se encontró una sesión activa, mostrar un error
	if !sessionFound {
		fmt.Println("  Error: No hay una Sesión Activa para Cerrar.")
	} else {
		fmt.Println("\t Sesión Finalizada  con éxito.")
	}

	fmt.Println("╚═══════════════════   CERRANDO SESION  ═══════════════════════╝")
}
