package Environment

import (
	"fmt"
	"os"
)

// Rmdisk elimina un archivo de disco en la ruta especificada
func Rmdisk(path string) {
	// Verificar si el archivo existe
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

	fmt.Printf("Disco en la ruta %s eliminado correctamente.\n", path)
}
