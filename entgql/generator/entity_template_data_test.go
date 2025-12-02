package generator

import (
	"testing"

	"github.com/GoLabra/labra/entgql/entity"
	"github.com/GoLabra/labra/strcase"
)

func TestEntityTemplateData_Validate(t *testing.T) {
	tests := []struct {
		name        string
		setupData   func() *EntityTemplateData
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid entity with required field",
			setupData: func() *EntityTemplateData {
				required := true
				return &EntityTemplateData{
					Entity: entity.Entity{
						Name:    "testEntity",
						EntName: "TestEntity",
						Caption: "Test Entity",
						Owner:   entity.EntityOwnerUser,
					},
					Fields: []entity.Field{
						{
							Name:     "name",
							Caption:  "Name",
							Type:     string(entity.FieldTypeShortText),
							Required: &required,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectError: false,
		},
		{
			name: "valid entity with unique field",
			setupData: func() *EntityTemplateData {
				unique := true
				return &EntityTemplateData{
					Entity: entity.Entity{
						Name:    "testEntity",
						EntName: "TestEntity",
						Caption: "Test Entity",
						Owner:   entity.EntityOwnerUser,
					},
					Fields: []entity.Field{
						{
							Name:    "email",
							Caption: "Email",
							Type:    string(entity.FieldTypeEmail),
							Unique:  &unique,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectError: false,
		},
		{
			name: "valid entity with default value",
			setupData: func() *EntityTemplateData {
				defaultVal := "default"
				return &EntityTemplateData{
					Entity: entity.Entity{
						Name:    "testEntity",
						EntName: "TestEntity",
						Caption: "Test Entity",
						Owner:   entity.EntityOwnerUser,
					},
					Fields: []entity.Field{
						{
							Name:         "status",
							Caption:      "Status",
							Type:         string(entity.FieldTypeShortText),
							DefaultValue: &defaultVal,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectError: false,
		},
		{
			name: "invalid: unique field with default value",
			setupData: func() *EntityTemplateData {
				unique := true
				defaultVal := "default"
				return &EntityTemplateData{
					Entity: entity.Entity{
						Name:    "testEntity",
						EntName: "TestEntity",
						Caption: "Test Entity",
						Owner:   entity.EntityOwnerUser,
					},
					Fields: []entity.Field{
						{
							Name:         "email",
							Caption:      "Email",
							Type:         string(entity.FieldTypeEmail),
							Unique:       &unique,
							DefaultValue: &defaultVal,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectError: true,
			errorMsg:    "cannot be unique and have default value set",
		},
		{
			name: "valid entity with min and max",
			setupData: func() *EntityTemplateData {
				minVal := "0"
				maxVal := "100"
				return &EntityTemplateData{
					Entity: entity.Entity{
						Name:    "testEntity",
						EntName: "TestEntity",
						Caption: "Test Entity",
						Owner:   entity.EntityOwnerUser,
					},
					Fields: []entity.Field{
						{
							Name:    "age",
							Caption: "Age",
							Type:    string(entity.FieldTypeInteger),
							Min:     &minVal,
							Max:     &maxVal,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectError: false,
		},
		{
			name: "valid entity with private field",
			setupData: func() *EntityTemplateData {
				private := true
				return &EntityTemplateData{
					Entity: entity.Entity{
						Name:    "testEntity",
						EntName: "TestEntity",
						Caption: "Test Entity",
						Owner:   entity.EntityOwnerUser,
					},
					Fields: []entity.Field{
						{
							Name:    "password",
							Caption: "Password",
							Type:    string(entity.FieldTypeShortText),
							Private: &private,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectError: false,
		},
		{
			name: "valid entity with required and unique",
			setupData: func() *EntityTemplateData {
				required := true
				unique := true
				return &EntityTemplateData{
					Entity: entity.Entity{
						Name:    "testEntity",
						EntName: "TestEntity",
						Caption: "Test Entity",
						Owner:   entity.EntityOwnerUser,
					},
					Fields: []entity.Field{
						{
							Name:     "email",
							Caption:  "Email",
							Type:     string(entity.FieldTypeEmail),
							Required: &required,
							Unique:   &unique,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectError: false,
		},
		{
			name: "valid entity with required and default",
			setupData: func() *EntityTemplateData {
				required := true
				defaultVal := "default"
				return &EntityTemplateData{
					Entity: entity.Entity{
						Name:    "testEntity",
						EntName: "TestEntity",
						Caption: "Test Entity",
						Owner:   entity.EntityOwnerUser,
					},
					Fields: []entity.Field{
						{
							Name:         "status",
							Caption:      "Status",
							Type:         string(entity.FieldTypeShortText),
							Required:     &required,
							DefaultValue: &defaultVal,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectError: false,
		},
		{
			name: "valid entity with nillable field",
			setupData: func() *EntityTemplateData {
				return &EntityTemplateData{
					Entity: entity.Entity{
						Name:    "testEntity",
						EntName: "TestEntity",
						Caption: "Test Entity",
						Owner:   entity.EntityOwnerUser,
					},
					Fields: []entity.Field{
						{
							Name:     "description",
							Caption:  "Description",
							Type:     string(entity.FieldTypeLongText),
							Nillable: true,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectError: false,
		},
		{
			name: "valid entity with UpdateDefault",
			setupData: func() *EntityTemplateData {
				defaultVal := "now()"
				return &EntityTemplateData{
					Entity: entity.Entity{
						Name:    "testEntity",
						EntName: "TestEntity",
						Caption: "Test Entity",
						Owner:   entity.EntityOwnerUser,
					},
					Fields: []entity.Field{
						{
							Name:          "updatedAt",
							Caption:       "Updated At",
							Type:          string(entity.FieldTypeDateTime),
							DefaultValue:  &defaultVal,
							UpdateDefault: true,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectError: false,
		},
		{
			name: "valid entity with accepted values",
			setupData: func() *EntityTemplateData {
				return &EntityTemplateData{
					Entity: entity.Entity{
						Name:    "testEntity",
						EntName: "TestEntity",
						Caption: "Test Entity",
						Owner:   entity.EntityOwnerUser,
					},
					Fields: []entity.Field{
						{
							Name:           "status",
							Caption:        "Status",
							Type:           string(entity.FieldTypeSingleChoice),
							AcceptedValues: []string{"active", "inactive", "pending"},
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectError: false,
		},
		{
			name: "valid entity with JSON field and valid default",
			setupData: func() *EntityTemplateData {
				defaultVal := `{"key": "value"}`
				return &EntityTemplateData{
					Entity: entity.Entity{
						Name:    "testEntity",
						EntName: "TestEntity",
						Caption: "Test Entity",
						Owner:   entity.EntityOwnerUser,
					},
					Fields: []entity.Field{
						{
							Name:         "metadata",
							Caption:      "Metadata",
							Type:         string(entity.FieldTypeJson),
							DefaultValue: &defaultVal,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectError: false,
		},
		{
			name: "invalid: JSON field with invalid default",
			setupData: func() *EntityTemplateData {
				defaultVal := `{"key": "value"` // Invalid JSON
				return &EntityTemplateData{
					Entity: entity.Entity{
						Name:    "testEntity",
						EntName: "TestEntity",
						Caption: "Test Entity",
						Owner:   entity.EntityOwnerUser,
					},
					Fields: []entity.Field{
						{
							Name:         "metadata",
							Caption:      "Metadata",
							Type:         string(entity.FieldTypeJson),
							DefaultValue: &defaultVal,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectError: true,
			errorMsg:    "json format invalid",
		},
		{
			name: "valid entity with multiple field combinations",
			setupData: func() *EntityTemplateData {
				required := true
				private := true
				minVal := "1"
				maxVal := "100"
				return &EntityTemplateData{
					Entity: entity.Entity{
						Name:    "testEntity",
						EntName: "TestEntity",
						Caption: "Test Entity",
						Owner:   entity.EntityOwnerUser,
					},
					Fields: []entity.Field{
						{
							Name:     "name",
							Caption:  "Name",
							Type:     string(entity.FieldTypeShortText),
							Required: &required,
						},
						{
							Name:    "age",
							Caption: "Age",
							Type:    string(entity.FieldTypeInteger),
							Min:     &minVal,
							Max:     &maxVal,
						},
						{
							Name:    "password",
							Caption: "Password",
							Type:    string(entity.FieldTypeShortText),
							Private: &private,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectError: false,
		},
		{
			name: "valid entity with all field types",
			setupData: func() *EntityTemplateData {
				required := true
				unique := true
				defaultVal := "default"
				private := true
				minVal := "0"
				maxVal := "100"
				return &EntityTemplateData{
					Entity: entity.Entity{
						Name:    "testEntity",
						EntName: "TestEntity",
						Caption: "Test Entity",
						Owner:   entity.EntityOwnerUser,
					},
					Fields: []entity.Field{
						{Name: strcase.ToLowerCamel("ShortText"), Type: string(entity.FieldTypeShortText), Caption: "ShortText"},
						{Name: strcase.ToLowerCamel("LongText"), Type: string(entity.FieldTypeLongText), Caption: "LongText"},
						{Name: strcase.ToLowerCamel("RichText"), Type: string(entity.FieldTypeRichText), Caption: "RichText"},
						{Name: strcase.ToLowerCamel("Email"), Type: string(entity.FieldTypeEmail), Caption: "Email"},
						{Name: strcase.ToLowerCamel("Integer"), Type: string(entity.FieldTypeInteger), Caption: "Integer"},
						{Name: strcase.ToLowerCamel("Decimal"), Type: string(entity.FieldTypeDecimal), Caption: "Decimal"},
						{Name: strcase.ToLowerCamel("FloatValue"), Type: string(entity.FieldTypeFloat), Caption: "FloatValue"},
						{Name: strcase.ToLowerCamel("Boolean"), Type: string(entity.FieldTypeBoolean), Caption: "Boolean"},
						{Name: strcase.ToLowerCamel("SingleChoice"), Type: string(entity.FieldTypeSingleChoice), Caption: "SingleChoice", AcceptedValues: []string{"option1", "option2"}},
						{Name: strcase.ToLowerCamel("MultipleChoice"), Type: string(entity.FieldTypeMultipleChoice), Caption: "MultipleChoice", AcceptedValues: []string{"option1", "option2"}},
						{Name: strcase.ToLowerCamel("DateTime"), Type: string(entity.FieldTypeDateTime), Caption: "DateTime"},
						{Name: strcase.ToLowerCamel("Date"), Type: string(entity.FieldTypeDate), Caption: "Date"},
						{Name: strcase.ToLowerCamel("Time"), Type: string(entity.FieldTypeTime), Caption: "Time"},
						{Name: strcase.ToLowerCamel("Json"), Type: string(entity.FieldTypeJson), Caption: "Json"},
						{Name: strcase.ToLowerCamel("Enum"), Type: string(entity.FieldTypeEnum), Caption: "Enum"},
						{Name: strcase.ToLowerCamel("Enums"), Type: string(entity.FieldTypeEnums), Caption: "Enums"},
						{Name: strcase.ToLowerCamel("Required"), Type: string(entity.FieldTypeShortText), Caption: "Required", Required: &required},
						{Name: strcase.ToLowerCamel("Unique"), Type: string(entity.FieldTypeEmail), Caption: "Unique", Unique: &unique},
						{Name: strcase.ToLowerCamel("DefaultValue"), Type: string(entity.FieldTypeShortText), Caption: "DefaultValue", DefaultValue: &defaultVal},
						{Name: strcase.ToLowerCamel("Private"), Type: string(entity.FieldTypeShortText), Caption: "Private", Private: &private},
						{Name: strcase.ToLowerCamel("MinMax"), Type: string(entity.FieldTypeInteger), Caption: "MinMax", Min: &minVal, Max: &maxVal},
					},
					Edges: []entity.Edge{},
				}
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.setupData()
			err := data.Validate()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

