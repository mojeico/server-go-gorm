package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Driver struct {
	gorm.Model

	Type      string `json:"type"`
	FirstName string `json:"first_name" binding:"required"`
	CompanyId int    `json:"company_id"`
	LastName  string `json:"last_name" binding:"required"`
	Address   string `json:"address"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	BirthDay  int64  `json:"birth_day"`

	Zip string `json:"zip"`

	OrderList   pq.Int32Array `json:"order_list" gorm:"type:int[]"`
	ChargesList pq.Int32Array `json:"charges_list" gorm:"type:int[]"`
	FilesList   pq.Int32Array `json:"files_list" gorm:"type:int[]"`

	PayTo      string  `json:"pay_to"` // pay to  company
	PayMethod  string  `json:"pay_method"`
	PayToCount float64 `json:"pay_to_count"`

	PayToOwner  string `json:"pay_to_owner"` // pay to driver
	FuelCards   string `json:"fuel_cards"`
	Transponder string `json:"transponder"`

	Status    string `json:"status"`
	IsActive  bool   `json:"is_active"`
	IsDeleted bool   `json:"is_deleted"`

	Active   bool `json:"active"`
	GrossPay bool `json:"gross_pay"`
	Eligible bool `json:"eligible"`
	Checks   bool `json:"checks"`

	SSN               string `json:"ssn"`
	HireDate          int64  `json:"hire_date"`
	ReviewDate        int64  `json:"review_date"`
	TerminationDate   int64  `json:"termination_date"`
	CdlDate           int64  `json:"cdl_date"`
	LicenseNumber     string `json:"license_number"`
	LicenseExpDate    int64  `json:"license_exp_date"`
	CdlClassification string `json:"cdl_classification"`
	CdlEndorsements   string `json:"cdl_endorsements"`
	TaxForm           string `json:"tax_form"`
	Restrictions      string `json:"restrictions"`

	Emergency        string `json:"emergency"`
	EmergencyPhone   string `json:"emergency_phone"`
	MedialCard       int64  `json:"medial_card"`
	MVRDate          int64  `json:"mvr_date"`
	LastDrugTestDate int64  `json:"last_drug_test_date"`
}

type DriverUpdateInput struct {
	Type      string `json:"type"`
	FirstName string `json:"first_name"`
	CompanyId int    `json:"company_id"`
	LastName  string `json:"last_name"`
	Address   string `json:"address"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	BirthDay  int64  `json:"birth_day"`

	Zip string `json:"zip"`

	PayTo      string  `json:"pay_to"` // pay to  company
	PayMethod  string  `json:"pay_method"`
	PayToCount float64 `json:"pay_to_count"`

	PayToOwner  string `json:"pay_to_owner"` // pay to driver
	FuelCards   string `json:"fuel_cards"`
	Transponder string `json:"transponder"`

	Status    string `json:"status"`
	IsActive  bool   `json:"is_active"`
	IsDeleted bool   `json:"is_deleted"`

	Active            bool `json:"active"`
	GrossPay          bool `json:"gross_pay"`
	EligibleForRehire bool `json:"eligible"`
	Checks            bool `json:"checks"`

	SSN               string `json:"ssn"`
	HireDate          int64  `json:"hire_date"`
	ReviewDate        int64  `json:"review_date"`
	TerminationDate   int64  `json:"termination_date"`
	CDLDate           int64  `json:"cdl_date"`
	LicenseNumber     string `json:"license_number"`
	LicenseExpDate    int64  `json:"license_exp_date"`
	CDLClassification string `json:"cdl_classification"`
	CDLEndorsements   string `json:"cdl_endorsements"`
	TaxForm           string `json:"tax_form"`
	Restrictions      string `json:"restrictions"`

	Emergency        string `json:"emergency"`
	EmergencyPhone   string `json:"emergency_phone"`
	MedialCard       int64  `json:"medial_card"`
	MVRDate          int64  `json:"mvr_date"`
	LastDrugTestDate int64  `json:"last_drug_test_date"`
}
