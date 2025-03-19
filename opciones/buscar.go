package opciones

import "fmt"

func BuscarUsuario(nombre string) {
	for _, usuario := range usuarios {
		if usuario == nombre {
			fmt.Println("El usuario", nombre, "existe en el sistema.")
			return
		}
	}
	fmt.Println("El usuario", nombre, "No existe en el sistema.")
}
