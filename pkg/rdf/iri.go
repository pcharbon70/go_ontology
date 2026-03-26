package rdf

import (
	"fmt"
	"strings"
)

// IRI represents an Internationalized Resource Identifier
type IRI struct {
	value string
}

// NewIRI creates a new IRI from a string
func NewIRI(value string) *IRI {
	return &IRI{value: value}
}

// String returns the string representation of the IRI
func (i *IRI) String() string {
	return i.value
}

// Equals checks if two IRIs are equal
func (i *IRI) Equals(other *IRI) bool {
	if i == nil || other == nil {
		return i == other
	}
	return i.value == other.value
}

// IsBlank checks if this is a blank node
func (i *IRI) IsBlank() bool {
	return strings.HasPrefix(i.value, "_:")
}

// Namespace returns the namespace part of the IRI (before the last # or /)
func (i *IRI) Namespace() string {
	if idx := strings.LastIndex(i.value, "#"); idx != -1 {
		return i.value[:idx+1]
	}
	if idx := strings.LastIndex(i.value, "/"); idx != -1 {
		return i.value[:idx+1]
	}
	return i.value
}

// Fragment returns the fragment part of the IRI (after # or /)
func (i *IRI) Fragment() string {
	if idx := strings.LastIndex(i.value, "#"); idx != -1 {
		return i.value[idx+1:]
	}
	if idx := strings.LastIndex(i.value, "/"); idx != -1 {
		return i.value[idx+1:]
	}
	return ""
}

// BlankNode represents a blank node (anonymous resource)
type BlankNode struct {
	ID string
}

// NewBlankNode creates a new blank node
func NewBlankNode() *BlankNode {
	return &BlankNode{ID: fmt.Sprintf("b%d", nextBlankID())}
}

// NewBlankNodeWithID creates a new blank node with a specific ID
func NewBlankNodeWithID(id string) *BlankNode {
	return &BlankNode{ID: id}
}

// String returns the string representation of the blank node
func (b *BlankNode) String() string {
	return "_:" + b.ID
}

var blankIDCounter uint64

func nextBlankID() uint64 {
	blankIDCounter++
	return blankIDCounter
}

// ResetBlankIDCounter resets the blank ID counter (useful for tests)
func ResetBlankIDCounter() {
	blankIDCounter = 0
}
