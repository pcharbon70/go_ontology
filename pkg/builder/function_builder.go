package builder

import (
	"fmt"
	"strings"

	"github.com/pascal/ontology/pkg/extractor"
	"github.com/pascal/ontology/pkg/rdf"
)

// FunctionBuilder builds RDF triples for Go functions and methods
type FunctionBuilder struct {
	helper *Helper
}

// NewFunctionBuilder creates a new function builder
func NewFunctionBuilder() *FunctionBuilder {
	return &FunctionBuilder{
		helper: NewHelper(),
	}
}

func (b *FunctionBuilder) CanBuild(result extractor.ExtractionResult) bool {
	_, ok := result.(*extractor.FunctionInfo)
	return ok
}

func (b *FunctionBuilder) Build(
	result extractor.ExtractionResult,
	ctx *Context,
) (*rdf.Graph, error) {
	fn, ok := result.(*extractor.FunctionInfo)
	if !ok {
		return nil, fmt.Errorf("expected FunctionInfo, got %T", result)
	}

	graph := rdf.NewGraph()

	// Generate function IRI
	fnIRI := b.functionIRI(ctx.BaseIRI, ctx.PackageImportPath, fn)

	// Determine type
	var classType string
	if fn.IsMethod {
		if fn.Receiver != nil && isPointerType(fn.Receiver.Type) {
			classType = "PointerMethod"
		} else {
			classType = "ValueMethod"
		}
	} else if fn.IsExported {
		classType = "ExportedFunction"
	} else {
		classType = "Function"
	}

	// Add type triple
	graph.Add(b.helper.TypeTriple(
		fnIRI,
		rdf.NewIRI(structNS+classType),
	))

	// Add name
	graph.Add(b.helper.StringProperty(
		fnIRI,
		rdf.NewIRI(structNS+"functionName"),
		fn.Name,
	))

	// Add parameter count
	graph.Add(b.helper.IntegerProperty(
		fnIRI,
		rdf.NewIRI(structNS+"parameterCount"),
		len(fn.Parameters),
	))

	// Add result count
	graph.Add(b.helper.IntegerProperty(
		fnIRI,
		rdf.NewIRI(structNS+"resultCount"),
		len(fn.Results),
	))

	// Add variadic flag
	if fn.Variadic {
		graph.Add(b.helper.BooleanProperty(
			fnIRI,
			rdf.NewIRI(structNS+"isVariadic"),
			true,
		))
	}

	// Add receiver info for methods
	if fn.IsMethod && fn.Receiver != nil {
		// Find the type IRI for the receiver
		typeIRI := b.typeIRI(ctx.BaseIRI, fn.Receiver.Type)
		graph.Add(b.helper.ObjectProperty(
			fnIRI,
			rdf.NewIRI(structNS+"hasReceiver"),
			typeIRI,
		))

		// Add receiver name if present
		if fn.Receiver.Name != "" {
			graph.Add(b.helper.StringProperty(
				fnIRI,
				rdf.NewIRI(structNS+"receiverName"),
				fn.Receiver.Name,
			))
		}

		// Add receiver type
		graph.Add(b.helper.StringProperty(
			fnIRI,
			rdf.NewIRI(structNS+"receiverType"),
			fn.Receiver.Type,
		))
	}

	// Add doc if present
	if fn.Doc != "" {
		graph.Add(b.helper.StringProperty(
			fnIRI,
			rdf.NewIRI(structNS+"hasDoc"),
			fn.Doc,
		))
	}

	// Add belongsTo relationship to package
	pkgIRI := b.packageIRI(ctx.BaseIRI, ctx.PackageImportPath)
	graph.Add(b.helper.ObjectProperty(
		fnIRI,
		rdf.NewIRI(structNS+"belongsTo"),
		pkgIRI,
	))

	return graph, nil
}

// functionIRI generates a unique IRI for a function
// Format: base#package.FunctionName
// For methods: base#package.ReceiverType.FunctionName
func (b *FunctionBuilder) functionIRI(base, pkg string, fn *extractor.FunctionInfo) *rdf.IRI {
	// pkg should be the full import path (PackageImportPath)
	var sb strings.Builder
	sb.WriteString(base)
	sb.WriteString("#")
	sb.WriteString(pkg)
	sb.WriteString(".")

	if fn.IsMethod && fn.Receiver != nil {
		// Include receiver type in method IRI
		sb.WriteString(fn.Receiver.Type)
		sb.WriteString(".")
	}

	sb.WriteString(fn.Name)

	return rdf.NewIRI(sb.String())
}

// typeIRI generates an IRI for a type
func (b *FunctionBuilder) typeIRI(base, typeName string) *rdf.IRI {
	// For built-in types and basic types, use simple IRIs
	// For user-defined types, would need to resolve the package
	return rdf.NewIRI(base + "#" + typeName)
}

// packageIRI generates an IRI for a package
func (b *FunctionBuilder) packageIRI(base, importPath string) *rdf.IRI {
	return rdf.NewIRI(base + "#" + importPath)
}

// isPointerType checks if a type string represents a pointer type
func isPointerType(typeStr string) bool {
	return strings.HasPrefix(typeStr, "*")
}
