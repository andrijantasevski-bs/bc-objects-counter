// Package counter provides functionality to aggregate and count BC objects.
package counter

import (
	"sort"

	"github.com/andrijan/bc-objects-counter/internal/scanner"
)

// ObjectCount represents the count for a specific object type.
type ObjectCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// Summary contains the aggregated results of scanning BC objects.
type Summary struct {
	TotalObjects  int                           `json:"totalObjects"`
	CountsByType  []ObjectCount                 `json:"countsByType"`
	Objects       []scanner.BCObject            `json:"objects"`
	ObjectsByType map[string][]scanner.BCObject `json:"objectsByType"`
}

// CountObjects aggregates the scanned BC objects into a summary.
func CountObjects(objects []scanner.BCObject) *Summary {
	summary := &Summary{
		TotalObjects:  len(objects),
		Objects:       objects,
		ObjectsByType: make(map[string][]scanner.BCObject),
	}

	// Group objects by type
	typeCounts := make(map[string]int)
	for _, obj := range objects {
		typeCounts[obj.Type]++
		summary.ObjectsByType[obj.Type] = append(summary.ObjectsByType[obj.Type], obj)
	}

	// Convert to sorted slice
	for objType, count := range typeCounts {
		summary.CountsByType = append(summary.CountsByType, ObjectCount{
			Type:  objType,
			Count: count,
		})
	}

	// Sort by count (descending), then by type name
	sort.Slice(summary.CountsByType, func(i, j int) bool {
		if summary.CountsByType[i].Count != summary.CountsByType[j].Count {
			return summary.CountsByType[i].Count > summary.CountsByType[j].Count
		}
		return summary.CountsByType[i].Type < summary.CountsByType[j].Type
	})

	return summary
}

// GetObjectsByType returns all objects of a specific type from the summary.
func (s *Summary) GetObjectsByType(objType string) []scanner.BCObject {
	return s.ObjectsByType[objType]
}

// GetCountByType returns the count for a specific object type.
func (s *Summary) GetCountByType(objType string) int {
	for _, c := range s.CountsByType {
		if c.Type == objType {
			return c.Count
		}
	}
	return 0
}
