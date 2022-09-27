package models

import (
	"github.com/jinzhu/gorm"
)

type Settlement struct {
	gorm.Model

	SettlementDate   int64   `json:"settlement_date"`
	InvoicingCompany string  `json:"invoicing_company"`
	DriverId         uint    `json:"driver_id" binding:"required"`
	DriverName       string  `json:"driver_name"`
	TotalMiles       int64   `json:"total_miles"`
	EmptyMiles       int64   `json:"empty_miles"`
	LoadedMiles      int64   `json:"loaded_miles"`
	Deductions       float32 `json:"deduction"`
	Reimbursement    float64 `json:"reimbursement"`
	Total            float64 `json:"total"`
	OrderId          uint    `json:"order_id" binding:"required"`

	Status    string `json:"status"`
	IsDeleted bool   `json:"is_deleted"`
	IsActive  bool   `json:"is_active"`
}

type SettlementUpdateInput struct {
	SettlementDate   int64  `json:"settlement_date"`
	InvoicingCompany string `json:"invoicing_company"`
	DriverId         uint   `json:"driver_id" binding:"required"`
	DriverName       string `json:"driver_name"`
	//SettlementJoin   []SettlementOrderAndChargesJoin `gorm:"foreignKey:ID"`
	TotalMiles    int64   `json:"total_miles"`
	OrderId       uint    `json:"order_id" binding:"required"`
	EmptyMiles    int64   `json:"empty_miles"`
	LoadedMiles   int64   `json:"loaded_miles"`
	Deductions    float32 `json:"deduction"`
	Reimbursement float64 `json:"reimbursement"`
	Total         float64 `json:"total"`
}

//type SettlementOrderAndChargesJoin struct {
//	ID            uint    `json:"primary_key"`
//	TotalMiles    int64   `json:"total_miles"`
//	EmptyMiles    int64   `json:"empty_miles"`
//	LoadedMiles   int64   `json:"loaded_miles"`
//	Deductions    float32 `json:"deduction"`
//	Reimbursement float64 `json:"reimbursement"`
//	Total         float64 `json:"total"`
//}

//type SettlementOrderAndChargesJoinUpdateInput struct {
//	TotalMiles    int64   `json:"total_miles"`
//	EmptyMiles    int64   `json:"empty_miles"`
//	LoadedMiles   int64   `json:"loaded_miles"`
//	Deductions    float32 `json:"deduction"`
//	Reimbursement float64 `json:"reimbursement"`
//	Total         float64 `json:"total"`
//}
//
//type Holder struct {
//	DriverId  string        `json:"driverId"`
//	OrderList pq.Int32Array `json:"order_list"`
//}
