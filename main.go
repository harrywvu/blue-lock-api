package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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


func main() {

	db := openDB()
	defer db.Close()
	app := application{db: db}

	router := gin.Default()
	router.GET("/players", app.getPlayers)
	router.GET("/players/:id", app.getPlayerById)

	router.Run("localhost:8080")
}


// var players = []player {
// 	{
// 		ID: 1, 
// 		Name: "Isagi Yoichi", 
// 		Age: 17, 
// 		Height: 175, 
// 		NELTeam: "Bastard Munchen", 
// 		PrimaryPosition: "Center Forward", 
// 		CurrentBlueLockRank: 1,
// 		Stats: Stats{Overall: 94, Offense: 97.5, Shooting: 90, Speed: 83, Defense: 80, Passing: 83, Dribbling: 77.5},
// 	},
// 	{
// 		ID: 2, 
// 		Name: "Itoshi Rin", 
// 		Age: 18, 
// 		Height: 186, 
// 		NELTeam: "Paris X Gen", 
// 		PrimaryPosition: "Left Winger", 
// 		CurrentBlueLockRank: 2,
// 		Stats: Stats{Overall: 95, Offense: 96, Shooting: 94, Speed: 88, Defense: 82, Passing: 85, Dribbling: 90},
// 	},
// 	{
// 		ID: 3, 
// 		Name: "Nagi Seishiro", 
// 		Age: 17, 
// 		Height: 190, 
// 		NELTeam: "Manshine City", 
// 		PrimaryPosition: "Center Forward", 
// 		CurrentBlueLockRank: 3,
// 		Stats: Stats{Overall: 93, Offense: 98, Shooting: 92, Speed: 80, Defense: 75, Passing: 88, Dribbling: 95},
// 	},
// 	{
// 		ID: 4, 
// 		Name: "Michael Kaiser", 
// 		Age: 20, 
// 		Height: 186, 
// 		NELTeam: "Bastard Munchen", 
// 		PrimaryPosition: "False Nine", 
// 		CurrentBlueLockRank: 4,
// 		Stats: Stats{Overall: 96, Offense: 99, Shooting: 95, Speed: 85, Defense: 78, Passing: 90, Dribbling: 97},
// 	},
// 	{
// 		ID: 5, 
// 		Name: "Julian Loki", 
// 		Age: 19, 
// 		Height: 181, 
// 		NELTeam: "Paris X Gen", 
// 		PrimaryPosition: "Left Back", 
// 		CurrentBlueLockRank: 5,
// 		Stats: Stats{Overall: 97, Offense: 92, Shooting: 85, Speed: 99, Defense: 98, Passing: 94, Dribbling: 88},
// 	},
// }


func (a *application) getPlayers(c *gin.Context){

	name := c.Query("name")

	if name == ""{
		var selectAllstmt = `SELECT p.id, p.name, p.age, p.nelTeam, p.primaryPosition, p.currentBlueLockRank,
       			s.overall, s.offense, s.shooting, s.speed, s.defense, s.passing, s.dribbling
				FROM players p
				JOIN stats s ON s.playerid = p.id`

	playerRows, err := a.db.Query(selectAllstmt)
	if err != nil {log.Fatal(err)}
	defer playerRows.Close()

	var players []player
	for playerRows.Next(){
		var p player
		err = playerRows.Scan(
			&p.ID, &p.Name, &p.Age, &p.NELTeam, &p.PrimaryPosition, &p.CurrentBlueLockRank,
			&p.Stats.Overall, &p.Stats.Offense, &p.Stats.Shooting, &p.Stats.Speed, &p.Stats.Defense, &p.Stats.Passing, &p.Stats.Dribbling,)
		if err != nil {log.Fatal(err)}
		players = append(players, p)
	}

	c.IndentedJSON(http.StatusOK, players)

	} else {
		var stmt string = `SELECT p.id, p.name, p.age, p.nelTeam, p.primaryPosition, p.currentBlueLockRank,
							s.overall, s.offense, s.shooting, s.speed, s.defense, s.passing, s.dribbling
							FROM players p
							JOIN stats s ON s.playerid = p.id
						WHERE p.name = ?`
		var p player
		err := a.db.QueryRow(stmt, name).Scan(
			&p.ID, &p.Name, &p.Age, &p.NELTeam, &p.PrimaryPosition, &p.CurrentBlueLockRank,
			&p.Stats.Overall, &p.Stats.Offense, &p.Stats.Shooting, &p.Stats.Speed, &p.Stats.Defense, &p.Stats.Passing, &p.Stats.Dribbling,
		)
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

	var stmt string = `SELECT p.id, p.name, p.age, p.nelTeam, p.primaryPosition, p.currentBlueLockRank,
						s.overall, s.offense, s.shooting, s.speed, s.defense, s.passing, s.dribbling
						FROM players p
						JOIN stats s ON s.playerid = p.id
					WHERE p.id = ?`


	var p player
	err = a.db.QueryRow(stmt, id).Scan(
		&p.ID, &p.Name, &p.Age, &p.NELTeam, &p.PrimaryPosition, &p.CurrentBlueLockRank,
		&p.Stats.Overall, &p.Stats.Offense, &p.Stats.Shooting, &p.Stats.Speed, &p.Stats.Defense, &p.Stats.Passing, &p.Stats.Dribbling,
	)
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