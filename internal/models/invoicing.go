package models

import "github.com/jinzhu/gorm"

type Invoicing struct { // Summary about order when order is  invoiced & completed
	gorm.Model

	OrderId      string  `json:"order_id"`
	Date         int64   `json:"date"`
	Pickup       string  `json:"pickup"`        // Order -> ShipperFromLocation
	PickUpDate   int64   `json:"pick_up_date"`  // Will be updated by
	Delivery     string  `json:"delivery"`      // address there
	DeliveryDate int64   `json:"delivery_date"` // Delivery time  - update
	LoadNumber   string  `json:"load_number"`   // Will be updated by
	CustomerName string  `json:"company_name"`  // Shipper from order
	Rate         float64 `json:"rate"`          // Will be updated by
	ExtraPay     float64 `json:"extra_pay"`     // only + // Will be updated by
	Total        float64 `json:"total"`         // rate + ExtraPay  // Will be updated by

	Status    string `json:"status"`
	IsDeleted bool   `json:"is_deleted"`
}

type InvoicingUpdateInput struct {
	gorm.Model

	PickUpDate   int64   `json:"pick_up_date"`
	DeliveryDate int64   `json:"delivery_date"`
	LoadNumber   string  `json:"load_number"`
	Rate         float64 `json:"rate"`
	ExtraPay     float64 `json:"extra_pay"`
}
