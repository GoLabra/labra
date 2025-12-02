package entity

import (
	"fmt"
	"regexp"
)

type Edge struct {
	Name             string
	EntName          string
	Caption          string
	Type             string
	BelongsToCaption *string
	Required         *bool
	RelationType     RelationType
	Private          *bool
	Ref              string
}

// Reserved names that cannot be used for edge names
var reservedEdgeNames = map[string]bool{
	"id":          true, // Reserved by ent
	"edges":       true, // Reserved by ent
	"type":        true, // Go keyword
	"func":        true, // Go keyword
	"var":         true, // Go keyword
	"const":       true, // Go keyword
	"package":     true, // Go keyword
	"import":      true, // Go keyword
	"return":      true, // Go keyword
	"if":          true, // Go keyword
	"else":        true, // Go keyword
	"for":         true, // Go keyword
	"range":       true, // Go keyword
	"switch":      true, // Go keyword
	"case":        true, // Go keyword
	"default":     true, // Go keyword
	"break":       true, // Go keyword
	"continue":    true, // Go keyword
	"fallthrough": true, // Go keyword
	"defer":       true, // Go keyword
	"go":          true, // Go keyword
	"goto":        true, // Go keyword
	"select":      true, // Go keyword
	"chan":        true, // Go keyword
	"interface":   true, // Go keyword
	"struct":      true, // Go keyword
	"map":         true, // Go keyword
	"string":      true, // Go built-in type
	"int":         true, // Go built-in type
	"bool":        true, // Go built-in type
	"float":       true, // Go built-in type
	"byte":        true, // Go built-in type
	"rune":        true, // Go built-in type
	"error":       true, // Go built-in type
	"nil":         true, // Go built-in
	"true":        true, // Go built-in
	"false":       true, // Go built-in
	"make":        true, // Go built-in function
	"new":         true, // Go built-in function
	"len":         true, // Go built-in function
	"cap":         true, // Go built-in function
	"append":      true, // Go built-in function
	"copy":        true, // Go built-in function
	"delete":      true, // Go built-in function
	"panic":       true, // Go built-in function
	"recover":     true, // Go built-in function
	"close":       true, // Go built-in function
	"complex":     true, // Go built-in function
	"real":        true, // Go built-in function
	"imag":        true, // Go built-in function
	"print":       true, // Go built-in function
	"println":     true, // Go built-in function
}

// Valid identifier regex: starts with letter or underscore, followed by letters, digits, or underscores
var validEdgeIdentifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func (e Edge) Validate() error {
	// Validate name is not reserved
	if err := validateName(e.Name, reservedEdgeNames, "edge"); err != nil {
		return fmt.Errorf("[Validate] edge %s invalid: %w", e.Caption, err)
	}

	// Validate name is a valid identifier
	if !validEdgeIdentifierRegex.MatchString(e.Name) {
		return fmt.Errorf("[Validate] edge %s invalid: name '%s' is not a valid identifier (must start with letter or underscore, contain only letters, digits, and underscores)", e.Caption, e.Name)
	}

	// Validate Caption is not empty
	if e.Caption == "" {
		return fmt.Errorf("[Validate] edge %s invalid: Caption cannot be empty", e.Name)
	}

	// Validate RelationType is valid
	if !e.RelationType.IsValid() {
		return fmt.Errorf("[Validate] edge %s invalid: RelationType '%s' is not a valid RelationType", e.Caption, e.RelationType)
	}

	return nil
}

func (f Edge) Equals(uniqueInput EdgeWhereUniqueInput) bool {
	if uniqueInput.Caption != nil && f.Caption == *uniqueInput.Caption {
		return true
	}

	if uniqueInput.Name != nil && f.Name == *uniqueInput.Name {
		return true
	}

	return false
}

func (e *Edge) ApplyUpdateInput(data UpdateEdgeInput) {
	if data.Caption != nil {
		e.Caption = *data.Caption
	}

	if data.Required != nil {
		e.Required = data.Required
	}

	if data.Private != nil {
		e.Private = data.Private
	}
}
