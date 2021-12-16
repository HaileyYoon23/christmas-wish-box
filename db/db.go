package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

var DB *sql.DB

func init() {
	var err error
	//dsn := "mysql://root@localhost/XMAS"
	//
	//u, err := url.Parse(dsn)
	////if err == nil {
		DB, err = sql.Open("mysql", "root:315931@tcp(127.0.0.1:3306)/XMAS")//u.User.String()+"@tcp("+u.Host+")"+u.RequestURI())
	//}

	if err != nil {
		panic(err)
	}

	maxOpen, _ := strconv.Atoi(os.Getenv("MAX_OPEN_CONN"))
	if maxOpen <= 0 {
		maxOpen = runtime.NumCPU() * 2
	}

	DB.SetConnMaxIdleTime(10 * time.Minute)
	DB.SetMaxIdleConns(maxOpen)
	DB.SetMaxOpenConns(maxOpen)

	createTableQuery := `
		create table IF NOT EXISTS presents ( 
		id BIGINT NOT NULL AUTO_INCREMENT,
		present VARCHAR(100),
		PRIMARY KEY (id),
		UNIQUE (present)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		panic(err)
	}
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