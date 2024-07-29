package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type Estacao struct {
	ID               int        `json:"id"`
	Estacao          string     `json:"estacao"`
	Cidade           string     `json:"cidade"`
	Coordenadas      string     `json:"coordenadas"`
	InicioDeOperacao *time.Time `json:"iniciodeoperacao"`
	FimDeOperacao    *time.Time `json:"fimdeoperacao"`
	EmUso            bool       `json:"emuso"`
}

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := DatabaseConn()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	estacoes, err := QueryStations(db)

	jsonData, err := json.MarshalIndent(estacoes, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))

}

func DatabaseConn() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbSsl := os.Getenv("DB_SSL")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, dbSsl)

	fmt.Println(connStr)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	return db, err

}

func QueryStations(db *sql.DB) ([]Estacao, error) {
	rows, err := db.Query("SELECT * from estacao")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var estacoes []Estacao

	for rows.Next() {
		var estacao Estacao

		var inicioDeOperacao, fimDeOperacao sql.NullTime

		err := rows.Scan(&estacao.ID, &estacao.Estacao, &estacao.Cidade, &estacao.Coordenadas, &inicioDeOperacao, &fimDeOperacao, &estacao.EmUso)
		if err != nil {
			log.Fatal(err)
		}
		if inicioDeOperacao.Valid {
			estacao.InicioDeOperacao = &inicioDeOperacao.Time
		}
		if fimDeOperacao.Valid {
			estacao.FimDeOperacao = &fimDeOperacao.Time
		}

		estacoes = append(estacoes, estacao)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return estacoes, nil
}
