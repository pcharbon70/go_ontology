package builder

import (
	"fmt"

	"github.com/pascal/ontology/pkg/extractor"
	"github.com/pascal/ontology/pkg/rdf"
)

// StructBuilder builds RDF triples for struct types
type StructBuilder struct {
	helper *Helper
}

// NewStructBuilder creates a new struct builder
func NewStructBuilder() *StructBuilder {
	return &StructBuilder{
		helper: NewHelper(),
	}
}

// Build transforms a StructInfo into RDF triples
func (b *StructBuilder) Build(result extractor.ExtractionResult, ctx *Context) (*rdf.Graph, error) {
	structInfo, ok := result.(*extractor.StructInfo)
	if !ok {
		return nil, fmt.Errorf("expected StructInfo, got %T", result)
	}

	graph := rdf.NewGraph()

	// Create IRI for the struct
	structIRI := b.createStructIRI(ctx, structInfo.Name)

	// Add type triple
	graph.Add(b.helper.TypeTriple(structIRI, rdf.GoStruct))

	// Add basic properties
	graph.Add(b.helper.StringProperty(structIRI, rdf.StructName, structInfo.Name))

	if structInfo.Doc != "" {
		graph.Add(b.helper.StringProperty(structIRI, rdf.RDFSComment, structInfo.Doc))
	}

	// Add package membership
	pkgIRI := rdf.NewIRI(ctx.BaseIRI + "#" + ctx.PackageImportPath)
	graph.Add(rdf.NewTriple(structIRI, rdf.BelongsTo, pkgIRI))

	// Add type parameters if present
	for _, tp := range structInfo.TypeParams {
		typeParamIRI := rdf.NewIRI(structIRI.String() + "/" + tp.Name)
		graph.Add(b.helper.TypeTriple(typeParamIRI, rdf.GoTypeParameter))
		graph.Add(b.helper.StringProperty(typeParamIRI, rdf.TypeParamName, tp.Name))
		graph.Add(b.helper.StringProperty(typeParamIRI, rdf.TypeParamConstraint, tp.Constraint))
		graph.Add(rdf.NewTriple(typeParamIRI, rdf.BelongsTo, structIRI))
	}

	// Add fields
	for i, field := range structInfo.Fields {
		fieldIRI := rdf.NewIRI(fmt.Sprintf("%s/field/%d", structIRI.String(), i))
		graph.Add(b.helper.TypeTriple(fieldIRI, rdf.GoStructField))
		graph.Add(b.helper.StringProperty(fieldIRI, rdf.FieldName, field.Name))
		graph.Add(b.helper.StringProperty(fieldIRI, rdf.FieldType, field.Type))

		if field.Tag != "" {
			graph.Add(b.helper.StringProperty(fieldIRI, rdf.FieldTag, field.Tag))
		}

		if field.Doc != "" {
			graph.Add(b.helper.StringProperty(fieldIRI, rdf.RDFSComment, field.Doc))
		}

		if field.IsEmbedded {
			graph.Add(b.helper.BooleanProperty(fieldIRI, rdf.IsEmbedded, true))
		}

		if field.IsExported {
			graph.Add(b.helper.BooleanProperty(fieldIRI, rdf.IsExported, true))
		}

		graph.Add(rdf.NewTriple(fieldIRI, rdf.BelongsTo, structIRI))
	}

	return graph, nil
}

// CanBuild checks if this builder can handle the given result
func (b *StructBuilder) CanBuild(result extractor.ExtractionResult) bool {
	_, ok := result.(*extractor.StructInfo)
	return ok
}

// createStructIRI creates an IRI for a struct
func (b *StructBuilder) createStructIRI(ctx *Context, structName string) *rdf.IRI {
	return rdf.NewIRI(ctx.BaseIRI + "#" + ctx.PackageImportPath + "." + structName)
}
