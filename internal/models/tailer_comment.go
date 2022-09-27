package models

import "github.com/jinzhu/gorm"

type TrailerComment struct {
	gorm.Model
	Description string `json:"description"`
	Comment     string `json:"comment"`
	TrailerId   int    `json:"trailer_id"`
}
