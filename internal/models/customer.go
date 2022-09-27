package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Customer struct { // company receiver and company sender
	gorm.Model

	Address     string `json:"address"`
	City        string `json:"city"`
	Country     string `json:"country"`
	PhoneNumber string `json:"phone_number"`
	FaxNumber   string `json:"fax_number"`
	LegalName   string `json:"legal_name" binding:"required"`
	McNumber    int    `json:"mc_number"`
	DotNumber   string `json:"dot_number"`

	FilesList pq.Int64Array `json:"files_list" gorm:"type:int[]"`

	BillingAddress      string  `json:"billing_address"`
	BillingMethod       string  `json:"billing_method"`
	BillingType         string  `json:"billing_type"`
	BillingEmail        string  `json:"billing_email"`
	BillingCreditLimit  float64 `json:"billing_credit_limit"`
	BillingBalance      float64 `json:"billing_balance"`
	BillingTotalBalance float64 `json:"billing_total_balance"`

	Status    string `json:"status"`
	IsDeleted bool   `json:"is_deleted"`
	IsActive  bool   `json:"is_active"`
}

type CustomerUpdateInput struct {
	Address     string `json:"address"`
	City        string `json:"city"`
	Country     string `json:"country"`
	PhoneNumber string `json:"phone_number"`
	FaxNumber   string `json:"fax_number"`
	LegalName   string `json:"legal_name"`
	McNumber    int    `json:"mc_number"`
	DotNumber   string `json:"dot_number"`

	BillingAddress string `json:"billing_address"`

	BillingMethod       string  `json:"billing_method"`
	BillingType         string  `json:"billing_type"`
	BillingEmail        string  `json:"billing_email"`
	BillingCreditLimit  float64 `json:"billing_credit_limit"`
	BillingBalance      float64 `json:"billing_balance"`
	BillingTotalBalance float64 `json:"billing_total_balance"`
}

func (s *Customer) IsValid() error {

	if len(s.LegalName) == 0 {
		return errors.New("LegalName is empty")
	}

	if s.BillingCreditLimit < 0 ||
		s.BillingBalance < 0 ||
		s.BillingTotalBalance < 0 {

		return errors.New("negative sum of money")

	}

	return nil
}
