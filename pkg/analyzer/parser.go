package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Parser handles parsing Go source files
type Parser struct {
	fset      *token.FileSet
	parseMode parser.Mode
}

// NewParser creates a new parser
func NewParser() *Parser {
	return &Parser{
		fset:      token.NewFileSet(),
		parseMode: parser.ParseComments | parser.SpuriousErrors,
	}
}

// FileSet returns the token.FileSet used by the parser
func (p *Parser) FileSet() *token.FileSet {
	return p.fset
}

// ParseFile parses a single Go source file
func (p *Parser) ParseFile(path string) (*ast.File, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	// Parse the file
	file, err := parser.ParseFile(p.fset, path, nil, p.parseMode)
	if err != nil {
		return nil, fmt.Errorf("parse error in %s: %w", path, err)
	}

	return file, nil
}

// ParseString parses Go source from a string
func (p *Parser) ParseString(src string) (*ast.File, error) {
	file, err := parser.ParseFile(p.fset, "", strings.NewReader(src), p.parseMode)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// IsGoFile checks if a file is a Go source file
func IsGoFile(path string) bool {
	return strings.HasSuffix(path, ".go")
}

// IsTestFile checks if a file is a Go test file
func IsTestFile(path string) bool {
	base := filepath.Base(path)
	return strings.HasSuffix(base, "_test.go")
}

// ShouldExclude checks if a file should be excluded from analysis
func ShouldExclude(path string, excludeTests bool) bool {
	if !IsGoFile(path) {
		return true
	}

	// Skip generated files
	base := filepath.Base(path)
	if strings.HasPrefix(base, "gen_") || strings.HasSuffix(base, "_gen.go") {
		return true
	}

	// Skip test files if configured
	if excludeTests && IsTestFile(path) {
		return true
	}

	// Check for build tags that might indicate non-standard files
	// This is a simple heuristic - a full implementation would parse the file
	return false
}

// GetPosition gets the position information for a node
func (p *Parser) GetPosition(node ast.Node) token.Position {
	return p.fset.Position(node.Pos())
}

// GetPositionRange gets the start and end position for a node
func (p *Parser) GetPositionRange(node ast.Node) (token.Position, token.Position) {
	return p.fset.Position(node.Pos()), p.fset.Position(node.End())
}
