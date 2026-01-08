// Package export provides functionality to export BC object summaries to various formats.
package export

import (
	"encoding/json"
	"os"

	"github.com/andrijan/bc-objects-counter/internal/counter"
)

// ToJSON exports the summary to a JSON file.
func ToJSON(summary *counter.Summary, filePath string) error {
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// ToJSONString returns the summary as a JSON string.
func ToJSONString(summary *counter.Summary) (string, error) {
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
