package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Trailer struct {
	gorm.Model

	CompanyId   int            `json:"company_id"`
	DriverID    int            `json:"driver_id" gorm:"type:integer"`
	Name        string         `json:"name" binding:"required" gorm:"type:text"`
	TrailerType pq.StringArray `json:"trailer_type" binding:"required" gorm:"type:text[]"`
	UnitNumber  string         `json:"unit_number" binding:"required" gorm:"type:text"`
	Make        string         `json:"make"`
	Year        int            `json:"year"`
	Plate       string         `json:"plate" binding:"required" gorm:"type:text"`
	State       string         `json:"state" binding:"required" gorm:"type:text"`
	Expiration  int64          `json:"expiration"`
	VinNumber   string         `json:"vin_number" binding:"required" gorm:"type:text"`
	Weight      int            `json:"weight" gorm:"type:integer"`

	OwnerName    string `json:"owner_name" binding:"required" gorm:"type:text"` // pay to  company -> companyOwner
	PurchaseDate int64  `json:"purchase_date"`

	Location string `json:"location" gorm:"type:text"` // create crud for location ---

	FilesList pq.Int64Array `json:"files_list" gorm:"type:int[]"`

	CargoLength int32 `json:"cargo_length"`
	CargoWidth  int32 `json:"cargo_width"`
	CargoHeight int32 `json:"cargo_height"`

	AirRide       bool `json:"air_ride"`
	LoadBars      bool `json:"load_bars"`
	PlatedTrailer bool `json:"plated_trailer"`

	DoorType   pq.StringArray `json:"door_type" gorm:"type:text[]"`
	DoorWidth  int            `json:"door_width"`
	DoorHeight int            `json:"door_height"`

	Straps   bool `json:"straps"`
	LiftGate bool `json:"lift_gate"`
	DockHigh bool `json:"dock_high"`

	DoorWeightCapacity int  `json:"door_weight_capacity"`
	PalletJacks        bool `json:"pallet_jacks"`

	Status    string `json:"status"`
	IsActive  bool   `json:"is_active"`
	IsDeleted bool   `json:"is_deleted"`
}

type TrailerInputUpdate struct {
	CompanyId   int            `json:"company_id"`
	DriverID    int            `json:"driver_id" gorm:"type:integer"`
	Name        string         `json:"name" gorm:"type:text"`
	TrailerType pq.StringArray `json:"trailer_type" gorm:"type:text[]"`
	UnitNumber  string         `json:"unit_number" gorm:"type:text"`
	Make        string         `json:"make"`
	Year        int            `json:"year"`
	Plate       string         `json:"plate" gorm:"type:text"`
	State       string         `json:"state" gorm:"type:text"`
	Expiration  int64          `json:"expiration" `
	VinNumber   string         `json:"vin_number" gorm:"type:text"`
	Weight      int            `json:"weight" gorm:"type:integer"`

	OwnerName    string `json:"owner_name"  gorm:"type:text"` // pay to  company -> companyOwner
	PurchaseDate int64  `json:"purchase_date"`

	Location string `json:"location" gorm:"type:text"` // create crud for location ---

	CargoLength int32 `json:"cargo_length"`
	CargoWidth  int32 `json:"cargo_width"`
	CargoHeight int32 `json:"cargo_height"`

	AirRide       bool `json:"air_ride"`
	LoadBars      bool `json:"load_bars"`
	PlatedTrailer bool `json:"plated_trailer"`

	DoorType   pq.StringArray `json:"door_type"`
	DoorWidth  int            `json:"door_width"`
	DoorHeight int            `json:"door_height"`

	Straps   bool `json:"straps"`
	LiftGate bool `json:"lift_gate"`
	DockHigh bool `json:"dock_high"`

	DoorWeightCapacity int  `json:"door_weight_capacity"`
	PalletJacks        bool `json:"pallet_jacks"`
}
