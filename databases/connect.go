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
	DataBase = db
	return err
}

/*func Connect(app *fiber.App) {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, "3306", dbName)
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	Migrate(db)
	DataBase = db
}*/
