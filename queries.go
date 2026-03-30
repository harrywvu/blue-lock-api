package main

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

