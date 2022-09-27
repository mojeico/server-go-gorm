package models

import "github.com/jinzhu/gorm"

type ExtraPay struct {
	gorm.Model

	OrderID int `json:"order_id" binding:"required"`

	Description      string `json:"description"`
	Type             string `json:"type"`
	InvoicingCompany string `json:"invoicing_company"`
	Distribute       string `json:"distribute"`
	Date             int64  `json:"date"`

	Status    string `json:"status"`
	IsActive  bool   `json:"is_active"`
	IsDeleted bool   `json:"is_deleted"`
}
