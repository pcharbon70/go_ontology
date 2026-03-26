package builder

import (
	"fmt"

	"github.com/pascal/ontology/pkg/extractor"
	"github.com/pascal/ontology/pkg/rdf"
)

const structNS = "https://w3id.org/go-code/structure#"

// PackageBuilder builds RDF triples for Go packages
type PackageBuilder struct {
	helper *Helper
}

// NewPackageBuilder creates a new package builder
func NewPackageBuilder() *PackageBuilder {
	return &PackageBuilder{
		helper: NewHelper(),
	}
}

func (b *PackageBuilder) CanBuild(result extractor.ExtractionResult) bool {
	_, ok := result.(*extractor.PackageInfo)
	return ok
}

func (b *PackageBuilder) Build(
	result extractor.ExtractionResult,
	ctx *Context,
) (*rdf.Graph, error) {
	pkg, ok := result.(*extractor.PackageInfo)
	if !ok {
		return nil, fmt.Errorf("expected PackageInfo, got %T", result)
	}

	graph := rdf.NewGraph()

	// Generate package IRI
	pkgIRI := b.packageIRI(ctx.BaseIRI, pkg.ImportPath)

	// Add type triple
	graph.Add(b.helper.TypeTriple(pkgIRI, rdf.NewIRI(structNS+"Package")))

	// Add name
	graph.Add(b.helper.StringProperty(
		pkgIRI,
		rdf.NewIRI(structNS+"packageName"),
		pkg.Name,
	))

	// Add import path
	graph.Add(b.helper.StringProperty(
		pkgIRI,
		rdf.NewIRI(structNS+"importPath"),
		pkg.ImportPath,
	))

	// Add doc if present
	if pkg.Doc != "" {
		graph.Add(b.helper.StringProperty(
			pkgIRI,
			rdf.NewIRI(structNS+"hasDoc"),
			pkg.Doc,
		))
	}

	// Add imports
	for _, imp := range pkg.Imports {
		importIRI := b.packageIRI(ctx.BaseIRI, imp.Path)
		graph.Add(b.helper.ObjectProperty(
			pkgIRI,
			rdf.NewIRI(structNS+"importsPackage"),
			importIRI,
		))
	}

	return graph, nil
}

func (b *PackageBuilder) packageIRI(base, importPath string) *rdf.IRI {
	return rdf.NewIRI(base + "#" + importPath)
}
