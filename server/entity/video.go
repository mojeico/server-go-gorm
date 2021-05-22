package entity

import "time"

type Person struct {
	ID        uint64 `gorm:"primary_key;auto_increment" json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
}
type Video struct {
	ID          uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Author      Person `json:"author" gorm:"foreignkey:PersonID"`
	PersonID    uint64 `json:"-"`

	CreatedAd time.Time `json:"-" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"-" gorm:"default:CURRENT_TIMESTAMP"`
}
