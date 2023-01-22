package models

import (
	"html"
	"strings"
	"time"
	"errors"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type User struct {
	gorm.Model
	ID        	string		`gorm:"primary_key; not null; unique" json:"id"`
	Username  	string    	`gorm:"size:255; not null" json:"username"`
	Email     	string    	`gorm:"size:255; not null; unique" json:"email"`
	Password	string    	`gorm:"size:255; not null;" json:"password"`
	Photo		Photo     	`gorm:"constraint:OnUpdate:CASCADE, OnDelete:SET NULL;" json:"photo"`
	CreatedAt	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
// USER METHODS

//Initialize user data
func (u *User) Init() (*User, error) {
	id := uuid.New()
	u.ID = id.String()                      
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Password = html.EscapeString(strings.TrimSpace(u.Password))

	return u, nil
}

func (u *User) InitLogin() (*User, error) {                           
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Password = html.EscapeString(strings.TrimSpace(u.Password))

	return u, nil
}


//Validate user data
func (u *User) ValidateInput(action string) error {
	switch strings.ToLower(action) { //Convert to lowercase

		case "login":
			if u.Email == "" {
				return errors.New("email is required")
			}
			if u.Password == "" {
				return errors.New("password is required")
			}
			return nil

		case "userinput":
			if u.ID == "" {
				return errors.New("ID is required")
			} else if u.Email == "" {
				return errors.New("email is required")
			} else if u.Username == "" {
				return errors.New("username is required")
			} else if u.Password == "" {
				return errors.New("password is required")
			} else if len(u.Password) < 6 {
				return errors.New("minimum password length is 6")
			}
			return nil

		default:
			return nil
	}
}