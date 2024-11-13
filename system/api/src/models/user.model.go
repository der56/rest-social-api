package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID          sql.NullString `json:"id"`
	Username    string         `json:"username"`
	FirstName   string         `json:"firstName"`
	LastName    string         `json:"lastName"`
	Email       string         `json:"email"`
	Password    string         `json:"password"`
	Picture     sql.NullString `json:"picture"`
	Description sql.NullString `json:"description"`

	Posts []Post `gorm:"foreignKey:UserID"`
}

type Followers struct {
	gorm.Model
	Followers []User `json:"followers"`
	Following []User `json:"following"`
}

type UserRelevantInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
