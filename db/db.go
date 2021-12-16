package db

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

var DB *sql.DB

var (
	ErrDuplicateGift = errors.New("already exist gift")
)

type XMASPresent struct {
	Present string
	Likes int
}

func init() {
	var err error

	DB, err = sql.Open("mysql", "doadmin:VmKcD8m0tWnKOv3E@tcp(xmas-do-user-10452914-0.b.db.ondigitalocean.com:25060)/xmas")

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
		likes	INT,
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
	var id int64

	err = db.QueryRow("select id from presents where present=?", gift).Scan(&id)
	if id > 0 {
		err = ErrDuplicateGift
		return
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into presents (present, likes) values (?, ?)")
	_, err = stmt.Exec(gift, 0)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	tx.Commit()
	return
}

func GetGift(db *sql.DB) (xmasList []XMASPresent, err error) {
	var gift string
	var likes int

	rows, err := db.Query("select present, likes from presents order by likes desc")

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gift, &likes)
		if err != nil {
			return
		}
		xmasList = append(xmasList, XMASPresent{
			Present: gift,
			Likes: likes,
		})
	}

	return
}

func UpdateLike(present string) (err error) {
	var likes int

	println(present)
	err = DB.QueryRow("select likes from presents where present=?", present).Scan(&likes)
	if err != nil {
		return
	}

	_, err = DB.Exec("update presents set likes=? where present=?", likes + 1, present)
	if err != nil {
		return
	}

	return
}