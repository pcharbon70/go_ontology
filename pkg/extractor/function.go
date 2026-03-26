package extractor

import (
	"fmt"
	"go/ast"
	"strings"
)

// FunctionExtractor extracts function and method declarations
type FunctionExtractor struct{}

// NewFunctionExtractor creates a new function extractor
func NewFunctionExtractor() *FunctionExtractor {
	return &FunctionExtractor{}
}

// Extract extracts function/method info from an AST node
func (e *FunctionExtractor) Extract(node ast.Node, ctx *Context) ([]ExtractionResult, error) {
	decl, ok := node.(*ast.FuncDecl)
	if !ok {
		return nil, nil
	}

		funcInfo := &FunctionInfo{
		Name:      decl.Name.Name,
		Doc:       getDocComment(decl),
		Location:  getLocation(decl, ctx),
		IsExported: decl.Name.IsExported(),
	}

	// Extract receiver if present (method)
	if decl.Recv != nil && len(decl.Recv.List) > 0 {
		funcInfo.IsMethod = true
		funcInfo.Receiver = e.extractReceiver(decl.Recv.List[0])
	}

	// Extract parameters
	if decl.Type.Params != nil {
		funcInfo.Parameters = e.extractParams(decl.Type.Params, ctx)
	}

	// Extract results
	if decl.Type.Results != nil {
		funcInfo.Results = e.extractParams(decl.Type.Results, ctx)
	}

	// Check for variadic
	if len(funcInfo.Parameters) > 0 {
		lastField := decl.Type.Params.List[len(decl.Type.Params.List)-1]
		if lastField != nil {
			// Check if the type is an ellipsis (variadic)
			if _, ok := lastField.Type.(*ast.Ellipsis); ok {
				funcInfo.Variadic = true
			}
		}
	}

	// Include body if configured
	if ctx.Config.IncludeExpressions && decl.Body != nil {
		funcInfo.Body = decl.Body
	}

	return []ExtractionResult{funcInfo}, nil
}

// extractReceiver extracts receiver information
func (e *FunctionExtractor) extractReceiver(field *ast.Field) *ReceiverInfo {
	typeName := exprToString(field.Type)

	var recvName string
	if len(field.Names) > 0 {
		recvName = field.Names[0].Name
	}

	return &ReceiverInfo{
		Name: recvName,
		Type: typeName,
	}
}

// extractParams extracts parameters or return values
func (e *FunctionExtractor) extractParams(params *ast.FieldList, ctx *Context) []*ParameterInfo {
	if params == nil || len(params.List) == 0 {
		return []*ParameterInfo{}
	}

	var results []*ParameterInfo

	for _, field := range params.List {
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

// getDocComment extracts documentation comment
func getDocComment(node ast.Node) string {
	if doc := node.(*ast.FuncDecl).Doc; doc != nil {
		return doc.Text()
	}
	return ""
}

// getLocation extracts location from a node
func getLocation(node ast.Node, ctx *Context) Location {
	pos := ctx.FileSet.Position(node.Pos())
	end := ctx.FileSet.Position(node.End())

	return Location{
		File:      ctx.FilePath,
		StartLine: pos.Line,
		StartCol:  pos.Column,
		EndLine:   end.Line,
		EndCol:    end.Column,
	}
}

// exprToString converts an AST expression to a string (simplified)
func exprToString(expr ast.Expr) string {
	if expr == nil {
		return ""
	}

	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", exprToString(v.X), v.Sel.Name)
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", exprToString(v.X))
	case *ast.ArrayType:
		return fmt.Sprintf("[]%s", exprToString(v.Elt))
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", exprToString(v.Key), exprToString(v.Value))
	case *ast.ChanType:
		if v.Dir == ast.SEND {
			return fmt.Sprintf("chan<- %s", exprToString(v.Value))
		} else if v.Dir == ast.RECV {
			return fmt.Sprintf("<-chan %s", exprToString(v.Value))
		}
		return fmt.Sprintf("chan %s", exprToString(v.Value))
	case *ast.FuncType:
		// Build function signature string
		return "func"
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.Ellipsis:
		return fmt.Sprintf("...%s", exprToString(v.Elt))
	case *ast.IndexExpr:
		// Generic type like T[K]
		return fmt.Sprintf("%s[%s]", exprToString(v.X), exprToString(v.Index))
	case *ast.IndexListExpr:
		// Generic type with multiple params like T[K, V]
		params := make([]string, 0)
		for _, indices := range v.Indices {
			params = append(params, exprToString(indices))
		}
		return fmt.Sprintf("%s[%s]", exprToString(v.X), strings.Join(params, ", "))
	default:
		return ""
	}
}
