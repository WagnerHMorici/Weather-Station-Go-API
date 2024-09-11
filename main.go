package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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

type Data struct {
	ID          int        `json:"id"`
	Temperatura float64    `json:"temperatura"`
	Umidade     float64    `json:"umidade"`
	DataHora    *time.Time `json:"datahora"`
	Estacao_FK  int        `json:estacao_fk`
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

	data, err := QueryDataStations(db)

	router := gin.Default()

	router.GET("/stations", func(c *gin.Context) {
		c.JSON(200, estacoes)
	})

	router.GET("/data", func(c *gin.Context) {
		c.JSON(200, data)
	})

	router.Run()

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

func QueryDataStations(db *sql.DB) ([]Data, error) {
	rows, err := db.Query("SELECT * from registros_estacoes")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var all_data []Data

	for rows.Next() {
		var data Data

		var datahora sql.NullTime

		err := rows.Scan(&data.ID, &data.Temperatura, &data.Umidade, &datahora, &data.Estacao_FK)
		if err != nil {
			log.Fatal(err)
		}
		if datahora.Valid {
			data.DataHora = &datahora.Time
		}

		all_data = append(all_data, data)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return all_data, nil
}
