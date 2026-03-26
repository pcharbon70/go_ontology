package builder

import (
	"github.com/pascal/ontology/pkg/rdf"
)

// Helper provides RDF building helper functions
type Helper struct{}

// NewHelper creates a new helper instance
func NewHelper() *Helper {
	return &Helper{}
}

// TypeTriple creates an rdf:type triple
func (h *Helper) TypeTriple(subject *rdf.IRI, class *rdf.IRI) *rdf.Triple {
	return rdf.NewTriple(subject, rdf.RDFType, class)
}

// StringProperty creates a string property triple (xsd:string)
func (h *Helper) StringProperty(
	subject *rdf.IRI,
	predicate *rdf.IRI,
	value string,
) *rdf.Triple {
	literal := rdf.NewLiteral(value)
	return rdf.NewTriple(subject, predicate, literal)
}

// IntegerProperty creates an integer property triple (xsd:integer)
func (h *Helper) IntegerProperty(
	subject *rdf.IRI,
	predicate *rdf.IRI,
	value int,
) *rdf.Triple {
	literal := rdf.NewIntegerLiteral(int64(value))
	return rdf.NewTriple(subject, predicate, literal)
}

// BooleanProperty creates a boolean property triple (xsd:boolean)
func (h *Helper) BooleanProperty(
	subject *rdf.IRI,
	predicate *rdf.IRI,
	value bool,
) *rdf.Triple {
	literal := rdf.NewBooleanLiteral(value)
	return rdf.NewTriple(subject, predicate, literal)
}

// ObjectProperty creates an object property triple
func (h *Helper) ObjectProperty(
	subject *rdf.IRI,
	predicate *rdf.IRI,
	object *rdf.IRI,
) *rdf.Triple {
	return rdf.NewTriple(subject, predicate, object)
}

// OptionalStringProperty creates a string property if value is not empty
func (h *Helper) OptionalStringProperty(
	subject *rdf.IRI,
	predicate *rdf.IRI,
	value string,
) *rdf.Triple {
	if value == "" {
		return nil
	}
	return h.StringProperty(subject, predicate, value)
}

// OptionalIRIProperty creates an object property if IRI is not nil
func (h *Helper) OptionalIRIProperty(
	subject *rdf.IRI,
	predicate *rdf.IRI,
	object *rdf.IRI,
) *rdf.Triple {
	if object == nil {
		return nil
	}
	return rdf.NewTriple(subject, predicate, object)
}

// NewIRI creates a new IRI with the builder's base
func (h *Helper) NewIRI(suffix string) *rdf.IRI {
	return rdf.NewIRI(suffix)
}
