package rdf

// Namespace definitions for RDF ontologies
const (
	// RDF (RDF Syntax)
	RDF = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"

	// RDFS (RDF Schema)
	RDFS = "http://www.w3.org/2000/01/rdf-schema#"

	// XSD (XML Schema Datatypes)
	XSD = "http://www.w3.org/2001/XMLSchema#"

	// OWL (Web Ontology Language)
	OWL = "http://www.w3.org/2002/07/owl#"

	// DC (Dublin Core)
	DC = "http://purl.org/dc/elements/1.1/"

	// DCTERMS (Dublin Core Terms)
	DCTERMS = "http://purl.org/dc/terms/"

	// SKOS (Simple Knowledge Organization System)
	SKOS = "http://www.w3.org/2004/02/skos/core#"

	// BFO (Basic Formal Ontology)
	BFO = "http://purl.obolibrary.org/obo/"

	// IAO (Information Artifact Ontology)
	IAO = "http://purl.obolibrary.org/obo/IAO_"

	// PROV (Provenance Ontology)
	PROV = "http://www.w3.org/ns/prov#"
)

// Go Ontology Namespaces
const (
	// Go Code Ontology
	GoCore        = "https://w3id.org/go-code/core#"
	GoStructure   = "https://w3id.org/go-code/structure#"
	GoConcurrency = "https://w3id.org/go-code/concurrency#"
	GoEvolution   = "https://w3id.org/go-code/evolution#"
	GoShapes      = "https://w3id.org/go-code/shapes#"
)

// Predefined IRIs for common RDF terms
var (
	// RDF terms
	RDFType  = NewIRI(RDF + "type")
	RDFNil   = NewIRI(RDF + "nil")
	RDFFirst = NewIRI(RDF + "first")
	RDFRest  = NewIRI(RDF + "rest")

	// RDFS terms
	RDFSLabel       = NewIRI(RDFS + "label")
	RDFSComment     = NewIRI(RDFS + "comment")
	RDFSsubClassOf  = NewIRI(RDFS + "subClassOf")
	RDFSDomain      = NewIRI(RDFS + "domain")
	RDFSRange       = NewIRI(RDFS + "range")
	RDFSsubPropertyOf = NewIRI(RDFS + "subPropertyOf")

	// XSD datatypes
	XSDString   = NewIRI(XSD + "string")
	XSDInteger  = NewIRI(XSD + "integer")
	XSDBoolean  = NewIRI(XSD + "boolean")
	XSDDouble   = NewIRI(XSD + "double")
	XSDFloat    = NewIRI(XSD + "float")
	XSDDecimal  = NewIRI(XSD + "decimal")
	XSDDateTime = NewIRI(XSD + "dateTime")
	XSDDate     = NewIRI(XSD + "date")

	// OWL terms
	OWLOntology         = NewIRI(OWL + "Ontology")
	OWLClass            = NewIRI(OWL + "Class")
	OWLThing            = NewIRI(OWL + "Thing")
	OWLNothing          = NewIRI(OWL + "Nothing")
	OWLImports          = NewIRI(OWL + "imports")
	OWLversionInfo      = NewIRI(OWL + "versionInfo")
	OWLversionIRI       = NewIRI(OWL + "versionIRI")
	OWLFunctionalProperty = NewIRI(OWL + "FunctionalProperty")
	OWLObjectProperty   = NewIRI(OWL + "ObjectProperty")
	OWLDatatypeProperty = NewIRI(OWL + "DatatypeProperty")
	OWLinverseOf        = NewIRI(OWL + "inverseOf")
	OWLTransitiveProperty = NewIRI(OWL + "TransitiveProperty")
	OWLSymmetricProperty = NewIRI(OWL + "SymmetricProperty")

	// RDF terms for language-tagged strings
	RDFLangString = NewIRI(RDF + "langString")

	// Go Ontology classes
	GoFunction      = NewIRI(GoStructure + "Function")
	GoMethod        = NewIRI(GoStructure + "Method")
	GoStruct        = NewIRI(GoStructure + "Struct")
	GoInterface     = NewIRI(GoStructure + "Interface")
	GoPackage       = NewIRI(GoStructure + "Package")
	GoStructField   = NewIRI(GoStructure + "StructField")
	GoInterfaceMethod = NewIRI(GoStructure + "InterfaceMethod")
	GoTypeParameter = NewIRI(GoStructure + "TypeParameter")
	GoTypeAlias     = NewIRI(GoStructure + "TypeAlias")
	GoGoroutine     = NewIRI(GoConcurrency + "Goroutine")
	GoChannel       = NewIRI(GoConcurrency + "Channel")
	GoSelectStmt    = NewIRI(GoConcurrency + "Select")

	// Go Ontology properties (core)
	BelongsTo       = NewIRI(GoStructure + "belongsTo")
	HasParameter    = NewIRI(GoStructure + "hasParameter")
	HasResult       = NewIRI(GoStructure + "hasResult")
	HasReceiver     = NewIRI(GoStructure + "hasReceiver")
	HasField        = NewIRI(GoStructure + "hasField")
	HasTypeParam    = NewIRI(GoStructure + "hasTypeParameter")

	// Package properties
	PackageName     = NewIRI(GoStructure + "packageName")
	ImportPath      = NewIRI(GoStructure + "importPath")
	ImportsPackage  = NewIRI(GoStructure + "importsPackage")

	// Function properties
	FunctionName    = NewIRI(GoStructure + "functionName")
	ReceiverName    = NewIRI(GoStructure + "receiverName")
	ReceiverType    = NewIRI(GoStructure + "receiverType")
	ParameterCount  = NewIRI(GoStructure + "parameterCount")
	ResultCount     = NewIRI(GoStructure + "resultCount")
	IsVariadic      = NewIRI(GoStructure + "isVariadic")
	IsExported      = NewIRI(GoStructure + "isExported")

	// Struct properties
	StructName      = NewIRI(GoStructure + "structName")
	FieldName       = NewIRI(GoStructure + "fieldName")
	FieldType       = NewIRI(GoStructure + "fieldType")
	FieldTag        = NewIRI(GoStructure + "fieldTag")
	IsEmbedded      = NewIRI(GoStructure + "isEmbedded")

	// Type parameter properties
	TypeParamName   = NewIRI(GoStructure + "typeParamName")
	TypeParamConstraint = NewIRI(GoStructure + "typeParamConstraint")

	// Concurrency properties
	GoroutineCount  = NewIRI(GoConcurrency + "goroutineCount")
	ChannelType     = NewIRI(GoConcurrency + "channelType")
	ChannelElement  = NewIRI(GoConcurrency + "elementType")
	ChannelDirection = NewIRI(GoConcurrency + "direction")
)

// PrefixMap returns a map of prefixes to namespace IRIs for Turtle serialization
func PrefixMap() map[string]string {
	return map[string]string{
		"rdf":       RDF,
		"rdfs":      RDFS,
		"xsd":       XSD,
		"owl":       OWL,
		"dc":        DC,
		"dcterms":   DCTERMS,
		"skos":      SKOS,
		"bfo":       BFO,
		"iao":       IAO,
		"prov":      PROV,
		"core":      GoCore,
		"struct":    GoStructure,
		"conc":      GoConcurrency,
		"evo":       GoEvolution,
		"shapes":    GoShapes,
	}
}
