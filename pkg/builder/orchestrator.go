package builder

import (
	"github.com/pascal/ontology/pkg/extractor"
	"github.com/pascal/ontology/pkg/rdf"
)

// Orchestrator coordinates all builders to transform extraction results into RDF
type Orchestrator struct {
	builders []Builder
}

// NewOrchestrator creates a new builder orchestrator with default builders
func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		builders: []Builder{
			NewPackageBuilder(),
			NewFunctionBuilder(),
			NewStructBuilder(),
			NewInterfaceBuilder(),
		},
	}
}

// BuildAll processes all extraction results and returns a merged RDF graph
func (o *Orchestrator) BuildAll(results []extractor.ExtractionResult, ctx *Context) (*rdf.Graph, error) {
	graph := rdf.NewGraph()

	for _, result := range results {
		// Find appropriate builder
		for _, b := range o.builders {
			if b.CanBuild(result) {
				subGraph, err := b.Build(result, ctx)
				if err != nil {
					return nil, err
				}

				// Merge sub-graph into main graph
				graph = graph.Merge(subGraph)
				break // Only use first matching builder
			}
		}
	}

	return graph, nil
}

// AddBuilder adds a custom builder to the orchestrator
func (o *Orchestrator) AddBuilder(builder Builder) {
	o.builders = append(o.builders, builder)
}
