package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Safety struct {
	gorm.Model

	CompanyID int `json:"company_id"`

	FilesList pq.Int32Array `json:"files_list" gorm:"type:int[]"`

	SafetyType string `json:"safety_type"  binding:"required"` // reports permits registration
	Comments   string `json:"comments"`

	Status   string `json:"status"`
	Deleted  bool   `json:"deleted"`
	IsActive bool   `json:"is_active"`
}

type SafetyInputUpdate struct {
	FilesList pq.Int32Array `json:"files_list" gorm:"type:int[]"`

	Comments string `json:"comments"`
}
