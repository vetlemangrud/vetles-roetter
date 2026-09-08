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
	EatenAt time.Time `json:"eaten_at"`
}


func (r CarrotRepository) initDatabase() error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS carrots (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"eaten_at" DATETIME DEFAULT CURRENT_TIMESTAMP
		);`
	_, err := r.DB.Exec(createTableSQL)
	return err
}

func (r CarrotRepository) findCarrots(page int, pageSize int, from time.Time, to time.Time) ([]Carrot, error) {
	limit := pageSize
	offset := page * pageSize
	query := `
		SELECT id, eaten_at 
		FROM carrots 
		WHERE eaten_at >= ?
		AND eaten_at <= ?
		ORDER BY id 
		DESC 
		LIMIT ? 
		OFFSET ?
	`
	rows, err := r.DB.Query(query, from, to, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	carrots := []Carrot{}
	for rows.Next() {
		var c Carrot
		if err := rows.Scan(&c.ID, &c.EatenAt); err != nil {
			return nil, err
		}
		carrots = append(carrots, c)
	}
	return carrots, nil
}

func (r CarrotRepository) countCarrots(from time.Time, to time.Time) (int, error) {
	query := `
		SELECT COUNT(*) FROM carrots
		WHERE eaten_at >= ?
		AND eaten_at <= ?
	`
	row := r.DB.QueryRow(query, from, to)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
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

func (r CarrotRepository) deleteCarrot(id int) error {
	query := `DELETE FROM carrots WHERE id LIKE ?`
	_, err := r.DB.Exec(query, id)
	return err
}
