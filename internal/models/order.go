package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Order struct {
	gorm.Model
	CompanyId int `json:"company_id"`

	Shipper               string `json:"shipper"` //  company sender
	ShipperFromLocation   string `json:"shipper_from_location"`
	ShipperPickUpTimeFrom int64  `json:"shipper_pick_up_time_from"`
	ShipperPickUpDataFrom int64  `json:"shipper_pick_up_data_from"`
	ShipperPickUpTimeTo   int64  `json:"shipper_pick_up_time_to"`
	ShipperPickUpDataTo   int64  `json:"shipper_pick_up_data_to"`
	ShipperPhone          string `json:"shipper_phone"` // phone company sender

	Consignee                 string `json:"consignee"` //  company receiver
	ConsigneeToLocation       string `json:"consignee_to_location"`
	ConsigneeDeliveryTimeFrom int64  `json:"consignee_delivery_time_from"`
	ConsigneeDeliveryDateFrom int64  `json:"consignee_delivery_date_from"`
	ConsigneeDeliveryTimeTo   int64  `json:"consignee_delivery_time_to"`
	ConsigneeDeliveryDateTo   int64  `json:"consignee_delivery_date_to"`
	ConsigneePhone            string `json:"consignee_phone"`

	ChargesList pq.Int32Array `json:"charges_list" gorm:"type:int[]"`
	FilesList   pq.Int32Array `json:"files_list" gorm:"type:int[]"`

	LoadNumber       string `json:"load_number" binding:"required"`
	PickupNumber     string `json:"pickup_number"`
	DeliveryNumber   string `json:"delivery_number"`
	SealNumber       string `json:"seal_number"`
	BOLNumber        string `json:"bol_number"`
	Commodity        string `json:"commodity"`
	Weight           int    `json:"weight"`
	EquipmentType    string `json:"equipment_type"`
	Temperature      int    `json:"temperature"`
	Pieces           string `json:"pieces"`
	Pallets          int    `json:"pallets"`
	ETAOrFreeDat     int64  `json:"eta_or_free_dat"`
	LTL              bool   `json:"ltl"`
	TotalDays        int    `json:"total_days"`
	BrokerLoadNumber string `json:"broker_load_number"`

	InvoicingCompany string `json:"invoicing_company" binding:"required"`

	BillingMethod string `json:"billing_method" binding:"required"`
	BillingType   string `json:"billing_type" binding:"required"`
	InvoiceID     int    `json:"invoice_id"`
	DriverName    string `json:"driver_name"`
	DriverId      int    `json:"driver_id" required:"required"`
	TruckID       int    `json:"truck_id"`
	TrailerID     int    `json:"trailer_id"`

	Rate          float64 `json:"rate" binding:"required" gorm:"type:decimal(20,2);"`
	GrossPay      float32 `json:"gross_pay"`
	Total         float64 `json:"total"`
	EmptyMiles    float32 `json:"empty_miles"`
	LoadedMiles   float32 `json:"loaded_miles"`
	TotalMiles    float32 `json:"total_miles"`
	ExternalNotes string  `json:"external_notes" type:"text"`

	Status    string `json:"status"`
	IsDeleted bool   `json:"is_deleted"`
	IsActive  bool   `json:"is_active"`
}

type OrderUpdateInput struct {
	gorm.Model
	CompanyId int `json:"company_id"`

	Shipper               string `json:"shipper"` //  company sender
	ShipperFromLocation   string `json:"shipper_from_location"`
	ShipperPickUpTimeFrom int64  `json:"shipper_pick_up_time_from"`
	ShipperPickUpDataFrom int64  `json:"shipper_pick_up_data_from"`
	ShipperPickUpTimeTo   int64  `json:"shipper_pick_up_time_to"`
	ShipperPickUpDataTo   int64  `json:"shipper_pick_up_data_to"`
	ShipperPhone          string `json:"shipper_phone"` // phone company sender

	Consignee                 string `json:"consignee"` //  company receiver
	ConsigneeToLocation       string `json:"consignee_to_location"`
	ConsigneeDeliveryTimeFrom int64  `json:"consignee_delivery_time_from"`
	ConsigneeDeliveryDateFrom int64  `json:"consignee_delivery_date_from"`
	ConsigneeDeliveryTimeTo   int64  `json:"consignee_delivery_time_to"`
	ConsigneeDeliveryDateTo   int64  `json:"consignee_delivery_date_to"`
	ConsigneePhone            string `json:"consignee_phone"`

	ChargesList pq.Int32Array `json:"charges_list" gorm:"type:int[]"`
	FilesList   pq.Int32Array `json:"files_list" gorm:"type:int[]"`

	LoadNumber       string `json:"load_number" binding:"required"`
	PickupNumber     string `json:"pickup_number"`
	DeliveryNumber   string `json:"delivery_number"`
	SealNumber       string `json:"seal_number"`
	BOLNumber        string `json:"bol_number"`
	Commodity        string `json:"commodity"`
	Weight           int    `json:"weight"`
	EquipmentType    string `json:"equipment_type"`
	Temperature      int    `json:"temperature"`
	Pieces           string `json:"pieces"`
	Pallets          int    `json:"pallets"`
	ETAOrFreeDat     int64  `json:"eta_or_free_dat"`
	LTL              bool   `json:"ltl"`
	TotalDays        int    `json:"total_days"`
	BrokerLoadNumber string `json:"broker_load_number"`

	InvoicingCompany string `json:"invoicing_company" binding:"required"`

	BillingMethod string `json:"billing_method" binding:"required"`
	BillingType   string `json:"billing_type" binding:"required"`
	InvoiceID     int    `json:"invoice_id"`
	DriverName    string `json:"driver_name"`
	DriverId      int    `json:"driver_id" required:"required"`
	TruckID       int    `json:"truck_id"`
	TrailerID     int    `json:"trailer_id"`

	Rate          float64 `json:"rate" binding:"required" gorm:"type:decimal(20,2);"`
	GrossPay      float32 `json:"gross_pay"`
	Total         float64 `json:"total"`
	EmptyMiles    float32 `json:"empty_miles"`
	LoadedMiles   float32 `json:"loaded_miles"`
	TotalMiles    float32 `json:"total_miles"`
	ExternalNotes string  `json:"external_notes" type:"text"`

	Status    string `json:"status"`
	IsDeleted bool   `json:"is_deleted"`
	IsActive  bool   `json:"is_active"`
}
