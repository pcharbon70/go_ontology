package extractor

import (
	"go/ast"
	"go/token"
)

// InterfaceTypeExtractor extracts interface type declarations
type InterfaceTypeExtractor struct{}

// NewInterfaceExtractor creates a new interface extractor
func NewInterfaceExtractor() *InterfaceTypeExtractor {
	return &InterfaceTypeExtractor{}
}

// Extract extracts interface info from an AST node
func (e *InterfaceTypeExtractor) Extract(node ast.Node, ctx *Context) ([]ExtractionResult, error) {
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

		ifaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}

		ifaceInfo := &InterfaceInfo{
			Name:          typeSpec.Name.Name,
			TypeParams:    e.extractTypeParams(typeSpec),
			Methods:       e.extractMethods(ifaceType, ctx),
			EmbeddedTypes: e.extractEmbeddedTypes(ifaceType),
			Doc:           getDeclDoc(decl, typeSpec),
			Location:      getLocation(typeSpec, ctx),
			IsExported:    typeSpec.Name.IsExported(),
		}

		results = append(results, ifaceInfo)
	}

	return results, nil
}

// extractTypeParams extracts generic type parameters
func (e *InterfaceTypeExtractor) extractTypeParams(spec *ast.TypeSpec) []*TypeParamInfo {
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

// extractMethods extracts interface methods
func (e *InterfaceTypeExtractor) extractMethods(ifaceType *ast.InterfaceType, ctx *Context) []*FunctionInfo {
	if ifaceType.Methods == nil {
		return nil
	}

	var methods []*FunctionInfo

	for _, meth := range ifaceType.Methods.List {
		switch ft := meth.Type.(type) {
		case *ast.FuncType:
			// Method with signature
			funcInfo := &FunctionInfo{
				Name:      e.getMethodName(meth),
				Doc:       e.getFieldDoc(meth),
				Location:  e.getFieldLocation(meth, ctx),
				IsMethod:  true,
				IsExported: true, // Interface methods are always exported
			}

			// Extract parameters
			if ft.Params != nil {
				funcInfo.Parameters = extractParamsFromFieldList(ft.Params, ctx)
			}

			// Extract results
			if ft.Results != nil {
				funcInfo.Results = extractParamsFromFieldList(ft.Results, ctx)
			}

			methods = append(methods, funcInfo)

		case *ast.Ident, *ast.SelectorExpr:
			// Embedded interface (not a method) - handled in extractEmbeddedTypes
		}
	}

	return methods
}

// extractEmbeddedTypes extracts embedded interface names
func (e *InterfaceTypeExtractor) extractEmbeddedTypes(ifaceType *ast.InterfaceType) []string {
	if ifaceType.Methods == nil {
		return nil
	}

	var embedded []string

	for _, meth := range ifaceType.Methods.List {
		switch meth.Type.(type) {
		case *ast.Ident, *ast.SelectorExpr:
			// This is an embedded interface, not a method
			typeName := exprToString(meth.Type)
			if typeName != "" {
				embedded = append(embedded, typeName)
			}
		}
	}

	return embedded
}

// getMethodName gets the name of a method from a field
func (e *InterfaceTypeExtractor) getMethodName(field *ast.Field) string {
	if len(field.Names) > 0 {
		return field.Names[0].Name
	}
	// Embedded interface - return type name
	return exprToString(field.Type)
}

// getFieldDoc gets documentation from a field
func (e *InterfaceTypeExtractor) getFieldDoc(field *ast.Field) string {
	if field.Doc != nil {
		return field.Doc.Text()
	}
	if field.Comment != nil {
		return field.Comment.Text()
	}
	return ""
}

// getFieldLocation gets location from a field
func (e *InterfaceTypeExtractor) getFieldLocation(field *ast.Field, ctx *Context) Location {
	pos := ctx.FileSet.Position(field.Pos())
	end := ctx.FileSet.Position(field.End())

	return Location{
		File:      ctx.FilePath,
		StartLine: pos.Line,
		StartCol:  pos.Column,
		EndLine:   end.Line,
		EndCol:    end.Column,
	}
}

// extractParamsFromFieldList extracts parameters from a field list
func extractParamsFromFieldList(fl *ast.FieldList, ctx *Context) []*ParameterInfo {
	if fl == nil || len(fl.List) == 0 {
		return nil
	}

	var results []*ParameterInfo

	for _, field := range fl.List {
		typeName := exprToString(field.Type)

		if len(field.Names) == 0 {
			// Unnamed parameter/result
			results = append(results, &ParameterInfo{
				Name: "",
				Type: typeName,
			})
		} else {
			// Named parameters/results
			for _, name := range field.Names {
				results = append(results, &ParameterInfo{
					Name: name.Name,
					Type: typeName,
				})
			}
		}
	}

	return results
}
