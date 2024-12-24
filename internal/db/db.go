package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	mysqlCfg "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/advent_calendar")
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
}

func NewMySqlStorage(cfg mysqlCfg.Config) (*sql.DB, error) {
	dsn := cfg.FormatDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
