package models

import "time"

type LoginUser struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
