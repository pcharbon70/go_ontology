package extractor

import (
	"go/ast"
	"go/token"
)

// Location represents source location information
type Location struct {
	File      string
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
}

// Context provides extraction context
type Context struct {
	FileSet    *token.FileSet
	FilePath   string
	Package    string
	PkgImports map[string]string // import path -> alias
	Config     *Config
}

// Config defines extraction configuration
type Config struct {
	IncludeExpressions bool
	IncludeSourceText  bool
	ExtractDocComments bool
}

// ExtractionResult is the interface for all extraction results
type ExtractionResult interface {
	GetType() string
	GetLocation() Location
}

// =============================================================================
// Package Information
// =============================================================================

// PackageInfo represents an extracted Go package
type PackageInfo struct {
	Name       string
	ImportPath string
	Doc        string
	Files      []string
	Imports    []*ImportInfo
	Location   Location
}

func (p *PackageInfo) GetType() string { return "Package" }
func (p *PackageInfo) GetLocation() Location { return p.Location }

// ImportInfo represents an import declaration
type ImportInfo struct {
	Path     string
	Alias    string // "" if no alias
	Location Location
}

// =============================================================================
// Function Information
// =============================================================================

// FunctionInfo represents an extracted function or method
type FunctionInfo struct {
	Name          string
	Receiver      *ReceiverInfo // nil for functions
	Parameters    []*ParameterInfo
	Results       []*ParameterInfo
	Variadic      bool
	Doc           string
	Location      Location
	IsMethod      bool
	IsExported    bool
	Body          ast.Node // nil if not including expressions
}

func (f *FunctionInfo) GetType() string {
	if f.IsMethod {
		return "Method"
	}
	return "Function"
}

func (f *FunctionInfo) GetLocation() Location { return f.Location }

// ReceiverInfo represents a method receiver
type ReceiverInfo struct {
	Name string // "" if anonymous
	Type string
}

// ParameterInfo represents a function parameter or result
type ParameterInfo struct {
	Name string // "" if unnamed
	Type string
}

// =============================================================================
// Type Information
// =============================================================================

// StructInfo represents an extracted struct type
type StructInfo struct {
	Name       string
	TypeParams []*TypeParamInfo // Generic type parameters (Go 1.18+)
	Fields     []*FieldInfo
	Doc        string
	Location   Location
	IsExported bool
}

func (s *StructInfo) GetType() string { return "Struct" }
func (s *StructInfo) GetLocation() Location { return s.Location }

// InterfaceInfo represents an extracted interface type
type InterfaceInfo struct {
	Name          string
	TypeParams    []*TypeParamInfo // Generic type parameters (Go 1.18+)
	Methods       []*FunctionInfo   // Interface methods (signatures only)
	EmbeddedTypes []string          // Embedded interface names
	Doc           string
	Location      Location
	IsExported    bool
}

func (i *InterfaceInfo) GetType() string { return "Interface" }
func (i *InterfaceInfo) GetLocation() Location { return i.Location }

// AliasInfo represents a type alias
type AliasInfo struct {
	Name          string
	TypeParams    []*TypeParamInfo
	UnderlyingType string
	Doc           string
	Location      Location
	IsExported    bool
}

func (a *AliasInfo) GetType() string { return "TypeAlias" }
func (a *AliasInfo) GetLocation() Location { return a.Location }

// FieldInfo represents a struct field
type FieldInfo struct {
	Name      string // "" if anonymous/embedded
	Type      string
	Tag       string
	Doc       string
	IsExported bool
	IsEmbedded bool
	Location   Location
}

// TypeParamInfo represents a generic type parameter (Go 1.18+)
type TypeParamInfo struct {
	Name string
	// Constraint is the interface type that constrains the type parameter
	Constraint string
}

// =============================================================================
// Concurrency Information
// =============================================================================

// ConcurrencyInfo represents concurrency-related constructs
type ConcurrencyInfo struct {
	Type     string // "goroutine", "channel", "select", "defer", "panic", "recover"
	Details   map[string]interface{}
	Location  Location
}

func (c *ConcurrencyInfo) GetType() string { return "Concurrency:" + c.Type }
func (c *ConcurrencyInfo) GetLocation() Location { return c.Location }

// =============================================================================
// Import Alias Information
// =============================================================================

// ImportAliasInfo represents an import with an alias
type ImportAliasInfo struct {
	Path     string
	Alias    string
	Location Location
}

func (i *ImportAliasInfo) GetType() string { return "ImportAlias" }
func (i *ImportAliasInfo) GetLocation() Location { return i.Location }
