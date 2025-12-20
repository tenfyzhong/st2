# Repository Guidelines

## Project Structure & Module Organization

This repository contains the `st2` Go library and CLI tool for parsing and generating structured data formats.

- **Root Directory**: Core library code (parsers, generators, utilities)
- `cmd/st2/`: CLI tool implementation
- `tmpl/`: Code generation templates
- Tests: Co-located with source files (e.g., `json_parser_test.go` for `json_parser.go`)

## Build, Test, and Development Commands

### CLI Tool Development
```bash
# Build the CLI tool
go build -o st2 ./cmd/st2

# Install the CLI tool globally
go install ./cmd/st2@latest

# Run the CLI tool
./st2 --help
```

### Library Development
```bash
# Run all tests with coverage
go test -v -cover ./...

# Run tests for a specific file
go test -v ./... -run TestJSONParser
```

## Coding Style & Naming Conventions

- Follow Go's [official coding conventions](https://go.dev/doc/effective_go)
- Use `go fmt ./...` to format code consistently
- Naming: PascalCase for exported types/functions, camelCase for private ones
- Filenames: lowercase with underscores (e.g., `json_parser.go` for JSON parser)

## Testing Guidelines

- Tests are written in Go's standard testing framework
- Test files follow `*_test.go` naming pattern
- Strive for >90% test coverage
- Test cases should cover edge cases and expected behavior

## Commit & Pull Request Guidelines

### Commits
- Follow [conventional commit messages](https://www.conventionalcommits.org/):
  - `feat(parser): Add XML parser support`
  - `fix(cli): Resolve output file path issue`
  - `docs(readme): Update installation instructions`

### Pull Requests
- Include a clear description of changes
- Link related issues if applicable
- Ensure all tests pass
- Maintain or improve test coverage
