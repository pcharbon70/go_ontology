package builder

import (
	"fmt"

	"github.com/pascal/ontology/pkg/extractor"
	"github.com/pascal/ontology/pkg/rdf"
)

// InterfaceBuilder builds RDF triples for interface types
type InterfaceBuilder struct {
	helper *Helper
}

// NewInterfaceBuilder creates a new interface builder
func NewInterfaceBuilder() *InterfaceBuilder {
	return &InterfaceBuilder{
		helper: NewHelper(),
	}
}

// Build transforms an InterfaceInfo into RDF triples
func (b *InterfaceBuilder) Build(result extractor.ExtractionResult, ctx *Context) (*rdf.Graph, error) {
	ifaceInfo, ok := result.(*extractor.InterfaceInfo)
	if !ok {
		return nil, fmt.Errorf("expected InterfaceInfo, got %T", result)
	}

	graph := rdf.NewGraph()

	// Create IRI for the interface
	ifaceIRI := b.createInterfaceIRI(ctx, ifaceInfo.Name)

	// Add type triple
	graph.Add(b.helper.TypeTriple(ifaceIRI, rdf.GoInterface))

	// Add basic properties
	graph.Add(b.helper.StringProperty(ifaceIRI, rdf.StructName, ifaceInfo.Name))

	if ifaceInfo.Doc != "" {
		graph.Add(b.helper.StringProperty(ifaceIRI, rdf.RDFSComment, ifaceInfo.Doc))
	}

	// Add package membership
	pkgIRI := rdf.NewIRI(ctx.BaseIRI + "#" + ctx.PackageImportPath)
	graph.Add(rdf.NewTriple(ifaceIRI, rdf.BelongsTo, pkgIRI))

	// Add type parameters if present
	for _, tp := range ifaceInfo.TypeParams {
		typeParamIRI := rdf.NewIRI(ifaceIRI.String() + "/" + tp.Name)
		graph.Add(b.helper.TypeTriple(typeParamIRI, rdf.GoTypeParameter))
		graph.Add(b.helper.StringProperty(typeParamIRI, rdf.TypeParamName, tp.Name))
		graph.Add(b.helper.StringProperty(typeParamIRI, rdf.TypeParamConstraint, tp.Constraint))
		graph.Add(rdf.NewTriple(typeParamIRI, rdf.BelongsTo, ifaceIRI))
	}

	// Add methods
	for _, method := range ifaceInfo.Methods {
		methodIRI := rdf.NewIRI(fmt.Sprintf("%s/%s", ifaceIRI.String(), method.Name))
		graph.Add(b.helper.TypeTriple(methodIRI, rdf.GoInterfaceMethod))
		graph.Add(b.helper.StringProperty(methodIRI, rdf.FunctionName, method.Name))

		// Add parameter and result counts
		graph.Add(b.helper.IntegerProperty(methodIRI, rdf.ParameterCount, len(method.Parameters)))
		graph.Add(b.helper.IntegerProperty(methodIRI, rdf.ResultCount, len(method.Results)))

		if method.Doc != "" {
			graph.Add(b.helper.StringProperty(methodIRI, rdf.RDFSComment, method.Doc))
		}

		graph.Add(rdf.NewTriple(methodIRI, rdf.BelongsTo, ifaceIRI))
	}

	// Add embedded types
	for _, embedded := range ifaceInfo.EmbeddedTypes {
		// Create a triple to indicate embedding
		embeddedIRI := rdf.NewIRI(ctx.BaseIRI + "#" + ctx.CurrentPackage + "." + embedded)
		graph.Add(rdf.NewTriple(ifaceIRI, rdf.NewIRI(rdf.GoStructure+"embeds"), embeddedIRI))
	}

	return graph, nil
}

// CanBuild checks if this builder can handle the given result
func (b *InterfaceBuilder) CanBuild(result extractor.ExtractionResult) bool {
	_, ok := result.(*extractor.InterfaceInfo)
	return ok
}

// createInterfaceIRI creates an IRI for an interface
func (b *InterfaceBuilder) createInterfaceIRI(ctx *Context, ifaceName string) *rdf.IRI {
	return rdf.NewIRI(ctx.BaseIRI + "#" + ctx.PackageImportPath + "." + ifaceName)
}
