package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Charges struct { // deductions(-) and earnings(+)
	gorm.Model
	OrderId        int            `json:"order_id" binding:"required" `
	ChargeDate     int64          `json:"charge_date"  binding:"required"`
	DriverName     string         `json:"driver_name"  binding:"required"`
	DriverId       int            `json:"driver_id"  binding:"required"`
	TypeDeductions string         `json:"type_deductions"  binding:"required"` // deductions(-) and earnings(+)
	CompanyName    string         `json:"company_name"`                        //company driver owner
	Rate           float64        `json:"rate"  binding:"required"`            // +100 or -100
	Description    string         `json:"description"`
	Files          pq.StringArray `json:"files" gorm:"type:text[]"` //  files list (with possibility to open)
	IsDeleted      bool           `json:"is_deleted"`
	IsActive       bool           `json:"is_active"`
	Status         string         `json:"status"`
}

type ChargesUpdateInput struct {
	OrderId        int            `json:"order_id"`
	ChargeDate     int64          `json:"charge_date"`
	DriverName     string         `json:"driver_name"`
	TypeDeductions string         `json:"type_deductions"`
	CompanyName    string         `json:"company_name"` //company driver owner
	Rate           float64        `json:"rate"`         // +100 or -100
	Description    string         `json:"description"`
	Files          pq.StringArray `json:"files" gorm:"type:text[]"` //  files list (with possibility to open)
	IsActive       bool           `json:"is_active"`
	Status         string         `json:"status"`
}
