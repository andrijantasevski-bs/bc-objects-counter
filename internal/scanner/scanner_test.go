package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseObjectLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected *BCObject
	}{
		{
			name: "table",
			line: `table 50100 "My Custom Table"`,
			expected: &BCObject{
				Type: "table",
				ID:   "50100",
				Name: "My Custom Table",
			},
		},
		{
			name: "codeunit",
			line: `codeunit 50100 "My Codeunit"`,
			expected: &BCObject{
				Type: "codeunit",
				ID:   "50100",
				Name: "My Codeunit",
			},
		},
		{
			name: "page",
			line: `page 50100 "My Page"`,
			expected: &BCObject{
				Type: "page",
				ID:   "50100",
				Name: "My Page",
			},
		},
		{
			name: "report",
			line: `report 50100 "My Report"`,
			expected: &BCObject{
				Type: "report",
				ID:   "50100",
				Name: "My Report",
			},
		},
		{
			name: "tableextension",
			line: `tableextension 50100 "My Table Ext" extends "Customer"`,
			expected: &BCObject{
				Type: "tableextension",
				ID:   "50100",
				Name: "My Table Ext",
			},
		},
		{
			name: "pageextension",
			line: `pageextension 50100 "My Page Ext" extends "Customer Card"`,
			expected: &BCObject{
				Type: "pageextension",
				ID:   "50100",
				Name: "My Page Ext",
			},
		},
		{
			name: "enum",
			line: `enum 50100 "My Enum"`,
			expected: &BCObject{
				Type: "enum",
				ID:   "50100",
				Name: "My Enum",
			},
		},
		{
			name: "enumextension",
			line: `enumextension 50100 "My Enum Ext" extends "Sales Document Type"`,
			expected: &BCObject{
				Type: "enumextension",
				ID:   "50100",
				Name: "My Enum Ext",
			},
		},
		{
			name: "xmlport",
			line: `xmlport 50100 "My XMLport"`,
			expected: &BCObject{
				Type: "xmlport",
				ID:   "50100",
				Name: "My XMLport",
			},
		},
		{
			name: "query",
			line: `query 50100 "My Query"`,
			expected: &BCObject{
				Type: "query",
				ID:   "50100",
				Name: "My Query",
			},
		},
		{
			name: "interface",
			line: `interface "IMyInterface"`,
			expected: &BCObject{
				Type: "interface",
				ID:   "",
				Name: "IMyInterface",
			},
		},
		{
			name: "permissionset",
			line: `permissionset 50100 "My Permission Set"`,
			expected: &BCObject{
				Type: "permissionset",
				ID:   "50100",
				Name: "My Permission Set",
			},
		},
		{
			name: "case insensitive - TABLE",
			line: `TABLE 50100 "My Table"`,
			expected: &BCObject{
				Type: "table",
				ID:   "50100",
				Name: "My Table",
			},
		},
		{
			name: "case insensitive - CodeUnit",
			line: `CodeUnit 50100 "My Codeunit"`,
			expected: &BCObject{
				Type: "codeunit",
				ID:   "50100",
				Name: "My Codeunit",
			},
		},
		{
			name: "with leading whitespace",
			line: `    table 50100 "My Table"`,
			expected: &BCObject{
				Type: "table",
				ID:   "50100",
				Name: "My Table",
			},
		},
		{
			name:     "not an object - comment",
			line:     `// table 50100 "Commented Table"`,
			expected: nil, // This is filtered at ScanFile level
		},
		{
			name:     "not an object - random code",
			line:     `var myVar: Integer;`,
			expected: nil,
		},
		{
			name:     "not an object - empty line",
			line:     ``,
			expected: nil,
		},
		{
			name: "profile without ID",
			line: `profile "My Profile"`,
			expected: &BCObject{
				Type: "profile",
				ID:   "",
				Name: "My Profile",
			},
		},
		{
			name: "controladdin without ID",
			line: `controladdin "My Control Addin"`,
			expected: &BCObject{
				Type: "controladdin",
				ID:   "",
				Name: "My Control Addin",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseObjectLine(tt.line, "test.al")

			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %+v", result)
				}
				return
			}

			if result == nil {
				t.Errorf("expected %+v, got nil", tt.expected)
				return
			}

			if result.Type != tt.expected.Type {
				t.Errorf("Type: expected %s, got %s", tt.expected.Type, result.Type)
			}
			if result.ID != tt.expected.ID {
				t.Errorf("ID: expected %s, got %s", tt.expected.ID, result.ID)
			}
			if result.Name != tt.expected.Name {
				t.Errorf("Name: expected %s, got %s", tt.expected.Name, result.Name)
			}
		})
	}
}

func TestScanFile(t *testing.T) {
	// Create a temporary AL file for testing
	content := `// This is a comment
table 50100 "Customer Extension"
{
    // Some table code
}

/* 
Block comment
codeunit 99999 "Should Be Ignored"
*/

codeunit 50101 "Sales Management"
{
    // Codeunit code
}

page 50100 "Customer Card Ext"
{
    // Page code
}
`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.al")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	objects, err := ScanFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(objects) != 3 {
		t.Errorf("expected 3 objects, got %d", len(objects))
	}

	// Verify objects
	expectedTypes := []string{"table", "codeunit", "page"}
	expectedIDs := []string{"50100", "50101", "50100"}
	expectedNames := []string{"Customer Extension", "Sales Management", "Customer Card Ext"}

	for i, obj := range objects {
		if obj.Type != expectedTypes[i] {
			t.Errorf("object %d: expected type %s, got %s", i, expectedTypes[i], obj.Type)
		}
		if obj.ID != expectedIDs[i] {
			t.Errorf("object %d: expected ID %s, got %s", i, expectedIDs[i], obj.ID)
		}
		if obj.Name != expectedNames[i] {
			t.Errorf("object %d: expected name %s, got %s", i, expectedNames[i], obj.Name)
		}
	}
}

func TestScanDirectory(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()

	// Create main directory with an AL file
	mainFile := filepath.Join(tmpDir, "main.al")
	mainContent := `table 50100 "Main Table"
{
}
`
	if err := os.WriteFile(mainFile, []byte(mainContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create subdirectory with AL files
	subDir := filepath.Join(tmpDir, "sub")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	subFile := filepath.Join(subDir, "sub.al")
	subContent := `codeunit 50100 "Sub Codeunit"
{
}

page 50100 "Sub Page"
{
}
`
	if err := os.WriteFile(subFile, []byte(subContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Test recursive scan
	objects, err := ScanDirectory(tmpDir, true)
	if err != nil {
		t.Fatal(err)
	}

	if len(objects) != 3 {
		t.Errorf("recursive: expected 3 objects, got %d", len(objects))
	}

	// Test non-recursive scan
	objects, err = ScanDirectory(tmpDir, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(objects) != 1 {
		t.Errorf("non-recursive: expected 1 object, got %d", len(objects))
	}
}

func TestGetSupportedObjectTypes(t *testing.T) {
	types := GetSupportedObjectTypes()

	if len(types) == 0 {
		t.Error("expected at least one supported object type")
	}

	// Check that common types are included
	expectedTypes := []string{"table", "page", "codeunit", "report", "enum"}
	for _, expected := range expectedTypes {
		found := false
		for _, objType := range types {
			if objType == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected type %s not found in supported types", expected)
		}
	}
}
