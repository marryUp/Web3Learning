package model

import "gorm.io/gorm"

// id 、 username 、 password 、 email
type User struct {
	gorm.Model
	Id       uint
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}
