package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	username string `json:"username"`
	password string `json:"password"`
}
