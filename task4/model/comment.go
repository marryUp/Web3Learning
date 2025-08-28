package model

import (
	"gorm.io/gorm"
)

// id 、 content 、 user_id （关联 users 表的 id ）、 post_id （关联 posts 表的 id ）、 created_at
type Comment struct {
	gorm.Model
	Id      uint
	Content string `gorm:"not null"`
	UserId  uint   `gorm:"not null"`
	PostId  uint   `gorm:"not null"`
}
