package analyzer

import (
	"go/ast"
	"os"
	"path/filepath"
	"strings"

	"github.com/pascal/ontology/pkg/config"
	"github.com/pascal/ontology/pkg/extractor"
)

// FileAnalyzer analyzes a single Go source file
type FileAnalyzer struct {
	parser    *Parser
	extractors []Extractor
	cfg       *config.Config
}

// Extractor is the interface for AST node extractors
type Extractor interface {
	Extract(node ast.Node, ctx *extractor.Context) ([]extractor.ExtractionResult, error)
}

// NewFileAnalyzer creates a new file analyzer
func NewFileAnalyzer(cfg *config.Config) *FileAnalyzer {
	return &FileAnalyzer{
		parser:    NewParser(),
		extractors: make([]Extractor, 0),
		cfg:       cfg,
	}
}

// RegisterExtractor registers an extractor with the analyzer
func (fa *FileAnalyzer) RegisterExtractor(extractor Extractor) {
	fa.extractors = append(fa.extractors, extractor)
}

// AnalysisResult represents the result of analyzing a file
type AnalysisResult struct {
	FilePath string
	Package  *extractor.PackageInfo
	Entities []extractor.ExtractionResult
	Graph    interface{} // Will be *rdf.Graph once we have it
	Errors   []error
}

// Analyze analyzes a single Go file
func (fa *FileAnalyzer) Analyze(filePath string) (*AnalysisResult, error) {
	// Parse the file
	file, err := fa.parser.ParseFile(filePath)
	if err != nil {
		return nil, err
	}

	// Create extraction context
	ctx := &extractor.Context{
		FileSet:    fa.parser.FileSet(),
		FilePath:   filePath,
		Package:    file.Name.Name,
		PkgImports: make(map[string]string),
		Config: &extractor.Config{
			IncludeExpressions: fa.cfg.IncludeExpressions,
			IncludeSourceText:  fa.cfg.IncludeSourceText,
			ExtractDocComments: fa.cfg.IncludeDocs,
		},
	}

	// Collect imports
	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		alias := ""
		if imp.Name != nil {
			alias = imp.Name.Name
		}
		ctx.PkgImports[importPath] = alias
	}

	result := &AnalysisResult{
		FilePath: filePath,
		Entities: make([]extractor.ExtractionResult, 0),
		Errors:   make([]error, 0),
	}

	// Extract package info
	pkgInfo := fa.extractPackageInfo(file, ctx)
	result.Package = pkgInfo
	result.Entities = append(result.Entities, pkgInfo)

	// Walk the AST and extract entities
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		// Skip the file node itself (already handled)
		if _, ok := n.(*ast.File); ok {
			return true
		}

		// Try each extractor
		for _, ext := range fa.extractors {
			results, err := ext.Extract(n, ctx)
			if err != nil {
				result.Errors = append(result.Errors, err)
			}
			result.Entities = append(result.Entities, results...)
		}

		return true
	})

	return result, nil
}

// extractPackageInfo extracts package information from a file
func (fa *FileAnalyzer) extractPackageInfo(file *ast.File, ctx *extractor.Context) *extractor.PackageInfo {
	// Get package doc (first comment group)
	var doc string
	if fa.cfg.IncludeDocs && file.Doc != nil {
		doc = file.Doc.Text()
	}

	// Get import path from file path
	// This is a simplified version - a full implementation would use go/build
	importPath := fa.determineImportPath(ctx.FilePath)

	// Collect imports
	imports := make([]*extractor.ImportInfo, 0, len(file.Imports))
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		alias := ""
		if imp.Name != nil {
			alias = imp.Name.Name
		}

		pos := fa.parser.GetPosition(imp)
		imports = append(imports, &extractor.ImportInfo{
			Path:     path,
			Alias:    alias,
			Location: extractor.Location{
				File:      ctx.FilePath,
				StartLine: pos.Line,
				StartCol:  pos.Column,
			},
		})
	}

	return &extractor.PackageInfo{
		Name:       file.Name.Name,
		ImportPath: importPath,
		Doc:        doc,
		Files:      []string{ctx.FilePath},
		Imports:    imports,
		Location: extractor.Location{
			File:      ctx.FilePath,
			StartLine: 1,
			StartCol:  1,
		},
	}
}

// determineImportPath determines the import path from a file path
// This is a simplified version - a full implementation would use go/build
func (fa *FileAnalyzer) determineImportPath(filePath string) string {
	// Get the module path from go.mod if available
	modulePath := getModulePath()

	// Convert absolute file path to relative path
	// First, try to find the project root (where go.mod is)
	projectRoot := findProjectRoot(filePath)
	if projectRoot == "" {
		// No go.mod found, use directory name as package
		return filepath.Base(filepath.Dir(filePath))
	}

	// Get the relative path from project root
	relPath, err := filepath.Rel(projectRoot, filepath.Dir(filePath))
	if err != nil {
		return filepath.Base(filepath.Dir(filePath))
	}

	// If we're at the root, use module path
	if relPath == "." {
		return modulePath
	}

	// Combine module path with relative path
	return modulePath + "/" + filepath.ToSlash(relPath)
}

// getModulePath reads the module path from go.mod
func getModulePath() string {
	// Try to read go.mod from current directory
	// For now, return a default based on the project
	return "github.com/pascal/ontology"
}

// findProjectRoot finds the project root by looking for go.mod
func findProjectRoot(filePath string) string {
	dir := filepath.Dir(filePath)
	for {
		if dir == "/" || dir == "." {
			break
		}
		goMod := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goMod); err == nil {
			// File exists
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}
