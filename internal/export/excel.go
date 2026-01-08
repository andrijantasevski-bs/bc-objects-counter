package export

import (
	"fmt"

	"github.com/andrijan/bc-objects-counter/internal/counter"
	"github.com/xuri/excelize/v2"
)

// ToExcel exports the summary to an Excel file.
func ToExcel(summary *counter.Summary, filePath string) error {
	f := excelize.NewFile()
	defer f.Close()

	// Create Summary sheet
	summarySheet := "Summary"
	f.SetSheetName("Sheet1", summarySheet)

	// Summary headers
	f.SetCellValue(summarySheet, "A1", "BC Objects Summary")
	f.SetCellValue(summarySheet, "A3", "Object Type")
	f.SetCellValue(summarySheet, "B3", "Count")

	// Style for headers
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#4472C4"}, Pattern: 1},
	})
	f.SetCellStyle(summarySheet, "A3", "B3", headerStyle)

	// Title style
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14},
	})
	f.SetCellStyle(summarySheet, "A1", "A1", titleStyle)

	// Write counts by type
	row := 4
	for _, c := range summary.CountsByType {
		f.SetCellValue(summarySheet, fmt.Sprintf("A%d", row), c.Type)
		f.SetCellValue(summarySheet, fmt.Sprintf("B%d", row), c.Count)
		row++
	}

	// Total row
	row++
	f.SetCellValue(summarySheet, fmt.Sprintf("A%d", row), "TOTAL")
	f.SetCellValue(summarySheet, fmt.Sprintf("B%d", row), summary.TotalObjects)
	totalStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	f.SetCellStyle(summarySheet, fmt.Sprintf("A%d", row), fmt.Sprintf("B%d", row), totalStyle)

	// Set column widths
	f.SetColWidth(summarySheet, "A", "A", 25)
	f.SetColWidth(summarySheet, "B", "B", 12)

	// Create Details sheet
	detailsSheet := "Details"
	f.NewSheet(detailsSheet)

	// Details headers
	f.SetCellValue(detailsSheet, "A1", "Type")
	f.SetCellValue(detailsSheet, "B1", "ID")
	f.SetCellValue(detailsSheet, "C1", "Name")
	f.SetCellValue(detailsSheet, "D1", "File Path")
	f.SetCellStyle(detailsSheet, "A1", "D1", headerStyle)

	// Write all objects
	row = 2
	for _, obj := range summary.Objects {
		f.SetCellValue(detailsSheet, fmt.Sprintf("A%d", row), obj.Type)
		f.SetCellValue(detailsSheet, fmt.Sprintf("B%d", row), obj.ID)
		f.SetCellValue(detailsSheet, fmt.Sprintf("C%d", row), obj.Name)
		f.SetCellValue(detailsSheet, fmt.Sprintf("D%d", row), obj.FilePath)
		row++
	}

	// Set column widths for details
	f.SetColWidth(detailsSheet, "A", "A", 20)
	f.SetColWidth(detailsSheet, "B", "B", 10)
	f.SetColWidth(detailsSheet, "C", "C", 40)
	f.SetColWidth(detailsSheet, "D", "D", 60)

	return f.SaveAs(filePath)
}
