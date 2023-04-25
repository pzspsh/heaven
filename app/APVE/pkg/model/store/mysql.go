package store

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	User string
	Pass string
	Ip   string
	Port string
	DB   string
}

func (my *Mysql) Mysql() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", my.User, my.Pass, my.Ip, my.Port, my.DB))
	if err != nil {
	}
	return db
}
