package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	AuthorID    uint   `gorm:"foreignKey:AuthorID"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	Likes       int    `gorm:"default:0"`
	Responses   []Response
	Images      []string
	Videos      []string
}

type Response struct {
	gorm.Model
	PostID   uint   `gorm:"foreignKey:PostID"`
	AuthorID uint   `gorm:"foreignKey:AuthorID"`
	Content  string `gorm:"not null"`
	Likes    int    `gorm:"default:0"`
	Images   []string
	Videos   []string
}
