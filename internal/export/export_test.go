package export

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andrijan/bc-objects-counter/internal/counter"
	"github.com/andrijan/bc-objects-counter/internal/scanner"
)

func createTestSummary() *counter.Summary {
	objects := []scanner.BCObject{
		{Type: "table", ID: "50100", Name: "Test Table", FilePath: "src/table.al"},
		{Type: "codeunit", ID: "50100", Name: "Test Codeunit", FilePath: "src/codeunit.al"},
		{Type: "page", ID: "50100", Name: "Test Page", FilePath: "src/page.al"},
	}
	return counter.CountObjects(objects)
}

func TestToJSONString(t *testing.T) {
	summary := createTestSummary()

	jsonStr, err := ToJSONString(summary)
	if err != nil {
		t.Fatal(err)
	}

	// Check that JSON contains expected fields
	if !strings.Contains(jsonStr, `"totalObjects": 3`) {
		t.Error("JSON should contain totalObjects: 3")
	}
	if !strings.Contains(jsonStr, `"table"`) {
		t.Error("JSON should contain table type")
	}
	if !strings.Contains(jsonStr, `"Test Table"`) {
		t.Error("JSON should contain Test Table name")
	}
}

func TestToJSON(t *testing.T) {
	summary := createTestSummary()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")

	err := ToJSON(summary, filePath)
	if err != nil {
		t.Fatal(err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("JSON file was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), `"totalObjects": 3`) {
		t.Error("JSON file should contain totalObjects: 3")
	}
}

func TestToConsole(t *testing.T) {
	summary := createTestSummary()

	output := ToConsole(summary)

	// Check that console output contains expected content
	if !strings.Contains(output, "BC Objects Summary") {
		t.Error("console output should contain title")
	}
	if !strings.Contains(output, "TOTAL") {
		t.Error("console output should contain TOTAL")
	}
	if !strings.Contains(output, "3") {
		t.Error("console output should contain count 3")
	}
}

func TestToExcel(t *testing.T) {
	summary := createTestSummary()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.xlsx")

	err := ToExcel(summary, filePath)
	if err != nil {
		t.Fatal(err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Excel file was not created")
	}

	// Verify file has content (non-empty)
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Error("Excel file is empty")
	}
}

func TestToPDF(t *testing.T) {
	summary := createTestSummary()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.pdf")

	err := ToPDF(summary, filePath)
	if err != nil {
		t.Fatal(err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("PDF file was not created")
	}

	// Verify file has content (non-empty)
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Error("PDF file is empty")
	}
}

func TestToConsoleEmpty(t *testing.T) {
	summary := counter.CountObjects([]scanner.BCObject{})

	output := ToConsole(summary)

	if !strings.Contains(output, "TOTAL") {
		t.Error("console output should contain TOTAL even for empty summary")
	}
	if !strings.Contains(output, "0") {
		t.Error("console output should contain 0 for empty summary")
	}
}
