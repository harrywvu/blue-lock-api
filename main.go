package main

import (
	"database/sql"
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


var players = []player {
	{
		ID: 1, 
		Name: "Isagi Yoichi", 
		Age: 17, 
		Height: 175, 
		NELTeam: "Bastard Munchen", 
		PrimaryPosition: "Center Forward", 
		CurrentBlueLockRank: 1,
		Stats: Stats{Overall: 94, Offense: 97.5, Shooting: 90, Speed: 83, Defense: 80, Passing: 83, Dribbling: 77.5},
	},
	{
		ID: 2, 
		Name: "Itoshi Rin", 
		Age: 18, 
		Height: 186, 
		NELTeam: "Paris X Gen", 
		PrimaryPosition: "Left Winger", 
		CurrentBlueLockRank: 2,
		Stats: Stats{Overall: 95, Offense: 96, Shooting: 94, Speed: 88, Defense: 82, Passing: 85, Dribbling: 90},
	},
	{
		ID: 3, 
		Name: "Nagi Seishiro", 
		Age: 17, 
		Height: 190, 
		NELTeam: "Manshine City", 
		PrimaryPosition: "Center Forward", 
		CurrentBlueLockRank: 3,
		Stats: Stats{Overall: 93, Offense: 98, Shooting: 92, Speed: 80, Defense: 75, Passing: 88, Dribbling: 95},
	},
	{
		ID: 4, 
		Name: "Michael Kaiser", 
		Age: 20, 
		Height: 186, 
		NELTeam: "Bastard Munchen", 
		PrimaryPosition: "False Nine", 
		CurrentBlueLockRank: 4,
		Stats: Stats{Overall: 96, Offense: 99, Shooting: 95, Speed: 85, Defense: 78, Passing: 90, Dribbling: 97},
	},
	{
		ID: 5, 
		Name: "Julian Loki", 
		Age: 19, 
		Height: 181, 
		NELTeam: "Paris X Gen", 
		PrimaryPosition: "Left Back", 
		CurrentBlueLockRank: 5,
		Stats: Stats{Overall: 97, Offense: 92, Shooting: 85, Speed: 99, Defense: 98, Passing: 94, Dribbling: 88},
	},
}


func (a *application) getPlayers(c *gin.Context){

	name := c.Query("name")

	if name != "" {
		for _, p := range players {
			if p.Name == name {
				c.IndentedJSON(http.StatusOK, p)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message" : "player not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, players)
}

func (a *application) getPlayerById (c *gin.Context){
	idStr := c.Param("id")

	// if the conversion fails, it returns an error.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message" : "invalid id format"})
		return 
	}

	for _, p := range players {
		if p.ID == id {
			c.IndentedJSON(http.StatusOK, p)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "player not found"})
}