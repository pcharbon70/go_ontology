# Go Ontology

An OWL ontology for modeling Go code structure, concurrency patterns, and code evolution. Designed for semantic code analysis and LLM-based code understanding.

## Purpose

Existing code ontologies target object-oriented languages (Java, C#) and cannot represent Go's unique constructs:

- **Package system** with import paths and visibility
- **Functions with receivers** (methods vs functions)
- **Goroutines and channels** (CSP concurrency model)
- **Interfaces** (structural typing vs nominal inheritance)
- **Explicit pointers** and value vs reference semantics
- **Generics** (type parameters, introduced in Go 1.18)
- **Defer, panic, recover** (unique error handling patterns)

This ontology fills that gap by modeling Go's specific semantics while aligning with established foundational ontologies (BFO, IAO) and provenance standards (PROV-O).

## Features

- **Parse Go source code** into RDF knowledge graphs
- **Model Go-specific constructs** (packages, functions, methods, structs, interfaces)
- **Represent concurrency patterns** (goroutines, channels, select statements)
- **Track code evolution** with git provenance and changesets
- **Validate graphs** using SHACL constraints
- **Export to Turtle format**

## Quick Start

```bash
# Add dependencies
go mod tidy

# Analyze a single file
go run cmd/go-ontology/main.go analyze file.go -o output.ttl

# Analyze a project
go run cmd/go-ontology/main.go analyze ./myproject -o project.ttl

# Output to stdout
go run cmd/go-ontology/main.go analyze file.go
```

## Ontology Architecture

Four-layer modular architecture:

```
go-core.ttl          → Base AST primitives, BFO/IAO alignment
    ↓ (owl:imports)
go-structure.ttl     → Go: Package, Function, Method, Struct, Interface
    ↓ (owl:imports)
go-concurrency.ttl   → Goroutines, Channels, Select, Sync primitives
    ↓ (owl:imports)
go-evolution.ttl     → PROV-O integration, version tracking

go-shapes.ttl        → SHACL validation (cross-cutting)
```

## Namespaces

| Prefix | IRI | Ontology |
|--------|-----|----------|
| core | `https://w3id.org/go-code/core#` | Core Ontology |
| struct | `https://w3id.org/go-code/structure#` | Structure Ontology |
| conc | `https://w3id.org/go-code/concurrency#` | Concurrency Ontology |
| evo | `https://w3id.org/go-code/evolution#` | Evolution Ontology |
| shapes | `https://w3id.org/go-code/shapes#` | Shapes Ontology |

## IRI Convention

```
Base: https://w3id.org/go-code/
Package:  {base}#fmt
Function: {base}#fmt.Println
Method:   {base}#net/http.Handler.ServeHTTP
Struct:   {base}#time.Time
Field:   {base}#time.Time/Second
```

## Design Principles

1. **No function overloading**: Go doesn't have overloading, so function identity is by name only
2. **Methods are functions**: Methods are functions with receivers, treated distinctly
3. **Interfaces are structural**: Type satisfaction is implicit, not declared
4. **Pointer receivers matter**: Value vs pointer receivers affect method sets
5. **Package-level visibility**: Exported = capital first letter

## Project Structure

```
go-ontology/
├── cmd/go-ontology/    # CLI entry point
├── pkg/
│   ├── analyzer/        # Parser and file analyzers
│   ├── extractor/       # AST extractors
│   ├── builder/         # RDF builders
│   ├── rdf/             # RDF data structures
│   ├── config/          # Configuration
│   └── pipeline/        # End-to-end pipeline
├── ontologies/          # Turtle ontology files
└── testdata/            # Test fixtures
```

## Development Status

- [x] Phase 1: Foundation (RDF types, Turtle serialization, config)
- [x] Phase 2: Parser and Analyzer (go/parser integration)
- [x] Phase 3: Basic Extractors (package, function)
- [x] Phase 4: Basic Builders (package, function)
- [x] Phase 5: Concurrency Support (ontology created)
- [x] Phase 6: CLI Integration (basic CLI implemented)
- [ ] Phase 7: Generic Support (Go 1.18+ type parameters)
- [ ] Phase 8: Testing and Documentation

## License

MIT
