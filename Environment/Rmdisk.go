package Environment

import (
	"fmt"
	"os"
)

// Rmdisk elimina un archivo de disco en la ruta especificada
func Rmdisk(path string) {
	fmt.Println("╔═════════════════════ INICIO RMDISK  ═════════════════════════╗")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Error: El archivo no existe en la ruta especificada: %s\n", path)
		return
	}

	// Eliminar el archivo
	err := os.Remove(path)
	if err != nil {
		fmt.Printf("Error: No se pudo eliminar el disco en la ruta %s: %v\n", path, err)
		return
	}

	fmt.Printf("  Disco En la Ruta %s .\n", path)
	fmt.Printf("    Ha Siso Eliminado Correctamente.\n")

	fmt.Println("╚═════════════════════   FIN RMDISK   ═════════════════════════╝")
}
