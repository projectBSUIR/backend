package databases

import (
	"fiber-apis/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DataBase *gorm.DB

var (
	dbName = "backend"
	dbUser = "root"
	dbPass = "password"
	dbHost = "localhost"
)

func Connect(app *fiber.App) {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, "3306", dbName)
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	Migrate(db)
	DataBase = db
}

func Migrate(db *gorm.DB) {
	db.Migrator().CreateTable(&models.User{})
}
