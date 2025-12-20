# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`st2` is a Go library and CLI tool that converts between various structured data formats:
- **Input formats**: JSON, YAML, CSV, XML, TOML, Protocol Buffers, Thrift, Go structs
- **Output formats**: Go structs, Protocol Buffers, Thrift

## Common Development Commands

### Build and Install
```bash
# Build the CLI tool
go build -o st2 ./cmd/st2

# Install globally
go install ./cmd/st2@latest

# Run the CLI tool
./st2 --help
```

### Testing
```bash
# Run all tests with coverage
go test -v -cover ./...

# Run tests for specific file
go test -v ./... -run TestJSONParser

# Run specific test
go test -v -run TestConvert
```

### Code Quality
```bash
# Format all code
go fmt ./...

# Lint (if golangci-lint is installed)
golangci-lint run
```

## Architecture

### Core Library (`/` root directory)
The main library follows a **factory pattern** with these key components:

1. **Entry Point**: `st2.go` - `Convert()` function orchestrates the conversion
   - Creates parser via `CreateParser(ctx)`
   - Creates template via `CreateTmpl(ctx)`
   - Creates formatter via `CreateFormater(ctx)`
   - Executes template with parsed structs

2. **Parser Interface**: `parser.go`
   - `Parse(r io.Reader) ([]*Struct, error)` - parses source to `Struct` objects

3. **Data Model**: `struct.go` and `type.go`
   - `Struct` - represents a struct/enum with members
   - `Member` - represents a field with type, tags, comments
   - `Type` interface - handles type conversions between formats (Go, Proto, Thrift, JSON)

4. **Structured Parser**: `structured_parser.go`
   - Generic parser for JSON/YAML/TOML/XML/CSV
   - Uses `UnmarshalTagFormat` interface for format-specific unmarshaling
   - Builds `rawNode` tree, then converts to `Struct` objects
   - Handles deduplication via fingerprinting

5. **Specialized Parsers**:
   - `json_parser.go` - uses json-iterator with `UseNumber` for precision
   - `yaml_parser.go` - YAML parsing
   - `toml_parser.go` - TOML parsing
   - `xml_parser.go` - XML parsing with attribute/content tag prefixes
   - `csv_parser.go` - CSV parsing
   - `go_parser.go` - Go source parsing (uses AST)
   - `proto_parser.go` - Protocol Buffers parsing
   - `thrift_parser.go` - Thrift parsing

6. **Templates**: `/tmpl/` directory
   - `go.go` - Go struct/enum template
   - `proto.go` - Protocol Buffers message/enum template
   - `thrift.go` - Thrift struct/union/enum template

7. **Formatters**: `format.go`
   - `GoFormater` - uses `go/format` package
   - `EmptyFormater` - for other formats

### CLI Tool (`cmd/st2/`)

1. **main.go** - CLI entry using `urfave/cli/v3`
   - Flags: `--src`, `--dst`, `--input`, `--output`, `--root`, `--prefix`, `--suffix`
   - Clipboard support: `--rc` (read), `--wc` (write)
   - XML-specific: `--xml-content-tag-prefix`, `--xml-attribute-tag-prefix`

2. **clipboard.go** - Clipboard I/O using `golang.design/x/clipboard`

3. **config/vars.go** - Version configuration

### Key Design Patterns

1. **Factory Pattern**: `factory.go` creates parsers/templates/formatters based on context
2. **Fingerprinting**: `raw_node.go` - deduplicates structs by structural similarity
3. **Template-based Code Generation**: Uses Go's `text/template` with nested templates
4. **Type System**: `type.go` defines a unified type interface with format-specific methods

### Data Flow

```
Input → Parser → []*Struct → Template → Formatted Output
                    ↓
              Type Conversion
              (Go/Proto/Thrift)
```

## Testing Patterns

- Tests are co-located with source files (`*_test.go`)
- Uses `testify/assert` for assertions
- Table-driven tests with `name`, `args`, `want*` fields
- Test data in `cmd/st2/testdata/`

## Supported Conversions

| Source → Destination | Go | Proto | Thrift |
|---------------------|----|-------|--------|
| JSON | ✓ | ✓ | ✓ |
| YAML | ✓ | ✓ | ✓ |
| CSV | ✓ | ✓ | ✓ |
| XML | ✓ | ✓ | ✓ |
| TOML | ✓ | ✓ | ✓ |
| Go | ✓ | ✓ | ✓ |
| Proto | ✓ | ✓ | ✓ |
| Thrift | ✓ | ✓ | ✓ |

## Important Constants

- `Lang*` constants in `const.go` define supported languages
- `RootDefault = "Root"` - default root struct name
- `FlagXMLAttributeTagPrefixDefault = ","` - default XML attribute prefix

## CI/CD

- GitHub Actions: `.github/workflows/`
- Uses reusable workflows from `tenfyzhong/workflows`
- Code coverage via Codecov
- Homebrew releases via workflow

## Notes for Future Development

1. **Adding new source formats**: Implement `Parse` interface, add to `CreateParser()` factory
2. **Adding new destination formats**: Add template to `/tmpl/`, add to `CreateTmpl()` factory
3. **Type mappings**: Update `type.go` with new format methods
4. **CLI flags**: Add flags in `main.go`, wire through `Context` and `st2.Convert()`
5. **XML handling**: Uses `goxml2json` library, supports attribute/content tag prefixes
6. **Clipboard**: Platform-specific, uses `golang.design/x/clipboard`
