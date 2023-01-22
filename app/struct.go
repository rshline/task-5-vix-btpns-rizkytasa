package app

import (
	"time"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginData struct {
	ID       string	`json:"id"`
	Username string `json:"users.username"`
	Email    string `json:"users.email"`
	Token    string `json:"users.token"`
	Password string `json:"users.password"`
	Title    string `json:"photos.title"`
	Caption  string `json:"photos.caption"`
	PhotoUrl string `json:"photos.photo_url"`
}

type RegisterInput struct {
	ID        string	`json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string 	`json:"users.password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}