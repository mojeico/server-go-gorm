package reports

import (
	"fmt"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/trucktrace/internal/models"
	"time"
)

func CreateSettlementReport(settlement models.Settlement, order models.Order, chargesList []models.Charges) pdf.Maroto {

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(5, 15, 5)
	m.SetBorder(true)
	m.RegisterHeader(func() {
		m.Row(25, func() {

			m.Col(4, func() {
				m.Text("AnyCompany Name Inc. 851 Any Street Name, Suite 120, Any City, CA 45123.", props.Text{
					Size:  12,
					Align: consts.Left,
				})
			})
			m.Col(8, func() {
				m.Text("Logo", props.Text{
					Size:  12,
					Align: consts.Right,
				})
			})
		})
		m.Row(20, func() {
			m.Col(6, func() {
				m.Text(fmt.Sprintf("Settlement #%d:", settlement.ID), props.Text{
					Size:  12,
					Align: consts.Left,
				})
			})
			m.Col(6, func() {
				m.Text(fmt.Sprintf("Settlement Data : %s", time.Now()), props.Text{
					Size:  12,
					Align: consts.Right,
				})
			})
		})

	})

	m.Row(20, func() {
		m.Col(12, func() {
			m.Text(settlement.DriverName, props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Left,
				Size:  14,
			})
		})
	})

	header := []string{"Load#", "Origin", "Destination", "Distance", "Ship Date", "Delivery Date", "Gross Pay", "Pay Method", "Net Pay"}
	contents := [][]string{{order.LoadNumber}, {order.ShipperFromLocation}, {order.ConsigneeToLocation}, {fmt.Sprint(settlement.TotalMiles)}, {fmt.Sprint(order.ConsigneeDeliveryDateFrom)}, {fmt.Sprint(order.ConsigneeDeliveryDateTo)},
		{fmt.Sprintf("%f", order.GrossPay)}, {order.BillingMethod}, {fmt.Sprintf("%f", order.Total)}}

	m.TableList(header, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{1, 2, 2, 2, 1, 1, 1, 1, 1},
		},
		ContentProp: props.TableListContent{
			Size:      10,
			GridSizes: []uint{1, 2, 2, 2, 1, 1, 1, 1, 1},
		},
		Align: consts.Center,

		HeaderContentSpace: 1,
		Line:               false,
	},
	)

	m.Row(15, func() {

		m.Col(4, func() {
			m.Text("Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Left,
			})
		})
		m.Col(8, func() {
			m.Text("R$ 2.567,00", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
			})
		})
	})
	m.Row(15, func() {
		m.Col(12, func() {
			m.Text("Earnings", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  12,
				Align: consts.Center,
			})
		})
	})
	earningHeader := []string{"Description", "Date", "Rate"}
	var earningContent [][]string
	for _, charge := range chargesList {
		if charge.TypeDeductions == "Deduction" {
			chargeDescription := GetValuesAsArray(charge.Description)
			chargeDate := GetValuesAsArray(charge.ChargeDate)
			chargeRate := GetValuesAsArray(charge.Rate)
			earningContent = append(earningContent, chargeDescription)
			earningContent = append(earningContent, chargeDate)
			earningContent = append(earningContent, chargeRate)
		}
	}
	m.TableList(earningHeader, earningContent, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{4, 4, 4},
		},
		ContentProp: props.TableListContent{
			Size:      10,
			GridSizes: []uint{4, 4, 4},
		},
		Align: consts.Center,

		HeaderContentSpace: 1,
		Line:               false,
	},
	)
	m.Row(15, func() {

		m.Col(4, func() {
			m.Text("Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Left,
			})
		})
		m.Col(8, func() {
			m.Text(fmt.Sprintf("%f", settlement.Reimbursement), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
			})
		})
	})
	m.Row(15, func() {
		m.Col(12, func() {
			m.Text("Deduction", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  12,
				Align: consts.Center,
			})
		})
	})
	deductionHeader := []string{"Description", "Date", "Rate"}

	deductionContent := [][]string{}
	for _, charge := range chargesList {
		if charge.TypeDeductions == "Deduction" {
			chargeDescription := GetValuesAsArray(charge.Description)
			chargeDate := GetValuesAsArray(charge.ChargeDate)
			chargeRate := GetValuesAsArray(charge.Rate)
			deductionContent = append(deductionContent, chargeDescription)
			deductionContent = append(deductionContent, chargeDate)
			deductionContent = append(deductionContent, chargeRate)
		}
	}
	m.TableList(deductionHeader, deductionContent, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{4, 4, 4},
		},
		ContentProp: props.TableListContent{
			Size:      10,
			GridSizes: []uint{4, 4, 4},
		},
		Align: consts.Center,

		HeaderContentSpace: 1,
		Line:               false,
	},
	)
	m.Row(15, func() {

		m.Col(4, func() {
			m.Text("Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Left,
			})
		})
		m.Col(8, func() {
			m.Text(fmt.Sprintf("%f", settlement.Deductions), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
			})
		})
	})
	m.Row(15, func() {

		m.Col(4, func() {
			m.Text("Grand Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Left,
			})
		})
		m.Col(8, func() {
			m.Text(fmt.Sprintf("%f", settlement.Total), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
			})
		})
	})

	return m
}
