package rdf

import "fmt"

// Triple represents an RDF triple (subject-predicate-object)
type Triple struct {
	Subject   Node
	Predicate *IRI
	Object    Node
}

// Node represents any RDF node (IRI, Literal, or BlankNode)
type Node interface {
	isNode()
}

func (*IRI) isNode()     {}
func (*Literal) isNode() {}
func (*BlankNode) isNode() {}

// NewTriple creates a new triple
func NewTriple(subject Node, predicate *IRI, object Node) *Triple {
	return &Triple{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
	}
}

// String returns the Turtle representation of the triple
func (t *Triple) String() string {
	return fmt.Sprintf("%s %s %s .", nodeToString(t.Subject), t.Predicate.String(), nodeToString(t.Object))
}

func nodeToString(n Node) string {
	switch v := n.(type) {
	case *IRI:
		return v.String()
	case *Literal:
		return v.String()
	case *BlankNode:
		return v.String()
	default:
		return ""
	}
}

// TripleWithIRI creates a triple with an IRI object
func TripleWithIRI(subject *IRI, predicate *IRI, object *IRI) *Triple {
	return &Triple{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
	}
}

// TripleWithLiteral creates a triple with a literal object
func TripleWithLiteral(subject *IRI, predicate *IRI, object *Literal) *Triple {
	return &Triple{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
	}
}
