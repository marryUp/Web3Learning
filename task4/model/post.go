package model

import (
	"gorm.io/gorm"
)

// id 、 title 、 content 、 user_id （关联 users 表的 id ）、
// created_at 、 updated_at
type Post struct {
	gorm.Model
	Id      uint
	Title   string `gorm:"unique;not null"`
	Content string
	UserId  uint `gorm:"not null"`
}
