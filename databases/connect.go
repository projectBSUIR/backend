package databases

import (
	"database/sql"
	"fmt"
)

var DataBase *sql.DB

var (
	dbName = "backend"
	dbUser = "root"
	dbPass = "password"
	dbHost = "localhost"
)

func ConnectDB() error {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, "3306", dbName)
	db, err := sql.Open("mysql", conn)
	db.SetMaxIdleConns(300)
	DataBase = db
	return err
}
