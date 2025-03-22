package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// Obtener el nombre del host
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		// Si ocurre un error, intentar obtener el hostname mediante el comando "hostname"
		cmd := exec.Command("hostname")
		output, err := cmd.Output()
		if err != nil {
			return "Host desconocido"
		}
		return strings.TrimSpace(string(output))
	}
	return hostname
}

// Estructura con campo exportado
type HostName struct {
	Hostname string
}

func Index(rw http.ResponseWriter, r *http.Request) {
	hostnameS := getHostname()
	hostName := HostName{Hostname: hostnameS}

	fmt.Println("Hostname:", hostnameS) // Imprime correctamente en consola

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(rw, "Error cargando plantilla", http.StatusInternalServerError)
		fmt.Println("Error al cargar la plantilla:", err)
		return
	}

	err = tmpl.Execute(rw, hostName)
	if err != nil {
		http.Error(rw, "Error ejecutando plantilla", http.StatusInternalServerError)
		fmt.Println("Error al ejecutar plantilla:", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <puerto>")
		return
	}
	port := os.Args[1]

	http.HandleFunc("/", Index)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	fmt.Println("El servidor se corre en el puerto", port)
	fmt.Printf("Run Server: http://localhost:%s/\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
