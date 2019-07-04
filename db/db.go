package db

import (
	"log"

	// mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type dBClient struct {
	db *sqlx.DB
}

// DBCtl 数据库操作客户端
var DBCtl *dBClient

// InitDB init tidb
func InitDB() {
	DBCtl = &dBClient{}
	var err error
	DBCtl.db, err = sqlx.Connect("mysql", "root:@(172.21.0.6:4000)/wcore")
	if err != nil {
		log.Fatalln(err)
	}
}

// GetDB xs
func (me *dBClient) GetDB() *sqlx.DB {
	return me.db
}
