package rdf

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// TurtleSerializer serializes an RDF graph to Turtle format
type TurtleSerializer struct {
	prefixes map[string]string
}

// NewTurtleSerializer creates a new Turtle serializer
func NewTurtleSerializer() *TurtleSerializer {
	return &TurtleSerializer{
		prefixes: PrefixMap(),
	}
}

// Serialize writes the graph to Turtle format
func (s *TurtleSerializer) Serialize(w io.Writer, graph *Graph) error {
	// Write prefixes
	if err := s.writePrefixes(w); err != nil {
		return err
	}

	// Group triples by subject
	groups := graph.GroupBySubject()

	// Sort subjects for consistent output
	sortedSubjects := make([]string, 0, len(groups))
	for subj := range groups {
		sortedSubjects = append(sortedSubjects, subj)
	}
	sort.Strings(sortedSubjects)

	// Write triples grouped by subject
	for _, subject := range sortedSubjects {
		if err := s.writeSubjectTriples(w, subject, groups[subject]); err != nil {
			return err
		}
	}

	return nil
}

// SerializeToString returns the Turtle representation as a string
func (s *TurtleSerializer) SerializeToString(graph *Graph) (string, error) {
	var sb strings.Builder
	if err := s.Serialize(&sb, graph); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (s *TurtleSerializer) writePrefixes(w io.Writer) error {
	prefixes := make([]string, 0, len(s.prefixes))
	for k := range s.prefixes {
		prefixes = append(prefixes, k)
	}
	sort.Strings(prefixes)

	for _, k := range prefixes {
		if _, err := fmt.Fprintf(w, "@prefix %s: <%s> .\n", k, s.prefixes[k]); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(w)
	return err
}

func (s *TurtleSerializer) writeSubjectTriples(w io.Writer, subject string, triples []*Triple) error {
	// Write subject
	fmt.Fprintf(w, "<%s> ", subject)

	// Write predicates and objects
	for i, t := range triples {
		if i > 0 {
			fmt.Fprintf(w, " ;\n    ")
		}

		// Write predicate (abbreviated if possible)
		pred := s.abbreviate(t.Predicate.String())
		fmt.Fprintf(w, "%s ", pred)

		// Write object
		switch o := t.Object.(type) {
		case *IRI:
			obj := s.abbreviate(o.String())
			fmt.Fprintf(w, "%s", obj)
		case *Literal:
			fmt.Fprintf(w, "%s", o.String())
		case *BlankNode:
			fmt.Fprintf(w, "%s", o.String())
		}
	}

	_, err := fmt.Fprintln(w, " .")
	return err
}

func (s *TurtleSerializer) abbreviate(iri string) string {
	// Check if we can abbreviate with a prefix
	for prefix, expansion := range s.prefixes {
		if strings.HasPrefix(iri, expansion) {
			suffix := iri[len(expansion):]
			// Use prefixed name if the suffix is simple (no special chars)
			if isSafePrefix(suffix) {
				return fmt.Sprintf("%s:%s", prefix, suffix)
			}
		}
	}
	// Full IRI
	return fmt.Sprintf("<%s>", iri)
}

// isSafePrefix checks if a string is safe to use as a prefixed name
func isSafePrefix(s string) bool {
	if s == "" {
		return false
	}
	// Allow alphanumeric, underscore, dot, dash
	for _, c := range s {
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '_' || c == '.' || c == '-') {
			return false
		}
	}
	return true
}

// ToTurtle serializes a graph to Turtle format (convenience function)
func ToTurtle(graph *Graph) (string, error) {
	serializer := NewTurtleSerializer()
	return serializer.SerializeToString(graph)
}

// MustToTurtle serializes a graph to Turtle format, panicking on error
func MustToTurtle(graph *Graph) string {
	result, err := ToTurtle(graph)
	if err != nil {
		panic(err)
	}
	return result
}
