package main

type stats struct {
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
	Stats               stats `json:"stats"`
}
