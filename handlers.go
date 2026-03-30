package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

func (a *application) getPlayers(c *gin.Context){

	name := c.Query("name")

	if name == ""{

		players, err := a.getAllPlayers()

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "database error"})
			return
		}

		c.IndentedJSON(http.StatusOK, players)

	} else {
		
		p, err := a.getPlayerByName(name)

		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message" : "player not found"})
			return
		}
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "database error"})
			return
		}

		c.IndentedJSON(http.StatusOK, p)
	}
}

func (a *application) getPlayerById (c *gin.Context){

	// GET THE ID FROM THE PARAM
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message" : "invalid id"})
		return
	}

	p, err := a.getPlayerByID(id)

	if err == sql.ErrNoRows {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message" : "player not found"})
		return
	}

	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "database error"})
		return
	}

	c.IndentedJSON(http.StatusOK, p)
}