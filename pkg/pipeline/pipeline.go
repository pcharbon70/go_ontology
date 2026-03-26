package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pascal/ontology/pkg/analyzer"
	"github.com/pascal/ontology/pkg/builder"
	"github.com/pascal/ontology/pkg/config"
	"github.com/pascal/ontology/pkg/extractor"
	"github.com/pascal/ontology/pkg/rdf"
)

// Pipeline orchestrates the end-to-end analysis and RDF generation
type Pipeline struct {
	cfg         *config.Config
	fileAnalyzer *analyzer.FileAnalyzer
	orchestrator *builder.Orchestrator
}

// New creates a new pipeline with the given configuration
func New(cfg *config.Config) *Pipeline {
	fileAnalyzer := analyzer.NewFileAnalyzer(cfg)

	// Register extractors
	fileAnalyzer.RegisterExtractor(extractor.NewFunctionExtractor())
	fileAnalyzer.RegisterExtractor(extractor.NewStructExtractor())
	fileAnalyzer.RegisterExtractor(extractor.NewInterfaceExtractor())
	fileAnalyzer.RegisterExtractor(extractor.NewConcurrencyExtractor())

	return &Pipeline{
		cfg:         cfg,
		fileAnalyzer: fileAnalyzer,
		orchestrator: builder.NewOrchestrator(),
	}
}

// Result represents the result of pipeline execution
type Result struct {
	Graph    *rdf.Graph
	Turtle   string
	Metadata Metadata
}

// Metadata contains analysis metadata
type Metadata struct {
	FileCount  int
	PackageCount int
	FunctionCount int
	ErrorCount int
}

// AnalyzeFile analyzes a single Go file
func (p *Pipeline) AnalyzeFile(filePath string) (*Result, error) {
	// Analyze the file
	analysisResult, err := p.fileAnalyzer.Analyze(filePath)
	if err != nil {
		return nil, err
	}

	// Create builder context
	ctx := builder.NewContext(
		p.cfg.BaseIRI,
		filePath,
		analysisResult.Package.Name,
		analysisResult.Package.ImportPath,
	)

	// Build RDF graph
	graph, err := p.orchestrator.BuildAll(analysisResult.Entities, ctx)
	if err != nil {
		return nil, err
	}

	// Serialize to Turtle
	turtle, err := rdf.ToTurtle(graph)
	if err != nil {
		return nil, err
	}

	// Count entities
	metadata := Metadata{
		FileCount:  1,
		PackageCount: 1,
		FunctionCount: countEntities(analysisResult.Entities, "Function", "Method"),
		ErrorCount: len(analysisResult.Errors),
	}

	return &Result{
		Graph:    graph,
		Turtle:   turtle,
		Metadata: metadata,
	}, nil
}

// AnalyzeProject analyzes an entire Go project directory
func (p *Pipeline) AnalyzeProject(dirPath string) (*Result, error) {
	// Find all Go files
	goFiles, err := findGoFiles(dirPath, p.cfg.ExcludeTests)
	if err != nil {
		return nil, err
	}

	if len(goFiles) == 0 {
		return nil, fmt.Errorf("no Go files found in %s", dirPath)
	}

	// Collect all entities from all files
	allEntities := make([]extractor.ExtractionResult, 0)
	allPackages := make(map[string]*extractor.PackageInfo)
	totalErrors := 0

	for _, filePath := range goFiles {
		analysisResult, err := p.fileAnalyzer.Analyze(filePath)
		if err != nil {
			totalErrors++
			continue
		}

		// Track unique packages
		pkgKey := analysisResult.Package.ImportPath
		if _, exists := allPackages[pkgKey]; !exists {
			allPackages[pkgKey] = analysisResult.Package
			allEntities = append(allEntities, analysisResult.Package)
		}

		// Add non-package entities
		for _, entity := range analysisResult.Entities {
			if entity.GetType() != "Package" {
				allEntities = append(allEntities, entity)
			}
		}

		totalErrors += len(analysisResult.Errors)
	}

	// Use the first package for context (or main package)
	var primaryPkg *extractor.PackageInfo
	for _, pkg := range allPackages {
		primaryPkg = pkg
		if pkg.Name == "main" || pkg.ImportPath == "main" {
			break
		}
	}

	if primaryPkg == nil && len(allPackages) > 0 {
		for _, pkg := range allPackages {
			primaryPkg = pkg
			break
		}
	}

	// Create builder context
	ctx := builder.NewContext(
		p.cfg.BaseIRI,
		dirPath,
		primaryPkg.Name,
		primaryPkg.ImportPath,
	)

	// Build RDF graph
	graph, err := p.orchestrator.BuildAll(allEntities, ctx)
	if err != nil {
		return nil, err
	}

	// Serialize to Turtle
	turtle, err := rdf.ToTurtle(graph)
	if err != nil {
		return nil, err
	}

	// Count entities
	metadata := Metadata{
		FileCount:    len(goFiles),
		PackageCount: len(allPackages),
		FunctionCount: countEntities(allEntities, "Function", "Method"),
		ErrorCount:   totalErrors,
	}

	return &Result{
		Graph:    graph,
		Turtle:   turtle,
		Metadata: metadata,
	}, nil
}

// findGoFiles finds all Go files in a directory recursively
func findGoFiles(dirPath string, excludeTests bool) ([]string, error) {
	var goFiles []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Skip hidden directories and common exclusions
			if strings.HasPrefix(filepath.Base(path), ".") ||
			   strings.HasPrefix(filepath.Base(path), "_") {
				return filepath.SkipDir
			}

			// Skip vendor directory
			if filepath.Base(path) == "vendor" {
				return filepath.SkipDir
			}

			return nil
		}

		// Check if file should be excluded
		if analyzer.ShouldExclude(path, excludeTests) {
			return nil
		}

		goFiles = append(goFiles, path)
		return nil
	})

	return goFiles, err
}

// countEntities counts entities of specific types
func countEntities(entities []extractor.ExtractionResult, types ...string) int {
	count := 0
	for _, e := range entities {
		type_ := e.GetType()
		for _, t := range types {
			if type_ == t {
				count++
				break
			}
		}
	}
	return count
}

// WriteToFile writes the Turtle output to a file
func (r *Result) WriteToFile(filePath string) error {
	return os.WriteFile(filePath, []byte(r.Turtle), 0644)
}
