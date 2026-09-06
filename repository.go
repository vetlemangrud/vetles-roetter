package main

import "database/sql"

type CarrotRepository struct {
	DB *sql.DB
}

func createCarrotTable(db *sql.DB) {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS carrots (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"eaten_at" DATETIME DEFAULT CURRENT_TIMESTAMP
		);`
	_, err := db.Exec(createTableSQL);
	if err != nil {
		panic(err)
	}
}

func (r CarrotRepository) initDatabase() {
	createCarrotTable(r.DB)	
}

func (r CarrotRepository) getCarrots(){

}
