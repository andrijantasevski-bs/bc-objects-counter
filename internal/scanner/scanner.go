// Package scanner provides functionality to scan directories for Business Central AL files
// and extract object declarations.
package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// BCObject represents a Business Central object found in an AL file.
type BCObject struct {
	Type     string `json:"type"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	FilePath string `json:"filePath"`
}

// objectPattern matches BC object declarations.
// Captures: 1=type, 2=id (optional for some types), 3=name
var objectPattern = regexp.MustCompile(`(?i)^\s*(table|tableextension|page|pageextension|report|reportextension|codeunit|xmlport|query|enum|enumextension|interface|permissionset|permissionsetextension|profile|controladdin|entitlement)\s+(\d+)?\s*"([^"]+)"`)

// ScanDirectory recursively scans a directory for .al files and extracts BC objects.
func ScanDirectory(root string, recursive bool) ([]BCObject, error) {
	var objects []BCObject

	walkFn := func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories if not recursive (except root)
		if d.IsDir() {
			if !recursive && path != root {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process .al files
		if !strings.HasSuffix(strings.ToLower(path), ".al") {
			return nil
		}

		fileObjects, err := ScanFile(path)
		if err != nil {
			// Log error but continue scanning other files
			return nil
		}

		objects = append(objects, fileObjects...)
		return nil
	}

	err := filepath.WalkDir(root, walkFn)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

// ScanFile scans a single AL file and extracts BC objects.
func ScanFile(filePath string) ([]BCObject, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var objects []BCObject
	scanner := bufio.NewScanner(file)
	inBlockComment := false

	for scanner.Scan() {
		line := scanner.Text()

		// Handle block comments
		if strings.Contains(line, "/*") {
			inBlockComment = true
		}
		if strings.Contains(line, "*/") {
			inBlockComment = false
			continue
		}
		if inBlockComment {
			continue
		}

		// Skip single-line comments
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		// Try to match object declaration
		if obj := ParseObjectLine(line, filePath); obj != nil {
			objects = append(objects, *obj)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}

// ParseObjectLine attempts to parse a BC object declaration from a single line.
func ParseObjectLine(line, filePath string) *BCObject {
	matches := objectPattern.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	objType := strings.ToLower(matches[1])
	objID := matches[2] // May be empty for interface, profile, controladdin, entitlement
	objName := matches[3]

	return &BCObject{
		Type:     objType,
		ID:       objID,
		Name:     objName,
		FilePath: filePath,
	}
}

// GetSupportedObjectTypes returns a list of all supported BC object types.
func GetSupportedObjectTypes() []string {
	return []string{
		"table",
		"tableextension",
		"page",
		"pageextension",
		"report",
		"reportextension",
		"codeunit",
		"xmlport",
		"query",
		"enum",
		"enumextension",
		"interface",
		"permissionset",
		"permissionsetextension",
		"profile",
		"controladdin",
		"entitlement",
	}
}
