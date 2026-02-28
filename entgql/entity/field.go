package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Field struct {
	Name           string
	EntName        string
	Caption        string
	Type           string
	Required       *bool
	Unique         *bool
	DefaultValue   *string
	Min            *string
	Max            *string
	Private        *bool
	Nillable       bool
	UpdateDefault  bool
	AcceptedValues []string
}

const validateErrorMessageTemplate = "[Validate] field %s invalid: %w"

// Reserved names that cannot be used for field names
var reservedFieldNames = map[string]bool{
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

// Field types that support min/max constraints
var numericFieldTypes = map[string]bool{
	string(FieldTypeInteger): true,
	string(FieldTypeDecimal): true,
	string(FieldTypeFloat):   true,
}

var dateTimeFieldTypes = map[string]bool{
	string(FieldTypeDateTime): true,
	string(FieldTypeDate):     true,
	string(FieldTypeTime):     true,
}

// Valid field types
var validFieldTypes = map[string]bool{
	string(FieldTypeID):             true,
	string(FieldTypeShortText):      true,
	string(FieldTypeLongText):       true,
	string(FieldTypeRichText):       true,
	string(FieldTypeFileContent):    true,
	string(FieldTypeEmail):          true,
	string(FieldTypeInteger):        true,
	string(FieldTypeDecimal):        true,
	string(FieldTypeFloat):          true,
	string(FieldTypeBoolean):        true,
	string(FieldTypeSingleChoice):   true,
	string(FieldTypeMultipleChoice): true,
	string(FieldTypeDateTime):       true,
	string(FieldTypeDate):           true,
	string(FieldTypeTime):           true,
	string(FieldTypeJson):           true,
	string(FieldTypeEnum):           true,
	string(FieldTypeEnums):          true,
}

// Field types that support AcceptedValues
var choiceFieldTypes = map[string]bool{
	string(FieldTypeSingleChoice):   true,
	string(FieldTypeMultipleChoice): true,
}

// Field types that support UpdateDefault
var updateDefaultFieldTypes = map[string]bool{
	string(FieldTypeDateTime): true,
	string(FieldTypeDate):     true,
	string(FieldTypeTime):     true,
}

// Valid identifier regex: starts with letter or underscore, followed by letters, digits, or underscores
var validIdentifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func (f Field) Validate() error {
	var err error

	// Validate name is not reserved
	if err = validateName(f.Name, reservedFieldNames, "field"); err != nil {
		return fmt.Errorf(validateErrorMessageTemplate, f.Caption, err)
	}

	// Validate name is a valid identifier
	if !validIdentifierRegex.MatchString(f.Name) {
		err = fmt.Errorf("name '%s' is not a valid identifier (must start with letter or underscore, contain only letters, digits, and underscores)", f.Name)
		return fmt.Errorf(validateErrorMessageTemplate, f.Caption, err)
	}

	// Validate Caption is not empty
	if f.Caption == "" {
		err = errors.New("caption cannot be empty")
		return fmt.Errorf(validateErrorMessageTemplate, f.Name, err)
	}

	// Validate Type is a valid field type
	if !validFieldTypes[f.Type] {
		err = fmt.Errorf("invalid field type '%s'", f.Type)
		return fmt.Errorf(validateErrorMessageTemplate, f.Caption, err)
	}

	// Validate AcceptedValues only for choice field types
	if len(f.AcceptedValues) > 0 && !choiceFieldTypes[f.Type] {
		err = fmt.Errorf("acceptedValues can only be set for SingleChoice and MultipleChoice field types, got %s", f.Type)
		return fmt.Errorf(validateErrorMessageTemplate, f.Caption, err)
	}

	// Validate AcceptedValues is not empty for choice field types
	if choiceFieldTypes[f.Type] && len(f.AcceptedValues) == 0 {
		err = errors.New("acceptedValues must be provided for SingleChoice and MultipleChoice field types")
		return fmt.Errorf(validateErrorMessageTemplate, f.Caption, err)
	}

	// Validate UpdateDefault only for date/time field types
	if f.UpdateDefault && !updateDefaultFieldTypes[f.Type] {
		err = fmt.Errorf("updateDefault can only be set for DateTime, Date, and Time field types, got %s", f.Type)
		return fmt.Errorf(validateErrorMessageTemplate, f.Caption, err)
	}

	// Validate unique and default value cannot be set together
	if f.Unique != nil && *f.Unique && f.DefaultValue != nil {
		err = errors.New("cannot be unique and have default value set")
		return fmt.Errorf(validateErrorMessageTemplate, f.Caption, err)
	}

	// Validate min/max only for numeric and date/time field types
	if f.Min != nil || f.Max != nil {
		if !numericFieldTypes[f.Type] && !dateTimeFieldTypes[f.Type] {
			err = fmt.Errorf("min/max constraints are only supported for numeric (Integer, Decimal, Float) and date/time (DateTime, Date, Time) field types, got %s", f.Type)
			return fmt.Errorf(validateErrorMessageTemplate, f.Caption, err)
		}
	}

	// Validate min < max for numeric types
	if numericFieldTypes[f.Type] && f.Min != nil && f.Max != nil {
		minVal, errMin := parseNumericValue(*f.Min)
		maxVal, errMax := parseNumericValue(*f.Max)
		if errMin == nil && errMax == nil {
			if minVal >= maxVal {
				err = fmt.Errorf("min value (%s) must be less than max value (%s)", *f.Min, *f.Max)
				return fmt.Errorf(validateErrorMessageTemplate, f.Caption, err)
			}
		}
	}

	// Validate default value format based on field type
	if f.DefaultValue != nil {
		switch f.Type {
		case string(FieldTypeJson):
			err = ValidateJSONDefaultValue(*f.DefaultValue)
		case string(FieldTypeBoolean):
			if *f.DefaultValue != "true" && *f.DefaultValue != "false" {
				err = fmt.Errorf("default value for Boolean field must be 'true' or 'false', got '%s'", *f.DefaultValue)
			}
		case string(FieldTypeInteger):
			if _, err = strconv.ParseInt(*f.DefaultValue, 10, 64); err != nil {
				err = fmt.Errorf("default value for Integer field must be a valid integer, got '%s'", *f.DefaultValue)
			}
		case string(FieldTypeFloat), string(FieldTypeDecimal):
			if _, err = strconv.ParseFloat(*f.DefaultValue, 64); err != nil {
				err = fmt.Errorf("default value for %s field must be a valid number, got '%s'", f.Type, *f.DefaultValue)
			}
		}

		if err != nil {
			return fmt.Errorf(validateErrorMessageTemplate, f.Caption, err)
		}
	}

	return nil
}

// parseNumericValue attempts to parse a string as a numeric value
func parseNumericValue(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

// validateName checks if a name is reserved
func validateName(name string, reservedNames map[string]bool, entityType string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if reservedNames[name] {
		return fmt.Errorf("name '%s' is reserved and cannot be used for %s", name, entityType)
	}
	return nil
}

func (f Field) Equals(uniqueInput FieldWhereUniqueInput) bool {
	if uniqueInput.Caption != nil && f.Caption == *uniqueInput.Caption {
		return true
	}

	if uniqueInput.Name != nil && f.Name == *uniqueInput.Name {
		return true
	}

	return false
}

func (f *Field) ApplyUpdateInput(data UpdateFieldInput) {
	if data.Caption != nil {
		f.Caption = *data.Caption
	}

	if data.Required != nil {
		f.Required = data.Required
	}

	if data.Unique != nil {
		f.Unique = data.Unique
	}

	if data.DefaultValue != nil {
		f.DefaultValue = data.DefaultValue
	}

	if data.Min != nil {
		f.Min = data.Min
	}

	if data.Max != nil {
		f.Max = data.Max
	}

	if data.Private != nil {
		f.Private = data.Private
	}

	if data.AcceptedValues != nil {
		f.AcceptedValues = data.AcceptedValues
	}
}

type FieldType string

const (
	FieldTypeID             FieldType = "ID"
	FieldTypeShortText      FieldType = "ShortText"
	FieldTypeLongText       FieldType = "LongText"
	FieldTypeRichText       FieldType = "RichText"
	FieldTypeFileContent    FieldType = "FileContent"
	FieldTypeEmail          FieldType = "Email"
	FieldTypeInteger        FieldType = "Integer"
	FieldTypeDecimal        FieldType = "Decimal"
	FieldTypeFloat          FieldType = "Float"
	FieldTypeBoolean        FieldType = "Boolean"
	FieldTypeSingleChoice   FieldType = "SingleChoice"
	FieldTypeMultipleChoice FieldType = "MultipleChoice"
	FieldTypeDateTime       FieldType = "DateTime"
	FieldTypeDate           FieldType = "Date"
	FieldTypeTime           FieldType = "Time"
	FieldTypeJson           FieldType = "Json"
	FieldTypeEnum           FieldType = "Enum"
	FieldTypeEnums          FieldType = "Enums"
)

func ValidateEmail(val string) error {
	// TODO implement
	return nil
}

func ValidatePhone(val string) error {
	// TODO implement
	return nil
}

func ValidateJSONDefaultValue(val string) error {
	_, err := JSONDefaultValueWithError[map[string]interface{}](val) // TODO any other type than map?
	if err != nil {
		return fmt.Errorf("json format invalid: %w", err)
	}
	return err
}

func JSONDefaultValueWithError[T any](val string) (T, error) {
	var value T

	err := json.Unmarshal([]byte(val), &value)

	return value, err
}

func JSONDefaultValue[T any](val string) T {
	value, _ := JSONDefaultValueWithError[T](val)
	return value
}
