package config

import "time"

// Config holds configuration for Go ontology analysis
type Config struct {
	// BaseIRI is the base IRI for generated resources
	BaseIRI string

	// IncludeExpressions determines whether to include full AST expressions
	IncludeExpressions bool

	// IncludeSourceText determines whether to include source code text
	IncludeSourceText bool

	// IncludeDocs determines whether to include documentation comments
	IncludeDocs bool

	// ExcludeTests determines whether to exclude test files
	ExcludeTests bool

	// RecursionDepth limits directory recursion (0 = unlimited)
	RecursionDepth int

	// Timeout for parsing operations
	Timeout time.Duration
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		BaseIRI:           "https://w3id.org/go-code/",
		IncludeExpressions: false,
		IncludeSourceText:  false,
		IncludeDocs:       true,
		ExcludeTests:      true,
		RecursionDepth:    0,
		Timeout:           30 * time.Second,
	}
}

// Option is a function that modifies a Config
type Option func(*Config)

// WithBaseIRI sets the base IRI
func WithBaseIRI(iri string) Option {
	return func(c *Config) {
		c.BaseIRI = iri
	}
}

// WithIncludeExpressions sets whether to include expressions
func WithIncludeExpressions(include bool) Option {
	return func(c *Config) {
		c.IncludeExpressions = include
	}
}

// WithIncludeSourceText sets whether to include source text
func WithIncludeSourceText(include bool) Option {
	return func(c *Config) {
		c.IncludeSourceText = include
	}
}

// WithIncludeDocs sets whether to include documentation
func WithIncludeDocs(include bool) Option {
	return func(c *Config) {
		c.IncludeDocs = include
	}
}

// WithExcludeTests sets whether to exclude test files
func WithExcludeTests(exclude bool) Option {
	return func(c *Config) {
		c.ExcludeTests = exclude
	}
}

// WithRecursionDepth sets the maximum recursion depth
func WithRecursionDepth(depth int) Option {
	return func(c *Config) {
		c.RecursionDepth = depth
	}
}

// WithTimeout sets the parsing timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// New creates a new config with the given options
func New(opts ...Option) *Config {
	c := Default()
	for _, opt := range opts {
		opt(c)
	}
	return c
}
