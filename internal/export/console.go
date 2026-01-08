package export

import (
	"fmt"
	"strings"

	"github.com/andrijan/bc-objects-counter/internal/counter"
)

// ToConsole formats the summary for console output.
func ToConsole(summary *counter.Summary) string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString("═══════════════════════════════════════════\n")
	sb.WriteString("       BC Objects Summary\n")
	sb.WriteString("═══════════════════════════════════════════\n\n")

	// Find max type name length for alignment (minimum 5 for "TOTAL")
	maxLen := 5
	for _, c := range summary.CountsByType {
		if len(c.Type) > maxLen {
			maxLen = len(c.Type)
		}
	}

	// Print counts by type
	for _, c := range summary.CountsByType {
		padding := strings.Repeat(" ", maxLen-len(c.Type))
		sb.WriteString(fmt.Sprintf("  %s%s : %d\n", c.Type, padding, c.Count))
	}

	sb.WriteString("\n───────────────────────────────────────────\n")
	padding := strings.Repeat(" ", maxLen-5)
	sb.WriteString(fmt.Sprintf("  TOTAL%s : %d\n", padding, summary.TotalObjects))
	sb.WriteString("═══════════════════════════════════════════\n")

	return sb.String()
}
