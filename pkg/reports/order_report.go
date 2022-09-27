package reports

import (
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/trucktrace/internal/models"
)

func CreateOrderReport(order models.Order) pdf.Maroto {

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	m.SetAliasNbPages("{nb}")
	m.SetFirstPageNb(1)

	m.RegisterHeader(func() {

		m.Row(18, func() {

			m.Col(6, func() {
				m.Text("TEST", props.Text{
					Top:   1,
					Align: consts.Left,
				})
				m.Text("123 STATE STR  BROADVIEW, IL 60155 USA", props.Text{
					Top:   5,
					Align: consts.Left,
				})
				m.Text("+1 (773) 369-8542", props.Text{
					Top:   9,
					Align: consts.Left,
				})

				m.Text("test@utechcorp.com", props.Text{
					Top:   13,
					Align: consts.Left,
				})
			})

		})

		m.Row(15, func() {
			m.Col(0, func() {
				m.Text("INVOICE", props.Text{
					Size:  15,
					Style: consts.Bold,
					Align: consts.Center,
					Top:   50,

					Color: color.Color{
						Blue: 180,
					},
				})

			})
		})

	})

	m.Row(50, func() {
		m.Col(7, func() {
			m.Text("BILL TO:", props.Text{
				Top:   20,
				Style: consts.Bold,
				Align: consts.Left,
			})
			m.Text("TEST CARRIER", props.Text{
				Top:   24,
				Align: consts.Left,
			})
			m.Text("2100 S 21ST STREET", props.Text{
				Top:   28,
				Align: consts.Left,
			})
			m.Text("BROADVIEW, IL 60155", props.Text{
				Top:   32,
				Align: consts.Left,
			})
			m.Text("USA", props.Text{
				Top:   36,
				Align: consts.Left,
			})
		})

		m.Col(3, func() {
			m.Text("INVOICE #", props.Text{
				Top:   20,
				Style: consts.Bold,
				Align: consts.Left,
			})
			m.Text("DATE:", props.Text{
				Top:   24,
				Align: consts.Left,
			})
			m.Text("DUE DATE:", props.Text{
				Top:   28,
				Align: consts.Left,
			})
			m.Text("LOAD NUMBER:", props.Text{
				Top:   32,
				Align: consts.Left,
			})

		})

		m.Col(2, func() {
			m.Text("814", props.Text{
				Top:   20,
				Style: consts.Bold,
				Align: consts.Left,
			})
			m.Text("12/13/2020", props.Text{
				Top:   24,
				Align: consts.Left,
			})
			m.Text("----", props.Text{
				Top:   28,
				Align: consts.Left,
			})
			m.Text("IP", props.Text{
				Top:   32,
				Align: consts.Left,
			})

		})

	})

	m.Row(10, func() {
		m.Col(7, func() {
			m.Text("Description", props.Text{
				Top:   12,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})

		m.Col(3, func() {
			m.Text("Price Each", props.Text{
				Top:   12,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})

		m.Col(2, func() {
			m.Text("Amount", props.Text{
				Top:   12,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})

	})

	m.Row(15, func() {
		m.Col(7, func() {
			m.Text("SCHAUMBURG, IL 60195, USA - BRIDGEVIEW, IL 60455, USA", props.Text{
				Top:   8,
				Align: consts.Left,
			})
		})

		m.Col(3, func() {
			m.Text("$5,000,000.00", props.Text{
				Top:   8,
				Align: consts.Left,
			})
		})

		m.Col(2, func() {
			m.Text("$5,000,000.00", props.Text{
				Top:   8,
				Align: consts.Left,
			})
		})

	})

	m.Line(1.0)

	m.Row(2, func() {

		m.Col(8, func() {
			m.Text("Total", props.Text{
				Top:   8,
				Align: consts.Right,
			})
		})

		m.Col(5, func() {
			m.Text("$5,000,000.00", props.Text{
				Top:   8,
				Align: consts.Center,
			})
		})

	})

	return m
}
