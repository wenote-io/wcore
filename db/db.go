package db

import (
	"log"

	// mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DB Connect
var db *sqlx.DB

// InitDB init tidb
func InitDB() {
	var err error
	db, err = sqlx.Connect("mysql", "root:@(172.21.0.6:4000)/wcore")
	if err != nil {
		log.Fatalln(err)
	}
}

// GetDB xs
func GetDB() *sqlx.DB {
	return db
}
