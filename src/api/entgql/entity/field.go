package entity

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (f Field) Validate() error {
	var err error
	if f.Unique != nil && *f.Unique && f.DefaultValue != nil {
		err = errors.New("cannot be unique and have default value set")
		return fmt.Errorf(validateErrorMessageTemplate, f.Caption, err)
	}

	if f.DefaultValue != nil {
		switch f.Type {
		case string(FieldTypeJson):
			err = ValidateJSONDefaultValue(*f.DefaultValue)
		}

		if err != nil {
			return fmt.Errorf(validateErrorMessageTemplate, f.Caption, err)
		}
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
