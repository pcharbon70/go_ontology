package extractor

import (
	"go/ast"
	"go/token"
)

// StructExtractor extracts struct type declarations
type StructExtractor struct{}

// NewStructExtractor creates a new struct extractor
func NewStructExtractor() *StructExtractor {
	return &StructExtractor{}
}

// Extract extracts struct info from an AST node
func (e *StructExtractor) Extract(node ast.Node, ctx *Context) ([]ExtractionResult, error) {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.TYPE {
		return nil, nil
	}

	var results []ExtractionResult

	for _, spec := range decl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		structInfo := &StructInfo{
			Name:       typeSpec.Name.Name,
			TypeParams: e.extractTypeParams(typeSpec),
			Fields:     e.extractFields(structType),
			Doc:        getDeclDoc(decl, typeSpec),
			Location:   getLocation(typeSpec, ctx),
			IsExported: typeSpec.Name.IsExported(),
		}

		results = append(results, structInfo)
	}

	return results, nil
}

// extractTypeParams extracts generic type parameters
func (e *StructExtractor) extractTypeParams(spec *ast.TypeSpec) []*TypeParamInfo {
	if spec.TypeParams == nil {
		return nil
	}

	var params []*TypeParamInfo
	for _, param := range spec.TypeParams.List {
		constraint := exprToString(param.Type)
		for _, name := range param.Names {
			params = append(params, &TypeParamInfo{
				Name:      name.Name,
				Constraint: constraint,
			})
		}
	}
	return params
}

// extractFields extracts struct fields
func (e *StructExtractor) extractFields(structType *ast.StructType) []*FieldInfo {
	if structType.Fields == nil {
		return nil
	}

	var fields []*FieldInfo
	for _, field := range structType.Fields.List {
		typeName := exprToString(field.Type)

		// Get field doc from comment
		var doc string
		if field.Doc != nil {
			doc = field.Doc.Text()
		} else if field.Comment != nil {
			doc = field.Comment.Text()
		}

		// Get tag
		var tag string
		if field.Tag != nil {
			tag = field.Tag.Value
		}

		if len(field.Names) == 0 {
			// Anonymous/embedded field
			fields = append(fields, &FieldInfo{
				Name:       "",
				Type:       typeName,
				Tag:        tag,
				Doc:        doc,
				IsExported: isExportedType(typeName),
				IsEmbedded: true,
			})
		} else {
			// Named fields
			for _, name := range field.Names {
				fields = append(fields, &FieldInfo{
					Name:       name.Name,
					Type:       typeName,
					Tag:        tag,
					Doc:        doc,
					IsExported: name.IsExported(),
					IsEmbedded: false,
				})
			}
		}
	}
	return fields
}

// getDeclDoc gets documentation for a declaration
func getDeclDoc(decl *ast.GenDecl, spec *ast.TypeSpec) string {
	if decl.Doc != nil {
		return decl.Doc.Text()
	}
	if spec.Doc != nil {
		return spec.Doc.Text()
	}
	return ""
}

// isExportedType checks if a type name is exported
func isExportedType(typeName string) bool {
	if len(typeName) == 0 {
		return false
	}
	// First character is uppercase
	return typeName[0] >= 'A' && typeName[0] <= 'Z'
}
