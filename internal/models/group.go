package models

import (
	"github.com/jinzhu/gorm"
	pq "github.com/lib/pq"
)

type Groups struct {
	gorm.Model
	Name       string         `json:"name" binding:"required" gorm:"type:text"`
	CompanyId  string         `json:"company_id" gorm:"type:text"`
	Priveleges pq.StringArray `json:"priveleges" binding:"required" gorm:"type:text[]"`
	Users      pq.StringArray `json:"users" gorm:"type:text[]"`

	IsDeleted bool   `gorm:"type:bool" json:"is_deleted"`
	Status    string `json:"status"`
	IsActive  bool   `json:"is_active"`
}

type GroupUpdateInput struct {
	Name       string         `gorm:"type:text," json:"name"`
	Priveleges pq.StringArray `gorm:"type:text[]" json:"priveleges"`
	Users      pq.StringArray `gorm:"type:integer[]" json:"users"`
}
