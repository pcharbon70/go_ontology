package rdf

import (
	"fmt"
	"strings"
)

// Literal represents an RDF literal value
type Literal struct {
	Value    string
	Datatype *IRI
	Language string // If set, datatype should be rdf:langString
}

// NewLiteral creates a new literal with xsd:string datatype
func NewLiteral(value string) *Literal {
	return &Literal{
		Value:    value,
		Datatype: XSDString,
	}
}

// NewLiteralWithType creates a new literal with a specific datatype
func NewLiteralWithType(value string, datatype *IRI) *Literal {
	return &Literal{
		Value:    value,
		Datatype: datatype,
	}
}

// NewLiteralWithLanguage creates a new language-tagged literal
func NewLiteralWithLanguage(value, language string) *Literal {
	return &Literal{
		Value:    value,
		Datatype: RDFLangString,
		Language: language,
	}
}

// NewIntegerLiteral creates an integer literal
func NewIntegerLiteral(value int64) *Literal {
	return &Literal{
		Value:    fmt.Sprintf("%d", value),
		Datatype: XSDInteger,
	}
}

// NewBooleanLiteral creates a boolean literal
func NewBooleanLiteral(value bool) *Literal {
	val := "false"
	if value {
		val = "true"
	}
	return &Literal{
		Value:    val,
		Datatype: XSDBoolean,
	}
}

// String returns the Turtle representation of the literal
func (l *Literal) String() string {
	escaped := escapeString(l.Value)

	if l.Language != "" {
		return fmt.Sprintf(`"%s"@%s`, escaped, l.Language)
	}

	// Omit datatype for xsd:string (it's the default)
	if l.Datatype != nil && l.Datatype.Equals(XSDString) {
		return fmt.Sprintf(`"%s"`, escaped)
	}

	datatype := l.Datatype
	if datatype == nil {
		datatype = XSDString
	}

	return fmt.Sprintf(`"%s"^^%s`, escaped, datatype.String())
}

// escapeString escapes special characters in a string literal
func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}
