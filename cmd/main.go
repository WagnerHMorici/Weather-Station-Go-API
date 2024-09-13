package main

import (
	"log"

	"github.com/WagnerHMorici/Weather-Station-Go-API/internal"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := internal.DatabaseConn()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	estacoes, err := internal.QueryStations(db)
	if err != nil {
		log.Fatal(err)
	}

	data, err := internal.QueryDataStations(db)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("/stations", func(c *gin.Context) {
		c.JSON(200, estacoes)
	})

	router.GET("/data", func(c *gin.Context) {
		c.JSON(200, data)
	})

	router.Run()

}
