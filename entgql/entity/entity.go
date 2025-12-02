package entity

import (
	"fmt"
	"regexp"
)

type Entity struct {
	Name             string
	EntName          string
	Caption          string
	Owner            EntityOwner
	DisplayFieldName string
}

// Reserved names that cannot be used for entity names
var reservedEntityNames = map[string]bool{
	"entity":       true, // Reserved by entgql
	"entities":     true, // Reserved by entgql
	"field":        true, // Reserved by entgql
	"fields":       true, // Reserved by entgql
	"edge":         true, // Reserved by entgql
	"edges":        true, // Reserved by entgql
	"query":        true, // GraphQL reserved
	"mutation":     true, // GraphQL reserved
	"subscription": true, // GraphQL reserved
	"type":         true, // Go keyword
	"func":         true, // Go keyword
	"var":          true, // Go keyword
	"const":        true, // Go keyword
	"package":      true, // Go keyword
	"import":       true, // Go keyword
	"return":       true, // Go keyword
	"if":           true, // Go keyword
	"else":         true, // Go keyword
	"for":          true, // Go keyword
	"range":        true, // Go keyword
	"switch":       true, // Go keyword
	"case":         true, // Go keyword
	"default":      true, // Go keyword
	"break":        true, // Go keyword
	"continue":     true, // Go keyword
	"fallthrough":  true, // Go keyword
	"defer":        true, // Go keyword
	"go":           true, // Go keyword
	"goto":         true, // Go keyword
	"select":       true, // Go keyword
	"chan":         true, // Go keyword
	"interface":    true, // Go keyword
	"struct":       true, // Go keyword
	"map":          true, // Go keyword
	"string":       true, // Go built-in type
	"int":          true, // Go built-in type
	"bool":         true, // Go built-in type
	"float":        true, // Go built-in type
	"byte":         true, // Go built-in type
	"rune":         true, // Go built-in type
	"error":        true, // Go built-in type
	"nil":          true, // Go built-in
	"true":         true, // Go built-in
	"false":        true, // Go built-in
	"make":         true, // Go built-in function
	"new":          true, // Go built-in function
	"len":          true, // Go built-in function
	"cap":          true, // Go built-in function
	"append":       true, // Go built-in function
	"copy":         true, // Go built-in function
	"delete":       true, // Go built-in function
	"panic":        true, // Go built-in function
	"recover":      true, // Go built-in function
	"close":        true, // Go built-in function
	"complex":      true, // Go built-in function
	"real":         true, // Go built-in function
	"imag":         true, // Go built-in function
	"print":        true, // Go built-in function
	"println":      true, // Go built-in function
}

// Valid identifier regex: starts with letter or underscore, followed by letters, digits, or underscores
var validEntityIdentifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func (f Entity) Validate() error {
	// Validate name is not reserved
	if err := validateName(f.Name, reservedEntityNames, "entity"); err != nil {
		return fmt.Errorf("[Validate] entity %s invalid: %w", f.Caption, err)
	}

	// Validate name is a valid identifier
	if !validEntityIdentifierRegex.MatchString(f.Name) {
		return fmt.Errorf("[Validate] entity %s invalid: name '%s' is not a valid identifier (must start with letter or underscore, contain only letters, digits, and underscores)", f.Caption, f.Name)
	}

	// Validate EntName is not empty
	if f.EntName == "" {
		return fmt.Errorf("[Validate] entity %s invalid: EntName cannot be empty", f.Caption)
	}

	// Validate EntName is a valid identifier (PascalCase typically)
	if !validEntityIdentifierRegex.MatchString(f.EntName) {
		return fmt.Errorf("[Validate] entity %s invalid: EntName '%s' is not a valid identifier (must start with letter or underscore, contain only letters, digits, and underscores)", f.Caption, f.EntName)
	}

	// Validate Caption is not empty
	if f.Caption == "" {
		return fmt.Errorf("[Validate] entity %s invalid: Caption cannot be empty", f.Name)
	}

	// Validate Owner is valid
	if !f.Owner.IsValid() {
		return fmt.Errorf("[Validate] entity %s invalid: Owner '%s' is not a valid EntityOwner", f.Caption, f.Owner)
	}

	return nil
}

func (f Entity) Equals(uniqueInput EntityWhereUniqueInput) bool {
	if uniqueInput.Caption != nil && f.Caption == *uniqueInput.Caption {
		return true
	}

	if uniqueInput.Name != nil && f.Name == *uniqueInput.Name {
		return true
	}

	return false
}

func (e *Entity) ApplyUpdateInput(data UpdateEntityInput) {
	if data.Caption != nil {
		e.Caption = *data.Caption
	}
}
