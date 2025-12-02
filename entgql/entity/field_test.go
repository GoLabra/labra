package entity

import (
	"testing"
)

func TestField_Validate(t *testing.T) {
	tests := []struct {
		name    string
		field   Field
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid field with no constraints",
			field: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeShortText),
			},
			wantErr: false,
		},
		{
			name: "valid field with unique constraint",
			field: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeShortText),
				Unique:  boolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "valid field with default value",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeShortText),
				DefaultValue: stringPtr("default"),
			},
			wantErr: false,
		},
		{
			name: "invalid: unique and default value together",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeShortText),
				Unique:       boolPtr(true),
				DefaultValue: stringPtr("default"),
			},
			wantErr: true,
			errMsg:  "cannot be unique and have default value set",
		},
		{
			name: "valid: unique false and default value",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeShortText),
				Unique:       boolPtr(false),
				DefaultValue: stringPtr("default"),
			},
			wantErr: false,
		},
		{
			name: "valid JSON default value",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeJson),
				DefaultValue: stringPtr(`{"key": "value"}`),
			},
			wantErr: false,
		},
		{
			name: "invalid JSON default value",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeJson),
				DefaultValue: stringPtr(`{"key": "value"`), // Invalid JSON
			},
			wantErr: true,
			errMsg:  "json format invalid",
		},
		{
			name: "invalid JSON array default value (ValidateJSONDefaultValue only validates map[string]interface{})",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeJson),
				DefaultValue: stringPtr(`["item1", "item2"]`),
			},
			wantErr: true, // Arrays are not map[string]interface{}, so validation fails
			errMsg:  "json format invalid",
		},
		{
			name: "valid JSON with nested objects",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeJson),
				DefaultValue: stringPtr(`{"nested": {"key": "value"}}`),
			},
			wantErr: false,
		},
		{
			name: "non-JSON field with invalid JSON default (should not validate JSON)",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeShortText),
				DefaultValue: stringPtr(`{"key": "value"`), // Invalid JSON but not validated for non-JSON fields
			},
			wantErr: false,
		},
		{
			name: "field with reserved name 'id'",
			field: Field{
				Name:    "id",
				Caption: "ID",
				Type:    string(FieldTypeShortText),
			},
			wantErr: true,
			errMsg:  "reserved",
		},
		{
			name: "field with reserved name 'type'",
			field: Field{
				Name:    "type",
				Caption: "Type",
				Type:    string(FieldTypeShortText),
			},
			wantErr: true,
			errMsg:  "reserved",
		},
		{
			name: "field with empty name",
			field: Field{
				Name:    "",
				Caption: "Test",
				Type:    string(FieldTypeShortText),
			},
			wantErr: true,
			errMsg:  "cannot be empty",
		},
		{
			name: "field with min on non-numeric type (ShortText)",
			field: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeShortText),
				Min:     stringPtr("0"),
			},
			wantErr: true,
			errMsg:  "min/max constraints are only supported",
		},
		{
			name: "field with max on non-numeric type (Email)",
			field: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeEmail),
				Max:     stringPtr("100"),
			},
			wantErr: true,
			errMsg:  "min/max constraints are only supported",
		},
		{
			name: "field with min on Integer (valid)",
			field: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeInteger),
				Min:     stringPtr("0"),
			},
			wantErr: false,
		},
		{
			name: "field with max on Float (valid)",
			field: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeFloat),
				Max:     stringPtr("100.0"),
			},
			wantErr: false,
		},
		{
			name: "field with min/max on Decimal (valid)",
			field: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeDecimal),
				Min:     stringPtr("0.0"),
				Max:     stringPtr("100.0"),
			},
			wantErr: false,
		},
		{
			name: "field with min on DateTime (valid)",
			field: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeDateTime),
				Min:     stringPtr("2020-01-01"),
			},
			wantErr: false,
		},
		{
			name: "field with max on Date (valid)",
			field: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeDate),
				Max:     stringPtr("2020-12-31"),
			},
			wantErr: false,
		},
		{
			name: "field with min on Time (valid)",
			field: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeTime),
				Min:     stringPtr("00:00:00"),
			},
			wantErr: false,
		},
		{
			name: "field with invalid name starting with number",
			field: Field{
				Name:    "123field",
				Caption: "Test",
				Type:    string(FieldTypeShortText),
			},
			wantErr: true,
			errMsg:  "not a valid identifier",
		},
		{
			name: "field with invalid name with special characters",
			field: Field{
				Name:    "field-name",
				Caption: "Test",
				Type:    string(FieldTypeShortText),
			},
			wantErr: true,
			errMsg:  "not a valid identifier",
		},
		{
			name: "field with empty caption",
			field: Field{
				Name:    "test",
				Caption: "",
				Type:    string(FieldTypeShortText),
			},
			wantErr: true,
			errMsg:  "caption cannot be empty",
		},
		{
			name: "field with invalid type",
			field: Field{
				Name:    "test",
				Caption: "Test",
				Type:    "InvalidType",
			},
			wantErr: true,
			errMsg:  "invalid field type",
		},
		{
			name: "field with AcceptedValues on non-choice type",
			field: Field{
				Name:           "test",
				Caption:        "Test",
				Type:           string(FieldTypeShortText),
				AcceptedValues: []string{"value1", "value2"},
			},
			wantErr: true,
			errMsg:  "acceptedValues can only be set for SingleChoice and MultipleChoice",
		},
		{
			name: "field with SingleChoice but no AcceptedValues",
			field: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeSingleChoice),
			},
			wantErr: true,
			errMsg:  "acceptedValues must be provided",
		},
		{
			name: "field with MultipleChoice but no AcceptedValues",
			field: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeMultipleChoice),
			},
			wantErr: true,
			errMsg:  "acceptedValues must be provided",
		},
		{
			name: "field with UpdateDefault on non-date/time type",
			field: Field{
				Name:          "test",
				Caption:       "Test",
				Type:          string(FieldTypeShortText),
				UpdateDefault: true,
			},
			wantErr: true,
			errMsg:  "updateDefault can only be set for DateTime, Date, and Time",
		},
		{
			name: "field with UpdateDefault on DateTime (valid)",
			field: Field{
				Name:          "test",
				Caption:       "Test",
				Type:          string(FieldTypeDateTime),
				UpdateDefault: true,
			},
			wantErr: false,
		},
		{
			name: "field with min greater than max (Integer)",
			field: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeInteger),
				Min:     stringPtr("100"),
				Max:     stringPtr("50"),
			},
			wantErr: true,
			errMsg:  "must be less than max value",
		},
		{
			name: "field with min equal to max (Integer)",
			field: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeInteger),
				Min:     stringPtr("50"),
				Max:     stringPtr("50"),
			},
			wantErr: true,
			errMsg:  "must be less than max value",
		},
		{
			name: "field with invalid Boolean default value",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeBoolean),
				DefaultValue: stringPtr("maybe"),
			},
			wantErr: true,
			errMsg:  "default value for Boolean field must be 'true' or 'false'",
		},
		{
			name: "field with valid Boolean default value 'true'",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeBoolean),
				DefaultValue: stringPtr("true"),
			},
			wantErr: false,
		},
		{
			name: "field with valid Boolean default value 'false'",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeBoolean),
				DefaultValue: stringPtr("false"),
			},
			wantErr: false,
		},
		{
			name: "field with invalid Integer default value",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeInteger),
				DefaultValue: stringPtr("not-a-number"),
			},
			wantErr: true,
			errMsg:  "default value for Integer field must be a valid integer",
		},
		{
			name: "field with valid Integer default value",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeInteger),
				DefaultValue: stringPtr("42"),
			},
			wantErr: false,
		},
		{
			name: "field with invalid Float default value",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeFloat),
				DefaultValue: stringPtr("not-a-float"),
			},
			wantErr: true,
			errMsg:  "default value for Float field must be a valid number",
		},
		{
			name: "field with valid Float default value",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeFloat),
				DefaultValue: stringPtr("3.14"),
			},
			wantErr: false,
		},
		{
			name: "field with valid Decimal default value",
			field: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeDecimal),
				DefaultValue: stringPtr("99.99"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.field.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Field.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !contains(err.Error(), tt.errMsg) {
					t.Errorf("Field.Validate() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestField_Equals(t *testing.T) {
	field := Field{
		Name:    "testField",
		Caption: "Test Field",
		Type:    string(FieldTypeShortText),
	}

	tests := []struct {
		name        string
		uniqueInput FieldWhereUniqueInput
		want        bool
	}{
		{
			name: "match by caption",
			uniqueInput: FieldWhereUniqueInput{
				Caption: stringPtr("Test Field"),
			},
			want: true,
		},
		{
			name: "match by name",
			uniqueInput: FieldWhereUniqueInput{
				Name: stringPtr("testField"),
			},
			want: true,
		},
		{
			name: "match by both caption and name",
			uniqueInput: FieldWhereUniqueInput{
				Caption: stringPtr("Test Field"),
				Name:    stringPtr("testField"),
			},
			want: true,
		},
		{
			name: "no match - different caption",
			uniqueInput: FieldWhereUniqueInput{
				Caption: stringPtr("Different Field"),
			},
			want: false,
		},
		{
			name: "no match - different name",
			uniqueInput: FieldWhereUniqueInput{
				Name: stringPtr("differentField"),
			},
			want: false,
		},
		{
			name:        "no match - both nil",
			uniqueInput: FieldWhereUniqueInput{},
			want:        false,
		},
		{
			name: "no match - caption nil, name different",
			uniqueInput: FieldWhereUniqueInput{
				Name: stringPtr("differentField"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := field.Equals(tt.uniqueInput); got != tt.want {
				t.Errorf("Field.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_ApplyUpdateInput(t *testing.T) {
	tests := []struct {
		name     string
		field    *Field
		input    UpdateFieldInput
		expected Field
	}{
		{
			name: "update caption",
			field: &Field{
				Name:    "test",
				Caption: "Old Caption",
				Type:    string(FieldTypeShortText),
			},
			input: UpdateFieldInput{
				Caption: stringPtr("New Caption"),
			},
			expected: Field{
				Name:    "test",
				Caption: "New Caption",
				Type:    string(FieldTypeShortText),
			},
		},
		{
			name: "update required",
			field: &Field{
				Name:     "test",
				Caption:  "Test",
				Type:     string(FieldTypeShortText),
				Required: boolPtr(false),
			},
			input: UpdateFieldInput{
				Required: boolPtr(true),
			},
			expected: Field{
				Name:     "test",
				Caption:  "Test",
				Type:     string(FieldTypeShortText),
				Required: boolPtr(true),
			},
		},
		{
			name: "update unique",
			field: &Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeShortText),
				Unique:  boolPtr(false),
			},
			input: UpdateFieldInput{
				Unique: boolPtr(true),
			},
			expected: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeShortText),
				Unique:  boolPtr(true),
			},
		},
		{
			name: "update default value",
			field: &Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeShortText),
			},
			input: UpdateFieldInput{
				DefaultValue: stringPtr("new default"),
			},
			expected: Field{
				Name:         "test",
				Caption:      "Test",
				Type:         string(FieldTypeShortText),
				DefaultValue: stringPtr("new default"),
			},
		},
		{
			name: "update min",
			field: &Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeInteger),
			},
			input: UpdateFieldInput{
				Min: stringPtr("10"),
			},
			expected: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeInteger),
				Min:     stringPtr("10"),
			},
		},
		{
			name: "update max",
			field: &Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeInteger),
			},
			input: UpdateFieldInput{
				Max: stringPtr("100"),
			},
			expected: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeInteger),
				Max:     stringPtr("100"),
			},
		},
		{
			name: "update private",
			field: &Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeShortText),
				Private: boolPtr(false),
			},
			input: UpdateFieldInput{
				Private: boolPtr(true),
			},
			expected: Field{
				Name:    "test",
				Caption: "Test",
				Type:    string(FieldTypeShortText),
				Private: boolPtr(true),
			},
		},
		{
			name: "update accepted values",
			field: &Field{
				Name:           "test",
				Caption:        "Test",
				Type:           string(FieldTypeSingleChoice),
				AcceptedValues: []string{"old1", "old2"},
			},
			input: UpdateFieldInput{
				AcceptedValues: []string{"new1", "new2", "new3"},
			},
			expected: Field{
				Name:           "test",
				Caption:        "Test",
				Type:           string(FieldTypeSingleChoice),
				AcceptedValues: []string{"new1", "new2", "new3"},
			},
		},
		{
			name: "update multiple fields",
			field: &Field{
				Name:     "test",
				Caption:  "Old Caption",
				Type:     string(FieldTypeShortText),
				Required: boolPtr(false),
				Unique:   boolPtr(false),
			},
			input: UpdateFieldInput{
				Caption:  stringPtr("New Caption"),
				Required: boolPtr(true),
				Unique:   boolPtr(true),
				Min:      stringPtr("0"),
				Max:      stringPtr("100"),
			},
			expected: Field{
				Name:     "test",
				Caption:  "New Caption",
				Type:     string(FieldTypeShortText),
				Required: boolPtr(true),
				Unique:   boolPtr(true),
				Min:      stringPtr("0"),
				Max:      stringPtr("100"),
			},
		},
		{
			name: "update with nil values (should not change)",
			field: &Field{
				Name:     "test",
				Caption:  "Test",
				Type:     string(FieldTypeShortText),
				Required: boolPtr(true),
			},
			input: UpdateFieldInput{
				Caption:  nil,
				Required: nil,
			},
			expected: Field{
				Name:     "test",
				Caption:  "Test",
				Type:     string(FieldTypeShortText),
				Required: boolPtr(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.field.ApplyUpdateInput(tt.input)
			if !fieldsEqual(*tt.field, tt.expected) {
				t.Errorf("Field.ApplyUpdateInput() = %+v, want %+v", *tt.field, tt.expected)
			}
		})
	}
}

func TestValidateJSONDefaultValue(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid JSON object",
			value:   `{"key": "value"}`,
			wantErr: false,
		},
		{
			name:    "invalid JSON array (ValidateJSONDefaultValue only validates map[string]interface{})",
			value:   `["item1", "item2"]`,
			wantErr: true, // Arrays are not map[string]interface{}, so validation fails
		},
		{
			name:    "valid JSON with nested objects",
			value:   `{"nested": {"key": "value"}}`,
			wantErr: false,
		},
		{
			name:    "valid JSON with array of objects (but ValidateJSONDefaultValue only validates map[string]interface{})",
			value:   `[{"key": "value"}, {"key2": "value2"}]`,
			wantErr: true, // Arrays are not map[string]interface{}, so validation fails
		},
		{
			name:    "invalid JSON - missing closing brace",
			value:   `{"key": "value"`,
			wantErr: true,
		},
		{
			name:    "invalid JSON - missing closing bracket",
			value:   `["item1", "item2"`,
			wantErr: true,
		},
		{
			name:    "invalid JSON - malformed",
			value:   `{key: value}`,
			wantErr: true,
		},
		{
			name:    "invalid JSON - empty string",
			value:   ``,
			wantErr: true,
		},
		{
			name:    "invalid JSON - just text",
			value:   `not json`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJSONDefaultValue(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJSONDefaultValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJSONDefaultValueWithError(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid JSON object",
			value:   `{"key": "value"}`,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			value:   `{"key": "value"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := JSONDefaultValueWithError[map[string]interface{}](tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONDefaultValueWithError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// Test with array type
	t.Run("valid JSON array with []string type", func(t *testing.T) {
		value, err := JSONDefaultValueWithError[[]string](`["item1", "item2"]`)
		if err != nil {
			t.Errorf("JSONDefaultValueWithError() error = %v, wantErr false", err)
		}
		if len(value) != 2 {
			t.Errorf("JSONDefaultValueWithError() = %v, want length 2", value)
		}
	})

	// Test with array of objects type
	t.Run("valid JSON array of objects with []map[string]interface{} type", func(t *testing.T) {
		value, err := JSONDefaultValueWithError[[]map[string]interface{}](`[{"key": "value"}, {"key2": "value2"}]`)
		if err != nil {
			t.Errorf("JSONDefaultValueWithError() error = %v, wantErr false", err)
		}
		if len(value) != 2 {
			t.Errorf("JSONDefaultValueWithError() = %v, want length 2", value)
		}
	})
}

func TestJSONDefaultValue(t *testing.T) {
	t.Run("valid JSON object", func(t *testing.T) {
		value := JSONDefaultValue[map[string]interface{}](`{"key": "value"}`)
		if value == nil {
			t.Error("JSONDefaultValue() returned nil")
		}
		if val, ok := value["key"]; !ok || val != "value" {
			t.Errorf("JSONDefaultValue() = %v, want map with key 'value'", value)
		}
	})

	t.Run("invalid JSON (should not panic)", func(t *testing.T) {
		value := JSONDefaultValue[map[string]interface{}](`{"key": "value"`)
		// Should return zero value (nil for map) without panicking
		// This is expected behavior - invalid JSON returns zero value
		_ = value // Just verify it doesn't panic
	})

	t.Run("valid JSON array", func(t *testing.T) {
		value := JSONDefaultValue[[]string](`["item1", "item2"]`)
		if len(value) != 2 {
			t.Errorf("JSONDefaultValue() = %v, want length 2", value)
		}
	})
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsMiddle(s, substr))))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func fieldsEqual(f1, f2 Field) bool {
	if f1.Name != f2.Name || f1.Caption != f2.Caption || f1.Type != f2.Type {
		return false
	}
	if !boolPtrEqual(f1.Required, f2.Required) ||
		!boolPtrEqual(f1.Unique, f2.Unique) ||
		!stringPtrEqual(f1.DefaultValue, f2.DefaultValue) ||
		!stringPtrEqual(f1.Min, f2.Min) ||
		!stringPtrEqual(f1.Max, f2.Max) ||
		!boolPtrEqual(f1.Private, f2.Private) {
		return false
	}
	if len(f1.AcceptedValues) != len(f2.AcceptedValues) {
		return false
	}
	for i := range f1.AcceptedValues {
		if f1.AcceptedValues[i] != f2.AcceptedValues[i] {
			return false
		}
	}
	return true
}
