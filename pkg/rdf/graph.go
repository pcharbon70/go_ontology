package rdf

import (
	"sort"
	"sync"
)

// Graph represents an RDF graph as a collection of triples
type Graph struct {
	triples []*Triple
	mu      sync.RWMutex
}

// NewGraph creates a new empty RDF graph
func NewGraph() *Graph {
	return &Graph{
		triples: make([]*Triple, 0),
	}
}

// Add adds one or more triples to the graph
func (g *Graph) Add(triples ...*Triple) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, t := range triples {
		if t != nil {
			g.triples = append(g.triples, t)
		}
	}
}

// Triples returns all triples in the graph
func (g *Graph) Triples() []*Triple {
	g.mu.RLock()
	defer g.mu.RUnlock()

	result := make([]*Triple, len(g.triples))
	copy(result, g.triples)
	return result
}

// Count returns the number of triples in the graph
func (g *Graph) Count() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.triples)
}

// Merge combines this graph with another graph
func (g *Graph) Merge(other *Graph) *Graph {
	g.mu.RLock()
	defer g.mu.RUnlock()

	other.mu.RLock()
	defer other.mu.RUnlock()

	result := NewGraph()
	result.triples = make([]*Triple, 0, len(g.triples)+len(other.triples))
	result.triples = append(result.triples, g.triples...)
	result.triples = append(result.triples, other.triples...)

	return result
}

// Subjects returns all unique subjects in the graph
func (g *Graph) Subjects() []*IRI {
	g.mu.RLock()
	defer g.mu.RUnlock()

	seen := make(map[string]bool)
	var subjects []*IRI

	for _, t := range g.triples {
		if iri, ok := t.Subject.(*IRI); ok {
			if !seen[iri.String()] {
				seen[iri.String()] = true
				subjects = append(subjects, iri)
			}
		}
	}

	sort.Slice(subjects, func(i, j int) bool {
		return subjects[i].String() < subjects[j].String()
	})

	return subjects
}

// Objects returns all unique objects in the graph that are IRIs
func (g *Graph) Objects() []*IRI {
	g.mu.RLock()
	defer g.mu.RUnlock()

	seen := make(map[string]bool)
	var objects []*IRI

	for _, t := range g.triples {
		if iri, ok := t.Object.(*IRI); ok {
			if !seen[iri.String()] {
				seen[iri.String()] = true
				objects = append(objects, iri)
			}
		}
	}

	sort.Slice(objects, func(i, j int) bool {
		return objects[i].String() < objects[j].String()
	})

	return objects
}

// TriplesForSubject returns all triples with the given subject
func (g *Graph) TriplesForSubject(subject *IRI) []*Triple {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var result []*Triple
	for _, t := range g.triples {
		if iri, ok := t.Subject.(*IRI); ok && iri.Equals(subject) {
			result = append(result, t)
		}
	}
	return result
}

// TriplesForPredicate returns all triples with the given predicate
func (g *Graph) TriplesForPredicate(predicate *IRI) []*Triple {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var result []*Triple
	for _, t := range g.triples {
		if t.Predicate.Equals(predicate) {
			result = append(result, t)
		}
	}
	return result
}

// TriplesForObject returns all triples with the given object
func (g *Graph) TriplesForObject(object Node) []*Triple {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var result []*Triple
	for _, t := range g.triples {
		switch v := t.Object.(type) {
		case *IRI:
			if iri, ok := object.(*IRI); ok && v.Equals(iri) {
				result = append(result, t)
			}
		case *Literal:
			if lit, ok := object.(*Literal); ok && v.Value == lit.Value {
				result = append(result, t)
			}
		case *BlankNode:
			if bn, ok := object.(*BlankNode); ok && v.ID == bn.ID {
				result = append(result, t)
			}
		}
	}
	return result
}

// Describe returns all triples where the subject is the given IRI
func (g *Graph) Describe(subject *IRI) []*Triple {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var result []*Triple
	for _, t := range g.triples {
		if iri, ok := t.Subject.(*IRI); ok && iri.Equals(subject) {
			result = append(result, t)
		}
	}
	return result
}

// GroupBySubject groups triples by their subject
func (g *Graph) GroupBySubject() map[string][]*Triple {
	g.mu.RLock()
	defer g.mu.RUnlock()

	groups := make(map[string][]*Triple)
	for _, t := range g.triples {
		subjStr := nodeToString(t.Subject)
		groups[subjStr] = append(groups[subjStr], t)
	}
	return groups
}

// Clear removes all triples from the graph
func (g *Graph) Clear() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.triples = make([]*Triple, 0)
}

// Clone creates a deep copy of the graph
func (g *Graph) Clone() *Graph {
	g.mu.RLock()
	defer g.mu.RUnlock()

	clone := NewGraph()
	clone.triples = make([]*Triple, len(g.triples))
	copy(clone.triples, g.triples)
	return clone
}
