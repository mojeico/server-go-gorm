package reports

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/jung-kurt/gofpdf"
	"reflect"
	"time"
)

func GetHeaderAsArray(data interface{}) []string {
	v := reflect.ValueOf(data)
	n := v.NumField()

	st := reflect.TypeOf(data)
	headers := make([]string, n)
	for i := 0; i < n; i++ {

		if st.Field(i).Name == "Model" {
			headers[i] = "ID"
		} else {
			headers[i] = st.Field(i).Name
		}
	}

	return headers
}

func GetValuesAsArray(data interface{}) []string {
	v := reflect.ValueOf(data)
	n := v.NumField()

	rowContents := make([]string, n)
	for i := 0; i < n; i++ {
		x := v.Field(i)

		if i == 0 {
			iii := x.Interface().(gorm.Model)
			s := fmt.Sprintf("%v", iii.ID)
			rowContents[i] = s
		} else {
			s := fmt.Sprintf("%v", x.Interface())
			rowContents[i] = s
		}

	}

	return rowContents
}

func NewReport(Wt float64, reportType string) *gofpdf.Fpdf {
	pdf := gofpdf.New("L", "mm", "Letter", "")

	pdf.SetFont("Times", "B", 28)

	s := gofpdf.SizeType{
		Wd: Wt,
		Ht: 400,
	}

	pdf.AddPageFormat("P", s)

	pdf.Cell(40, 10, fmt.Sprintf("%v Report", reportType))

	pdf.Ln(12)

	pdf.SetFont("Times", "", 20)
	pdf.Cell(40, 10, time.Now().Format("Mon Jan 2, 2006"))
	pdf.Ln(20)

	return pdf
}

func Header(pdf *gofpdf.Fpdf, hdr []string, w float64) *gofpdf.Fpdf {
	pdf.SetFont("Times", "B", 16)
	pdf.SetFillColor(240, 240, 240)
	for _, str := range hdr {
		pdf.CellFormat(w, 7, str, "1", 0, "C", true, 0, "")
	}

	pdf.Ln(-1)
	return pdf
}

func Table(pdf *gofpdf.Fpdf, tbl [][]string, w float64) *gofpdf.Fpdf {
	pdf.SetFont("Times", "", 16)
	pdf.SetFillColor(255, 255, 255)

	for _, line := range tbl {
		for _, str := range line {
			pdf.CellFormat(w, 7, str, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}
	return pdf
}

func Image(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.ImageOptions("stats.png", 225, 10, 25, 25, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	return pdf
}

func SavePDF(pdf *gofpdf.Fpdf) error {
	return pdf.OutputFileAndClose("report.pdf")
}
