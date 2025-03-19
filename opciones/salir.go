package opciones

import (
	"fmt"
	"os"
)

func OpcionSalida() {
	fmt.Println("Saliendo del programa...")
	os.Exit(0)
}
