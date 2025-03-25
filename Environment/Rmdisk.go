/*
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

<<<<<<< Updated upstream
	fmt.Printf("  Disco En la Ruta %s .\n", path)
	fmt.Printf("    Ha Siso Eliminado Correctamente.\n")

	fmt.Println("╚═════════════════════   FIN RMDISK   ═════════════════════════╝")
=======
	fmt.Printf("Disco en la ruta %s eliminado correctamente.\n", path)
}*/

package Environment

import (
	"fmt"
	"os"
)

// Rmdisk elimina un archivo de disco en la ruta especificada
func Rmdisk(path string) {
	// Verificar si la ruta está vacía
	if path == "" {
		fmt.Println("Error: La ruta del archivo no puede estar vacía.")
		return
	}

	// Verificar si el archivo existe
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("Error: El archivo no existe en la ruta especificada: %s\n", path)
		return
	} else if err != nil {
		fmt.Printf("Error: No se pudo verificar la existencia del archivo en %s: %v\n", path, err)
		return
	}

	// Verificar si es un archivo regular (no un directorio u otro tipo de archivo)
	if !fileInfo.Mode().IsRegular() {
		fmt.Printf("Error: La ruta especificada no es un archivo regular: %s\n", path)
		return
	}
	/*
		// Confirmar la eliminación con el usuario
		fmt.Printf("¿Está seguro de que desea eliminar el disco en la ruta %s? (s/n): ", path)
		var confirmacion string
		fmt.Scanln(&confirmacion)

		if confirmacion != "s" && confirmacion != "S" {
			fmt.Println("Eliminación cancelada.")
			return
		}
	*/
	// Eliminar el archivo
	err = os.Remove(path)
	if err != nil {
		fmt.Printf("Error: No se pudo eliminar el disco en la ruta %s: %v\n", path, err)
		return
	}
	Log("======INICIO RMDISK======")
	fmt.Printf("Disco en la ruta %s eliminado correctamente.\n", path)
	Log("======FIN RMDISK======")
>>>>>>> Stashed changes
}
