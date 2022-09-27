package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Company struct { // company owner
	gorm.Model

	LegalName   string `json:"legal_name" binding:"required"`
	Address     string `json:"address"`
	City        string `json:"city"`
	Country     string `json:"country"`
	PhoneNumber string `json:"phone_number"`
	FaxNumber   string `json:"fax_number"`
	McNumber    int    `json:"mc_number"`

	BillingAddress string `json:"billing_address"`
	BillingMethod  string `json:"billing_method"`
	BillingType    string `json:"billing_type"`
	BillingEmail   string `json:"billing_email"`

	Drivers pq.Int64Array `json:"drivers" gorm:"type:int[]"`

	PictureName string `json:"picture_name"`

	Status    string `json:"status"`
	IsDeleted bool   `json:"is_deleted"`
	IsActive  bool   `json:"is_active"`
}
type CompanyUpdateInput struct {
	gorm.Model

	LegalName   string `json:"legal_name" `
	Address     string `json:"address"`
	City        string `json:"city"`
	Country     string `json:"country"`
	PhoneNumber string `json:"phone_number"`
	FaxNumber   string `json:"fax_number"`
	McNumber    int    `json:"mc_number"`

	BillingAddress string `json:"billing_address"`
	BillingMethod  string `json:"billing_method"`
	BillingType    string `json:"billing_type"`
	BillingEmail   string `json:"billing_email"`

	Drivers pq.Int64Array `json:"drivers" gorm:"type:int[]"`

	PictureName string `json:"picture_name"`
	PictureBody []byte ` json:"picture_body" form:"file"`

	Status    string `json:"status"`
	IsDeleted bool   `json:"is_deleted"`
	IsActive  bool   `json:"is_active"`
}
