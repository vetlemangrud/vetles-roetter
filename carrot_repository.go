package main

import (
	"database/sql"
	"time"
)

type CarrotRepository struct {
	DB *sql.DB
}

type Carrot struct {
	ID      int       `json:"id"`
	EatenAt time.Time `json:eaten_at`
}


func (r CarrotRepository) initDatabase() {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS carrots (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"eaten_at" DATETIME DEFAULT CURRENT_TIMESTAMP
		);`
	_, err := r.DB.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}
}

func (r CarrotRepository) findCarrots() ([]Carrot, error) {
	query := `SELECT id, eaten_at FROM carrots`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carrots []Carrot
	for rows.Next() {
		var c Carrot
		if err := rows.Scan(&c.ID, &c.EatenAt); err != nil {
			return nil, err
		}
		carrots = append(carrots, c)
	}
	return carrots, nil
}

func (r CarrotRepository) addCarrot() (Carrot, error) {
	query := `INSERT INTO carrots DEFAULT VALUES RETURNING id, eaten_at`
	row := r.DB.QueryRow(query)
	var carrot Carrot
	if err := row.Scan(&carrot.ID, &carrot.EatenAt); err != nil {
		return Carrot{}, err
	}

	return carrot, nil
}
