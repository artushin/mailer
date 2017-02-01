package pdf

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io"
	"strings"
)

type Envelope struct {
	To   Address
	From Address
}

type Address []string

type Letter struct {
	From      Address
	Person    string
	Signature string
}

func GenerateEnvelope(envelope *Envelope, w io.Writer) error {
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "in",
		Size:    gofpdf.SizeType{Wd: 9.5, Ht: 4.125},
	})
	pdf.SetMargins(0.25, 0.15, 0.25)
	pdf.AddPage()
	pdf.SetFont("times", "", 10)
	pdf.MultiCell(0, .16, strings.Join(envelope.From, "\n"), "", "L", false)
	pdf.CellFormat(0, 1, "", "", 1, "L", false, 0, "")
	pdf.SetLeftMargin(3.5)
	pdf.MultiCell(0, .16, strings.Join(envelope.To, "\n"), "", "L", false)
	return pdf.Output(w)
}

func GenerateLetter(letter *Letter, w io.Writer) error {
	pdf := gofpdf.New("P", "pt", "A4", "")
	pdf.AddPage()
	pdf.SetFont("times", "", 10)
	pdf.CellFormat(0, 40, "", "", 1, "L", false, 0, "")
	pdf.MultiCell(0, 10.5, strings.Join(letter.From, "\n"), "", "L", false)
	pdf.CellFormat(0, 30, "", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 10.5, fmt.Sprintf("Dear %s,", letter.Person), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 500, "", "", 1, "L", false, 0, "")
	pdf.SetCellMargin(140)
	pdf.CellFormat(0, 10.5, "Regards,", "", 1, "R", false, 0, "")
	pdf.CellFormat(0, 45, "", "", 1, "R", false, 0, "")
	pdf.CellFormat(0, 10.5, letter.Signature, "", 1, "R", false, 0, "")
	pdf.SetCellMargin(0)
	return pdf.Output(w)
}
