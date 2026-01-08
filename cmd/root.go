package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andrijan/bc-objects-counter/internal/counter"
	"github.com/andrijan/bc-objects-counter/internal/export"
	"github.com/andrijan/bc-objects-counter/internal/scanner"
	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags
var Version = "dev"

var (
	outputFormat string
	outputFile   string
	recursive    bool
	verbose      bool
)

var rootCmd = &cobra.Command{
	Use:   "bc-objects-counter [path]",
	Short: "Count Business Central objects in AL files",
	Long: `BC Objects Counter scans a directory for Business Central AL files
and counts all object types (tables, pages, codeunits, etc.).

It can export the results to JSON, Excel, or PDF format.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCounter,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&outputFormat, "output", "o", "console", "Output format: console, json, xlsx, pdf, all")
	rootCmd.Flags().StringVarP(&outputFile, "file", "f", "", "Output filename (without extension, auto-generated if not specified)")
	rootCmd.Flags().BoolVarP(&recursive, "recursive", "r", true, "Scan subdirectories")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed output")
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("bc-objects-counter version {{.Version}}\n")
}

func runCounter(cmd *cobra.Command, args []string) error {
	// Determine scan path
	scanPath := "."
	if len(args) > 0 {
		scanPath = args[0]
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(scanPath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Verify path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", absPath)
	}

	if verbose {
		fmt.Printf("Scanning: %s\n", absPath)
		fmt.Printf("Recursive: %v\n", recursive)
	}

	// Scan for objects
	objects, err := scanner.ScanDirectory(absPath, recursive)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	if verbose {
		fmt.Printf("Found %d objects\n", len(objects))
	}

	// Create summary
	summary := counter.CountObjects(objects)

	// Generate output filename if not specified
	if outputFile == "" {
		timestamp := time.Now().Format("20060102-150405")
		outputFile = fmt.Sprintf("bc-objects-%s", timestamp)
	}

	// Handle output formats
	format := strings.ToLower(outputFormat)

	switch format {
	case "console":
		fmt.Print(export.ToConsole(summary))

	case "json":
		jsonFile := outputFile + ".json"
		if err := export.ToJSON(summary, jsonFile); err != nil {
			return fmt.Errorf("failed to export JSON: %w", err)
		}
		fmt.Printf("✓ Exported to %s\n", jsonFile)

	case "xlsx", "excel":
		xlsxFile := outputFile + ".xlsx"
		if err := export.ToExcel(summary, xlsxFile); err != nil {
			return fmt.Errorf("failed to export Excel: %w", err)
		}
		fmt.Printf("✓ Exported to %s\n", xlsxFile)

	case "pdf":
		pdfFile := outputFile + ".pdf"
		if err := export.ToPDF(summary, pdfFile); err != nil {
			return fmt.Errorf("failed to export PDF: %w", err)
		}
		fmt.Printf("✓ Exported to %s\n", pdfFile)

	case "all":
		// Print console output
		fmt.Print(export.ToConsole(summary))

		// Export all formats
		jsonFile := outputFile + ".json"
		if err := export.ToJSON(summary, jsonFile); err != nil {
			return fmt.Errorf("failed to export JSON: %w", err)
		}
		fmt.Printf("✓ Exported to %s\n", jsonFile)

		xlsxFile := outputFile + ".xlsx"
		if err := export.ToExcel(summary, xlsxFile); err != nil {
			return fmt.Errorf("failed to export Excel: %w", err)
		}
		fmt.Printf("✓ Exported to %s\n", xlsxFile)

		pdfFile := outputFile + ".pdf"
		if err := export.ToPDF(summary, pdfFile); err != nil {
			return fmt.Errorf("failed to export PDF: %w", err)
		}
		fmt.Printf("✓ Exported to %s\n", pdfFile)

	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}

	return nil
}
