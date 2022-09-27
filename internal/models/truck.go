package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Truck struct {
	gorm.Model
	CompanyId       int    `json:"company_id"`
	UnitNumber      string `json:"unit_number"`
	Make            string `json:"make" binding:"required"`
	TruckModel      string `json:"truck_model"`
	Year            int    `json:"year"`
	Plate           uint   `json:"plate" binding:"required"`
	State           string `json:"state" binding:"required"`
	PlateExpiration int64  `json:"plate_expiration"`
	VINNumber       string `json:"vin_number"`

	Transporter string  `json:"transporter"`
	FuelCard    string  `json:"fuel_card"`
	FuelType    string  `json:"fuel_type" binding:"required"`
	Mpg         float64 `json:"mpg"`

	CoordinatorName string `json:"coordinator_name"`
	CompanyName     string `json:"company_name"`
	OwnerName       string `json:"owner_name"`
	DriverName      string `json:"driver_name"`
	AssignDate      int64  `json:"assign_date"`
	CoDriverName    string `json:"co_driver_name"`
	TrailerName     string `json:"trailer_name"`
	Location        string `json:"location" binding:"required"` // create crud for location ---
	LocationDate    int64  `json:"location_date" time_format:"2006-01-02" binding:"required"`

	TaxDueDate   int64 `json:"tax_due_date"`
	PurchaseDate int64 `json:"purchase_date"`

	FilesList pq.Int64Array `json:"files_list" gorm:"type:int[]"`

	LiabilityEffDate int64 `json:"liability_eff_date"`
	LiabilityExpDate int64 `json:"liability_exp_date"`

	Status    string `json:"status"`
	IsActive  bool   `json:"is_active"`
	IsDeleted bool   `json:"is_deleted"`
}

type TruckUpdateInput struct {
	CompanyId       int    `json:"company_id"`
	UnitNumber      string `json:"unit_number"`
	Make            string `json:"make" binding:"required"`
	TruckModel      string `json:"truck_model"`
	Year            int    `json:"year"`
	Plate           uint   `json:"plate" binding:"required"`
	State           string `json:"state" binding:"required"`
	PlateExpiration int64  `json:"plate_expiration"`
	VINNumber       string `json:"vin_number"`

	Transporter string  `json:"transporter"`
	FuelCard    string  `json:"fuel_card"`
	FuelType    string  `json:"fuel_type" binding:"required"`
	Mpg         float32 `json:"mpg"`

	DriverName   string `json:"driver_name"`
	AssignDate   int64  `json:"assign_date"`
	CoDriverName string `json:"co_driver_name"`
	TrailerName  string `json:"trailer_name"`
	Location     string `json:"location" binding:"required"` // create crud for location ---
	LocationDate int64  `json:"location_date" time_format:"2006-01-02" binding:"required"`

	LiabilityEffDate int64 `json:"liability_eff_date"`
	LiabilityExpDate int64 `json:"liability_exp_date"`
}
