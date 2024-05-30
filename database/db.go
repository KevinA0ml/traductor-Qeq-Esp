package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	server   = "DESKTOP-8HU9MFB"
	port     = 1433
	user     = "sa"
	password = "123456"
	database = "traductor"
)
var db *sql.DB

func InitDB() {
	connString := fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s",
		server, port, user, password, database)

	var err error
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creando la conexión: ", err.Error())
	}

	// Probar la conexión
	err = db.Ping()
	if err != nil {
		log.Fatal("Error conectándose a la base de datos: ", err.Error())
	}
	fmt.Println("Conexión exitosa a SQL Server")
}

// TRADUCIR TEXTO
func TranslateText(text string, direction string) (string, error) {
	var query string
	switch direction {
	case "es_to_qeqchi":
		query = "SELECT Qeqchi FROM dbo.EQ WHERE REPLACE(LOWER(Espanol), 'á', 'a') = @p1"
	case "qeq_to_es":
		query = "SELECT Espanol FROM dbo.EQ WHERE REPLACE(LOWER(Qeqchi), 'á', 'a') = @p1"
	case "fes_to_qeqchi":
		query = "SELECT frase_Q FROM dbo.Frases WHERE REPLACE(LOWER(frase_Es), 'á', 'a') = @p1"
	case "fqeq_to_es":
		query = "SELECT frase_Es FROM dbo.Frases WHERE REPLACE(LOWER(frase_Q), 'á', 'a') = @p1"
	default:
		return "", fmt.Errorf("dirección de traducción no soportada")
	}

	var translation string
	err := db.QueryRow(query, text).Scan(&translation)
	if err != nil {
		return "", err
	}
	return translation, nil
}

// BANCO DE PALABRAS
type Word struct {
	Espanol string `json:"Espanol"`
	Qeqchi  string `json:"Qeqchi"`
}

func GetWordBank() ([]Word, error) {
	query := "SELECT Espanol, Qeqchi FROM dbo.EQ"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var word Word
		err := rows.Scan(&word.Espanol, &word.Qeqchi)
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}
	return words, nil
}

// Cerrar la conexión a la base de datos
func CloseDB() {
	db.Close()
}
