package opciones

import "fmt"

var usuarios []string

func IngresarUsuario() {
	var usuario string
	fmt.Print("Ingrese Su Nombre: ")
	fmt.Scanln(&usuario)

	usuarios = append(usuarios, usuario)
	fmt.Println("Usuario Agregado Correctamente")
}
