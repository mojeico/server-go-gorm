package models

import "github.com/jinzhu/gorm"

type File struct {
	gorm.Model

	Name             string `json:"name"`
	Extension        string `json:"extension"`
	Size             int64  `json:"size"`
	ExpirationDate   int64  `json:"expiration_date"  binding:"required"`
	ExpirationStatus string `json:"expiration_status"` // expiredSoon(30) expired (0)

	Comment string `json:"comment"`

	OwnerType string `json:"owner_type"`
	OwnerId   int    `json:"owner_id"`

	Status    string `json:"status"`
	IsActive  bool   `json:"is_active"`
	IsDeleted bool   `json:"is_deleted"`
}
