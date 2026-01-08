package export

import (
	"fmt"

	"github.com/andrijan/bc-objects-counter/internal/counter"
	"github.com/go-pdf/fpdf"
)

// ToPDF exports the summary to a PDF file.
func ToPDF(summary *counter.Summary, filePath string) error {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("BC Objects Summary", false)
	pdf.SetAuthor("BC Objects Counter", false)

	// Add summary page
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 12, "Business Central Objects Summary")
	pdf.Ln(16)

	// Summary table header
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(68, 114, 196)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(80, 8, "Object Type", "1", 0, "L", true, 0, "")
	pdf.CellFormat(40, 8, "Count", "1", 1, "C", true, 0, "")

	// Summary table data
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0)
	for i, c := range summary.CountsByType {
		fill := i%2 == 0
		if fill {
			pdf.SetFillColor(240, 240, 240)
		}
		pdf.CellFormat(80, 7, c.Type, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(40, 7, fmt.Sprintf("%d", c.Count), "1", 1, "C", fill, 0, "")
	}

	// Total row
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(80, 8, "TOTAL", "1", 0, "L", true, 0, "")
	pdf.CellFormat(40, 8, fmt.Sprintf("%d", summary.TotalObjects), "1", 1, "C", true, 0, "")

	// Details section (new page)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Object Details")
	pdf.Ln(14)

	// Details table header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(68, 114, 196)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(35, 7, "Type", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 7, "ID", "1", 0, "C", true, 0, "")
	pdf.CellFormat(70, 7, "Name", "1", 0, "L", true, 0, "")
	pdf.CellFormat(65, 7, "File", "1", 1, "L", true, 0, "")

	// Details table data
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(0, 0, 0)
	for i, obj := range summary.Objects {
		// Check if we need a new page
		if pdf.GetY() > 270 {
			pdf.AddPage()
			// Repeat header on new page
			pdf.SetFont("Arial", "B", 9)
			pdf.SetFillColor(68, 114, 196)
			pdf.SetTextColor(255, 255, 255)
			pdf.CellFormat(35, 7, "Type", "1", 0, "L", true, 0, "")
			pdf.CellFormat(20, 7, "ID", "1", 0, "C", true, 0, "")
			pdf.CellFormat(70, 7, "Name", "1", 0, "L", true, 0, "")
			pdf.CellFormat(65, 7, "File", "1", 1, "L", true, 0, "")
			pdf.SetFont("Arial", "", 8)
			pdf.SetTextColor(0, 0, 0)
		}

		fill := i%2 == 0
		if fill {
			pdf.SetFillColor(245, 245, 245)
		}

		// Truncate file path if too long
		filePath := obj.FilePath
		if len(filePath) > 40 {
			filePath = "..." + filePath[len(filePath)-37:]
		}

		// Truncate name if too long
		name := obj.Name
		if len(name) > 35 {
			name = name[:32] + "..."
		}

		pdf.CellFormat(35, 6, obj.Type, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(20, 6, obj.ID, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(70, 6, name, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(65, 6, filePath, "1", 1, "L", fill, 0, "")
	}

	return pdf.OutputFileAndClose(filePath)
}
