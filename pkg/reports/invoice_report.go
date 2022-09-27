package reports

import (
	"fmt"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/trucktrace/internal/models"
	"time"
)

func CreateInvoiceReport(invoice models.Invoicing) pdf.Maroto {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(5, 15, 5)

	m.RegisterHeader(func() {
		m.Row(20, func() {

			m.Col(6, func() {
				m.Text("AnyCompany Name Inc. 851 Any Street Name, Suite 120, Any City, CA 45123.", props.Text{
					Size:  12,
					Align: consts.Left,
				})
			})
			m.Col(6, func() {
				m.Text("Logo", props.Text{
					Size:  12,
					Align: consts.Right,
				})
			})
		})

	})

	m.Row(12, func() {
		m.Col(12, func() {
			m.Text("Invoice", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Center,
				Size:  16,
			})
		})
	})

	m.Row(10, func() {
		m.Col(2, func() {
			m.Text("Bill To :", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Left,
				Size:  12,
			})
		})
		m.Col(8, func() {
			m.Text("Invoice :", props.Text{
				Top:   3,
				Align: consts.Right,
				Size:  12,
			})
		})

		m.Col(2, func() {
			m.Text(fmt.Sprint(invoice.ID), props.Text{
				Top:   3,
				Align: consts.Right,
				Size:  12,
			})
		})
	})

	m.Row(10, func() {
		m.Col(4, func() {
			m.Text(invoice.CustomerName, props.Text{
				Align: consts.Left,
				Size:  12,
			})
		})
		m.Col(6, func() {
			m.Text("Date :", props.Text{
				Align: consts.Right,
				Size:  12,
			})
		})
		m.Col(2, func() {
			m.Text(time.Now().String(), props.Text{
				Align: consts.Right,
				Size:  12,
			})
		})
	})
	m.Row(10, func() {
		m.Col(10, func() {
			m.Text("Load Number", props.Text{
				Align: consts.Right,
				Size:  12,
			})
		})
		m.Col(2, func() {
			m.Text(fmt.Sprint(invoice.LoadNumber), props.Text{
				Align: consts.Right,
				Size:  12,
			})
		})
	})
	header := []string{"Description", "Price Each", "Amount"}
	contents := [][]string{
		GetValuesAsArray(invoice.Delivery), GetValuesAsArray(invoice.Rate), GetValuesAsArray(invoice.Total),
	}

	m.Row(20, func() {
		m.Col(12, func() {
			m.Text("white", props.Text{
				Size:  12,
				Align: consts.Center,
				Top:   20,
				Color: color.NewWhite(),
			})
		})
	})
	m.TableList(header, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{6, 3, 3},
		},
		ContentProp: props.TableListContent{
			Size:      10,
			GridSizes: []uint{6, 3, 3},
		},
		Align: consts.Center,

		HeaderContentSpace: 1,
		Line:               false,
	},
	)

	m.Row(40, func() {

		m.Col(9, func() {
			m.Text("Total :", props.Text{
				Top:   30,
				Align: consts.Right,
				Size:  12,
				Style: consts.Bold,
			})
		})
		m.Col(3, func() {
			m.Text(fmt.Sprintf("%f", invoice.Total), props.Text{
				Top:   30,
				Align: consts.Center,
				Size:  12,
				Style: consts.Bold,
			})
		})
	})

	m.Row(20, func() {
		m.Col(12, func() {
			m.Text("NOA:", props.Text{
				Top:   20,
				Align: consts.Center,
				Size:  12,
			})
		})
	})

	return m
}
