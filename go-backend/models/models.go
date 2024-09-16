package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique"`
	Password  string
	CreatedAt time.Time
}

type Post struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	User      User
	Content   string
	CreatedAt time.Time
}
