/*
package main

import (

	"Proyecto/Analyzer"
	"encoding/json"
	"fmt"
	"net/http"

)

// Estructura para recibir comandos desde el frontend

	type CommandRequest struct {
		Command string `json:"command"`
	}

// Estructura para enviar respuestas al frontend

	type CommandResponse struct {
		Output string `json:"output"`
	}

// Función para manejar las solicitudes del frontend

	func executeCommandHandler(w http.ResponseWriter, r *http.Request) {
		var cmdReq CommandRequest
		if err := json.NewDecoder(r.Body).Decode(&cmdReq); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Aquí llamamos al analizador de comandos
		output := Analyzer.AnalizadorEntrada(cmdReq.Command)

		// Enviamos la respuesta al frontend
		response := CommandResponse{Output: output}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

	func main() {
		// Configuramos el endpoint para recibir comandos
		http.HandleFunc("/execute", executeCommandHandler)

		// Iniciamos el servidor en el puerto 8080
		fmt.Println("Server started at :8080")
		http.ListenAndServe(":8080", nil)
	}
*/
/*
package main

import (
	"Proyecto/Analyzer"
)

func main() {
	Analyzer.AnalizadorEntrada()
}
*/
package main

import (
	"Proyecto/Analyzer"
	"encoding/json"
	"fmt"
	"net/http"
)

// Estructura para recibir comandos desde el frontend
type CommandRequest struct {
	Command string `json:"command"`
}

// Estructura para enviar respuestas al frontend
type CommandResponse struct {
	Output string `json:"output"`
}

func main() {
	// Configuramos el endpoint para recibir comandos
	http.HandleFunc("/execute", executeCommandHandler)

	// Iniciamos el servidor en el puerto 8080
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

// Función para manejar las solicitudes del frontend
func executeCommandHandler(w http.ResponseWriter, r *http.Request) {
	// Habilitar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Manejar solicitudes OPTIONS (preflight)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Decodificar el cuerpo de la solicitud
	var cmdReq CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&cmdReq); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Llamar al analizador de comandos
	output := Analyzer.AnalizadorEntrada(cmdReq.Command)

	// Enviar la respuesta al frontend
	response := CommandResponse{Output: output}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
