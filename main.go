package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/KevinA0ml/traductorQeqchi-Espanol/database"
)

var tmpl *template.Template

func main() {
	// Inicializar la base de datos
	database.InitDB()
	defer database.CloseDB()

	// Cargar la plantilla HTML
	tmpl = template.Must(template.ParseFiles(filepath.Join("static", "index.html")))

	// Servir archivos estáticos (CSS, JS, imágenes, etc.)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Manejar la solicitud HTTP principal
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error renderizando la plantilla", http.StatusInternalServerError)
		}
	})

	// Ruta para manejar las traducciones
	http.HandleFunc("/translate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			Text      string `json:"text"`
			Direction string `json:"direction"`
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Solicitud inválida", http.StatusBadRequest)
			return
		}

		translation, err := database.TranslateText(request.Text, request.Direction)
		if err != nil {
			http.Error(w, "Error en la traducción", http.StatusInternalServerError)
			return
		}

		response := struct {
			Translation string `json:"translation"`
		}{
			Translation: translation,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Ruta para manejar el banco de palabras
	http.HandleFunc("/wordbank", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		words, err := database.GetWordBank()
		if err != nil {
			http.Error(w, "Error recuperando el banco de palabras", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(words)
	})

	/// aca obtengo las sugerencias

	// Iniciar el servidor en el puerto 8080
	fmt.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
