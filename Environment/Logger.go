// Logger.go
package Environment

import (
	"fmt"
	"strings"
)

var logBuffer []string // Almacena los mensajes

// Log imprime en la consola y almacena el mensaje en el buffer
func Log(message string) {
	fmt.Println(message)                   // Imprime en la consola
	logBuffer = append(logBuffer, message) // Almacena el mensaje
}

// GetLogs devuelve todos los mensajes almacenados como una cadena
func GetLogs() string {
	return strings.Join(logBuffer, "\n") // Une los mensajes con saltos de línea
}

// ClearLogs limpia el buffer de mensajes
func ClearLogs() {
	logBuffer = nil
}
