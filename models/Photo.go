package models

import (
	"html"
	"strings"
	"errors"
	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model
	ID       	string	`gorm:"primary_key; not null; unique" json:"id"`
	Title    	string	`gorm:"size:255;not null" json:"title"`
	Caption  	string	`gorm:"size:255;not null" json:"caption"`
	PhotoUrl 	string	`gorm:"not null;" json:"photo_url"`
	UserID   	string	`gorm:"not null; unique" json:"user_id"`
}

func (p *Photo) Init() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Caption = html.EscapeString(strings.TrimSpace(p.Caption))
	p.PhotoUrl = html.EscapeString(strings.TrimSpace(p.PhotoUrl))
}

//Function to validate Photo data
func (p *Photo) Validate(action string) error {
	switch strings.ToLower(action) {
		case "upload": 
			if p.Title == "" {
				return errors.New("title is required")
			} else if p.Caption == "" {
				return errors.New("caption is required")
			} else if p.UserID == "" {
				return errors.New("UserID is required")
			}
			return nil

		case "update":
			if p.Title == "" {
				return errors.New("title is required")
			} else if p.Caption == "" {
				return errors.New("caption is required")
			} else if p.PhotoUrl == "" {
				return errors.New("PhotoUrl is required")
			}
			return nil

		default:
			return nil
	}
}