package counter

import (
	"testing"

	"github.com/andrijan/bc-objects-counter/internal/scanner"
)

func TestCountObjects(t *testing.T) {
	objects := []scanner.BCObject{
		{Type: "table", ID: "50100", Name: "Table 1", FilePath: "file1.al"},
		{Type: "table", ID: "50101", Name: "Table 2", FilePath: "file2.al"},
		{Type: "codeunit", ID: "50100", Name: "Codeunit 1", FilePath: "file3.al"},
		{Type: "page", ID: "50100", Name: "Page 1", FilePath: "file4.al"},
		{Type: "page", ID: "50101", Name: "Page 2", FilePath: "file5.al"},
		{Type: "page", ID: "50102", Name: "Page 3", FilePath: "file6.al"},
	}

	summary := CountObjects(objects)

	// Test total count
	if summary.TotalObjects != 6 {
		t.Errorf("expected TotalObjects 6, got %d", summary.TotalObjects)
	}

	// Test counts by type
	if len(summary.CountsByType) != 3 {
		t.Errorf("expected 3 object types, got %d", len(summary.CountsByType))
	}

	// Since sorted by count desc, page should be first
	if summary.CountsByType[0].Type != "page" || summary.CountsByType[0].Count != 3 {
		t.Errorf("expected first type to be page with count 3, got %s with count %d",
			summary.CountsByType[0].Type, summary.CountsByType[0].Count)
	}

	// Test ObjectsByType
	if len(summary.ObjectsByType["table"]) != 2 {
		t.Errorf("expected 2 tables, got %d", len(summary.ObjectsByType["table"]))
	}
	if len(summary.ObjectsByType["codeunit"]) != 1 {
		t.Errorf("expected 1 codeunit, got %d", len(summary.ObjectsByType["codeunit"]))
	}
	if len(summary.ObjectsByType["page"]) != 3 {
		t.Errorf("expected 3 pages, got %d", len(summary.ObjectsByType["page"]))
	}
}

func TestCountObjectsEmpty(t *testing.T) {
	objects := []scanner.BCObject{}

	summary := CountObjects(objects)

	if summary.TotalObjects != 0 {
		t.Errorf("expected TotalObjects 0, got %d", summary.TotalObjects)
	}
	if len(summary.CountsByType) != 0 {
		t.Errorf("expected 0 object types, got %d", len(summary.CountsByType))
	}
}

func TestGetObjectsByType(t *testing.T) {
	objects := []scanner.BCObject{
		{Type: "table", ID: "50100", Name: "Table 1", FilePath: "file1.al"},
		{Type: "codeunit", ID: "50100", Name: "Codeunit 1", FilePath: "file2.al"},
	}

	summary := CountObjects(objects)

	tables := summary.GetObjectsByType("table")
	if len(tables) != 1 {
		t.Errorf("expected 1 table, got %d", len(tables))
	}

	// Non-existent type
	reports := summary.GetObjectsByType("report")
	if len(reports) != 0 {
		t.Errorf("expected 0 reports, got %d", len(reports))
	}
}

func TestGetCountByType(t *testing.T) {
	objects := []scanner.BCObject{
		{Type: "table", ID: "50100", Name: "Table 1", FilePath: "file1.al"},
		{Type: "table", ID: "50101", Name: "Table 2", FilePath: "file2.al"},
		{Type: "codeunit", ID: "50100", Name: "Codeunit 1", FilePath: "file3.al"},
	}

	summary := CountObjects(objects)

	if count := summary.GetCountByType("table"); count != 2 {
		t.Errorf("expected table count 2, got %d", count)
	}

	if count := summary.GetCountByType("codeunit"); count != 1 {
		t.Errorf("expected codeunit count 1, got %d", count)
	}

	// Non-existent type
	if count := summary.GetCountByType("report"); count != 0 {
		t.Errorf("expected report count 0, got %d", count)
	}
}

func TestCountObjectsSorting(t *testing.T) {
	// Create objects where two types have the same count
	objects := []scanner.BCObject{
		{Type: "table", ID: "50100", Name: "Table 1", FilePath: "file1.al"},
		{Type: "codeunit", ID: "50100", Name: "Codeunit 1", FilePath: "file2.al"},
		{Type: "page", ID: "50100", Name: "Page 1", FilePath: "file3.al"},
		{Type: "page", ID: "50101", Name: "Page 2", FilePath: "file4.al"},
	}

	summary := CountObjects(objects)

	// Page should be first (count 2)
	if summary.CountsByType[0].Type != "page" {
		t.Errorf("expected first type to be page, got %s", summary.CountsByType[0].Type)
	}

	// Codeunit and table both have count 1, should be sorted alphabetically
	// codeunit comes before table
	if summary.CountsByType[1].Type != "codeunit" {
		t.Errorf("expected second type to be codeunit, got %s", summary.CountsByType[1].Type)
	}
	if summary.CountsByType[2].Type != "table" {
		t.Errorf("expected third type to be table, got %s", summary.CountsByType[2].Type)
	}
}
