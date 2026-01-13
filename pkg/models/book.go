package models

import "time"

type Book struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateBook struct {
	Title  string `json:"title" binding:"required,min=1,max=255"`
	Author string `json:"author" binding:"required,min=1,max=255"`
}

type UpdateBook struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}
