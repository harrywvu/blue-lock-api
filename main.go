package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

type application struct { db *sql.DB }


func main() {

	db := openDB()
	defer db.Close()
	app := application{db: db}

	router := gin.Default()
	router.GET("/players", app.getPlayers)
	router.GET("/players/:id", app.getPlayerById)

	router.Run("localhost:8080")
}