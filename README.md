# BC Objects Counter

A fast CLI tool to scan Business Central AL files and count object types.

## Features

- 🔍 Recursively scans directories for `.al` files
- 📊 Counts all BC object types (tables, pages, codeunits, reports, etc.)
- 📁 Exports to JSON, Excel (.xlsx), and PDF formats
- ⚡ Fast, single binary with no dependencies
- 🖥️ Cross-platform: Windows, Linux, macOS

## Installation

### Download Binary

Download the latest release from the [Releases page](https://github.com/andrijan/bc-objects-counter/releases).

### Build from Source

```bash
go install github.com/andrijan/bc-objects-counter@latest
```

Or clone and build:

```bash
git clone https://github.com/andrijan/bc-objects-counter.git
cd bc-objects-counter
go build -o bc-objects-counter
```

## Usage

```bash
# Scan current directory and print to console
bc-objects-counter

# Scan a specific directory
bc-objects-counter /path/to/al/files

# Export to JSON
bc-objects-counter -o json

# Export to Excel
bc-objects-counter -o xlsx

# Export to PDF
bc-objects-counter -o pdf

# Export all formats at once
bc-objects-counter -o all

# Specify output filename
bc-objects-counter -o xlsx -f my-report

# Non-recursive scan (current directory only)
bc-objects-counter -r=false

# Show verbose output
bc-objects-counter -v
```

## Supported Object Types

- `table`
- `tableextension`
- `page`
- `pageextension`
- `report`
- `reportextension`
- `codeunit`
- `xmlport`
- `query`
- `enum`
- `enumextension`
- `interface`
- `permissionset`
- `permissionsetextension`
- `profile`
- `controladdin`
- `entitlement`

## Output Example

### Console Output

```
═══════════════════════════════════════════
       BC Objects Summary
═══════════════════════════════════════════

  page              : 45
  codeunit          : 32
  table             : 28
  tableextension    : 15
  pageextension     : 12
  enum              : 8
  report            : 5

───────────────────────────────────────────
  TOTAL             : 145
═══════════════════════════════════════════
```

## Development

### Run Tests

```bash
go test -v ./...
```

### Build

```bash
go build -o bc-objects-counter
```

### Release

Create a new tag to trigger a release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## License

MIT License
