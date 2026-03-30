package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	db *sql.DB
}

type Stats struct {
	Overall   float64 `json:"overall"`
	Offense   float64 `json:"offense"`
	Shooting  float64 `json:"shooting"`
	Speed     float64 `json:"speed"`
	Defense   float64 `json:"defense"`
	Passing   float64 `json:"passing"`
	Dribbling float64 `json:"dribbling"`
}

type player struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Age                 int    `json:"age"`
	Height              int    `json:"height"`
	NELTeam             string `json:"nel_team"`
	PrimaryPosition     string `json:"primary_position"`
	CurrentBlueLockRank int    `json:"current_blue_lock_rank"`
	Stats               Stats `json:"stats"`
}

const basePlayerQuery = `SELECT p.id, p.name, p.age, p.nelTeam, p.primaryPosition, p.currentBlueLockRank,
							s.overall, s.offense, s.shooting, s.speed, s.defense, s.passing, s.dribbling
							FROM players p
							JOIN stats s ON s.playerid = p.id`


type scanner interface {
	Scan(dest ...any) error	
}

func scanPlayer(s scanner) (player, error) {
	var p player

	err := s.Scan(
		&p.ID, &p.Name, &p.Age, &p.NELTeam, &p.PrimaryPosition, &p.CurrentBlueLockRank,
		&p.Stats.Overall, &p.Stats.Offense, &p.Stats.Shooting, &p.Stats.Speed, &p.Stats.Defense, &p.Stats.Passing, &p.Stats.Dribbling,
	)

	return p, err
}

func main() {

	db := openDB()
	defer db.Close()
	app := application{db: db}

	router := gin.Default()
	router.GET("/players", app.getPlayers)
	router.GET("/players/:id", app.getPlayerById)

	router.Run("localhost:8080")
}

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

// HANDLERS
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

func (a *application) getPlayerByID (id int) (player, error) {
	stmt := basePlayerQuery + " WHERE p.id = ?"
	
	p, err := scanPlayer(a.db.QueryRow(stmt, id))

	return p, err
}

func (a *application) getPlayerByName(name string) (player, error) {
	stmt := basePlayerQuery + " WHERE p.name = ?"

	p, err := scanPlayer(a.db.QueryRow(stmt, name))

	return p, err
}

func (a *application) getAllPlayers() ([]player, error){
	var selectAllstmt = basePlayerQuery

	playerRows, err := a.db.Query(selectAllstmt)
		
	if err != nil {return nil, err}
		
	defer playerRows.Close()

	var players []player // store all the rows

	for playerRows.Next(){			
		p, err := scanPlayer(playerRows)
			
		if err != nil {return nil, err}

		players = append(players, p)
	}

	if err := playerRows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}