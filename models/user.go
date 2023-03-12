package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string `json:"username"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Password  string `json:"password"`
}
