package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/pascal/ontology/pkg/config"
	"github.com/pascal/ontology/pkg/pipeline"
)

var (
	output       string
	baseIRI      string
	includeExpr  bool
	includeDoc   bool
	excludeTests bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "go-ontology",
		Short: "Transform Go source code into RDF knowledge graphs",
		Long: `go-ontology analyzes Go source code and generates RDF knowledge
graphs representing code structure, types, functions, and concurrency
patterns following the Go Ontology specification.

Example:
  go-ontology analyze file.go -o output.ttl
  go-ontology analyze ./myproject -o project.ttl`,
	}

	analyzeCmd := &cobra.Command{
		Use:   "analyze [path]",
		Short: "Analyze Go code and generate RDF",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runAnalyze,
	}

	analyzeCmd.Flags().StringVarP(&output, "output", "o", "", "Output file (default: stdout)")
	analyzeCmd.Flags().StringVarP(&baseIRI, "base-iri", "b", "https://w3id.org/go-code/", "Base IRI for resources")
	analyzeCmd.Flags().BoolVar(&includeExpr, "include-expressions", false, "Include full AST expression triples")
	analyzeCmd.Flags().BoolVar(&includeDoc, "include-docs", true, "Include documentation comments")
	analyzeCmd.Flags().BoolVar(&excludeTests, "exclude-tests", true, "Exclude test files")

	rootCmd.AddCommand(analyzeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("error resolving path: %w", err)
	}

	// Check if path is a file or directory
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("error accessing path: %w", err)
	}

	// Create configuration
	cfg := config.New(
		config.WithBaseIRI(baseIRI),
		config.WithIncludeExpressions(includeExpr),
		config.WithIncludeDocs(includeDoc),
		config.WithExcludeTests(excludeTests),
	)

	// Create pipeline
	p := pipeline.New(cfg)

	var result *pipeline.Result

	if info.IsDir() {
		// Analyze project
		result, err = p.AnalyzeProject(absPath)
		if err != nil {
			return fmt.Errorf("project analysis failed: %w", err)
		}
	} else {
		// Analyze single file
		result, err = p.AnalyzeFile(absPath)
		if err != nil {
			return fmt.Errorf("file analysis failed: %w", err)
		}
	}

	// Output results
	if output == "" {
		fmt.Print(result.Turtle)
	} else {
		if err := result.WriteToFile(output); err != nil {
			return fmt.Errorf("write failed: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Output written to %s\n", output)
	}

	// Print summary to stderr
	fmt.Fprintf(os.Stderr, "Analyzed %d file(s), %d package(s), %d function(s)\n",
		result.Metadata.FileCount,
		result.Metadata.PackageCount,
		result.Metadata.FunctionCount,
	)

	if result.Metadata.ErrorCount > 0 {
		fmt.Fprintf(os.Stderr, "Warning: %d error(s) encountered during analysis\n",
			result.Metadata.ErrorCount)
	}

	return nil
}
