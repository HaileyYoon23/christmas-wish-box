package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB
var dbName = "tables.db"

func init() {
	var err error

	DB, err = InitDB(dbName)
	if err != nil {
		panic(err)
	}
}

func InitDB(file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	createTableQuery := `
		create table IF NOT EXISTS presents ( 
		id integer PRIMARY KEY autoincrement,
		present string,
		UNIQUE (gift)
		)
	`
	_, e := db.Exec(createTableQuery)
	if e != nil {
		return nil, e
	}

	return db, nil
}

func AddGift(db *sql.DB, gift string) (err error) {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into presents (present) values (?)")
	_, err = stmt.Exec(gift)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	tx.Commit()
	return
}

func GetGift(db *sql.DB) (gifts []string, err error) {
	var gift string

	rows, err := db.Query("select present from presents")

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gift)
		if err != nil {
			return
		}
		gifts = append(gifts, gift)
	}

	return
}