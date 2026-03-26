package extractor

import (
	"go/ast"
	"strings"
)

// PackageExtractor extracts package information
type PackageExtractor struct{}

// NewPackageExtractor creates a new package extractor
func NewPackageExtractor() *PackageExtractor {
	return &PackageExtractor{}
}

// Extract extracts package info from an AST node
func (e *PackageExtractor) Extract(node ast.Node, ctx *Context) ([]ExtractionResult, error) {
	file, ok := node.(*ast.File)
	if !ok {
		return nil, nil
	}

	// Get package doc
	var doc string
	if ctx.Config.ExtractDocComments && file.Doc != nil {
		doc = file.Doc.Text()
	}

	// Collect imports
	imports := make([]*ImportInfo, 0, len(file.Imports))
	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		alias := ""
		if imp.Name != nil {
			alias = imp.Name.Name
		}

		pos := ctx.FileSet.Position(imp.Pos())
		imports = append(imports, &ImportInfo{
			Path:     importPath,
			Alias:    alias,
			Location: Location{
				File:      ctx.FilePath,
				StartLine: pos.Line,
				StartCol:  pos.Column,
			},
		})
	}

	pkgInfo := &PackageInfo{
		Name:       file.Name.Name,
		ImportPath: determineImportPath(ctx.FilePath),
		Doc:        doc,
		Files:      []string{ctx.FilePath},
		Imports:    imports,
		Location: Location{
			File:      ctx.FilePath,
			StartLine: 1,
			StartCol:  1,
		},
	}

	return []ExtractionResult{pkgInfo}, nil
}

// determineImportPath determines the import path from a file path
// This is a simplified version - a full implementation would use go/build
func determineImportPath(filePath string) string {
	// Convert file path to potential import path
	// For now, return a placeholder based on directory name
	parts := strings.Split(filePath, "/")
	for i, part := range parts {
		if part == "src" && i < len(parts)-1 {
			// Found src directory, return the rest as import path
			return strings.Join(parts[i+1:], "/")
		}
	}

	// Fallback: use "main" for main packages or the directory name
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		return strings.TrimSuffix(lastPart, ".go")
	}

	return "main"
}
