package builder

import (
	"github.com/pascal/ontology/pkg/extractor"
	"github.com/pascal/ontology/pkg/rdf"
)

// Builder is the interface for RDF builders
type Builder interface {
	// Build transforms extraction results into RDF triples
	Build(result extractor.ExtractionResult, ctx *Context) (*rdf.Graph, error)

	// CanBuild checks if this builder can handle the given result
	CanBuild(result extractor.ExtractionResult) bool
}

// Context provides builder context
type Context struct {
	BaseIRI           string
	FilePath          string
	CurrentPackage    string
	PackageImportPath string
}

// NewContext creates a new builder context
func NewContext(baseIRI, filePath, currentPackage, importPath string) *Context {
	return &Context{
		BaseIRI:           baseIRI,
		FilePath:          filePath,
		CurrentPackage:    currentPackage,
		PackageImportPath: importPath,
	}
}
