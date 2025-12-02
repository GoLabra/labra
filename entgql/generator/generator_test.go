package generator

import (
	"os"
	"strings"
	"testing"

	"github.com/GoLabra/labra/entgql/entity"
	"github.com/GoLabra/labra/subscription"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// MockFileSystem implements FileSystem interface for testing
type MockFileSystem struct {
	files       map[string]string // filename -> temp file path
	createCalls []string
	openCalls   []string
	removeCalls []string
}

func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		files:       make(map[string]string),
		createCalls: []string{},
		openCalls:   []string{},
		removeCalls: []string{},
	}
}

func (m *MockFileSystem) Create(name string) (*os.File, error) {
	m.createCalls = append(m.createCalls, name)
	// Create a temporary file
	file, err := os.CreateTemp("", "test-*.go")
	if err != nil {
		return nil, err
	}
	// Store the temp file path so we can read it later
	m.files[name] = file.Name()
	return file, nil
}

func (m *MockFileSystem) Open(name string) (*os.File, error) {
	m.openCalls = append(m.openCalls, name)
	return os.Open(name)
}

func (m *MockFileSystem) CopyDirectory(from, to string) error {
	return nil
}

func (m *MockFileSystem) Remove(path string) error {
	m.removeCalls = append(m.removeCalls, path)
	if tempPath, ok := m.files[path]; ok {
		os.Remove(tempPath)
		delete(m.files, path)
	}
	return nil
}

func (m *MockFileSystem) Exit(code int) {
	// No-op for testing
}

func (m *MockFileSystem) GetFileContent(name string) (string, error) {
	if tempPath, ok := m.files[name]; ok {
		content, err := os.ReadFile(tempPath)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}
	return "", os.ErrNotExist
}

// MockSubscriptionClient implements SubscriptionClient interface
type MockSubscriptionClient struct{}

func (m *MockSubscriptionClient) PublishAppStatusMessage(appStatus subscription.AppStatus) {
	// No-op for testing
}

func TestImports(t *testing.T) {
	tests := []struct {
		name            string
		fields          []entity.Field
		expectedImports map[string]bool
	}{
		{
			name:            "no fields",
			fields:          []entity.Field{},
			expectedImports: map[string]bool{},
		},
		{
			name: "ShortText field",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeShortText), Caption: "Name"},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect": true,
			},
		},
		{
			name: "LongText field",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeLongText), Caption: "Description"},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect": true,
			},
		},
		{
			name: "RichText field",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeRichText), Caption: "Content"},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect": true,
			},
		},
		{
			name: "Email field",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeEmail), Caption: "Email"},
			},
			expectedImports: map[string]bool{},
		},
		{
			name: "Integer field",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeInteger), Caption: "Age"},
			},
			expectedImports: map[string]bool{},
		},
		{
			name: "Decimal field",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeDecimal), Caption: "Price"},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect": true,
			},
		},
		{
			name: "Float field",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeFloat), Caption: "Rating"},
			},
			expectedImports: map[string]bool{},
		},
		{
			name: "Boolean field",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeBoolean), Caption: "Active"},
			},
			expectedImports: map[string]bool{},
		},
		{
			name: "SingleChoice field",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeSingleChoice), Caption: "Status"},
			},
			expectedImports: map[string]bool{
				"fmt": true,
			},
		},
		{
			name: "MultipleChoice field",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeMultipleChoice), Caption: "Tags"},
			},
			expectedImports: map[string]bool{
				"fmt": true,
			},
		},
		{
			name: "DateTime field without default",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeDateTime), Caption: "Created At"},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect": true,
			},
		},
		{
			name: "DateTime field with default value",
			fields: []entity.Field{
				{
					Type:         string(entity.FieldTypeDateTime),
					Caption:      "Created At",
					DefaultValue: stringPtr("now()"),
				},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect":                 true,
				"time":                                 true,
				"github.com/GoLabra/labra/entgql/date": true,
			},
		},
		{
			name: "DateTime field with UpdateDefault",
			fields: []entity.Field{
				{
					Type:          string(entity.FieldTypeDateTime),
					Caption:       "Updated At",
					UpdateDefault: true,
				},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect": true,
				"time":                 true,
			},
		},
		{
			name: "DateTime field with default and UpdateDefault",
			fields: []entity.Field{
				{
					Type:          string(entity.FieldTypeDateTime),
					Caption:       "Updated At",
					DefaultValue:  stringPtr("now()"),
					UpdateDefault: true,
				},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect":                 true,
				"time":                                 true,
				"github.com/GoLabra/labra/entgql/date": true,
			},
		},
		{
			name: "Date field without default",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeDate), Caption: "Birth Date"},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect": true,
			},
		},
		{
			name: "Date field with default value",
			fields: []entity.Field{
				{
					Type:         string(entity.FieldTypeDate),
					Caption:      "Birth Date",
					DefaultValue: stringPtr("now()"),
				},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect":                 true,
				"time":                                 true,
				"github.com/GoLabra/labra/entgql/date": true,
			},
		},
		{
			name: "Time field without default",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeTime), Caption: "Start Time"},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect": true,
			},
		},
		{
			name: "Time field with default value",
			fields: []entity.Field{
				{
					Type:         string(entity.FieldTypeTime),
					Caption:      "Start Time",
					DefaultValue: stringPtr("now()"),
				},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect":                 true,
				"time":                                 true,
				"github.com/GoLabra/labra/entgql/date": true,
			},
		},
		{
			name: "Json field",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeJson), Caption: "Metadata"},
			},
			expectedImports: map[string]bool{},
		},
		{
			name: "Enum field",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeEnum), Caption: "Status"},
			},
			expectedImports: map[string]bool{},
		},
		{
			name: "Enums field",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeEnums), Caption: "Tags"},
			},
			expectedImports: map[string]bool{},
		},
		{
			name: "multiple fields with different types",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeShortText), Caption: "Name"},
				{Type: string(entity.FieldTypeEmail), Caption: "Email"},
				{Type: string(entity.FieldTypeInteger), Caption: "Age"},
				{Type: string(entity.FieldTypeBoolean), Caption: "Active"},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect": true,
			},
		},
		{
			name: "multiple fields requiring fmt import",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeSingleChoice), Caption: "Status"},
				{Type: string(entity.FieldTypeMultipleChoice), Caption: "Tags"},
			},
			expectedImports: map[string]bool{
				"fmt": true,
			},
		},
		{
			name: "DateTime and Date fields with defaults",
			fields: []entity.Field{
				{
					Type:         string(entity.FieldTypeDateTime),
					Caption:      "Created At",
					DefaultValue: stringPtr("now()"),
				},
				{
					Type:         string(entity.FieldTypeDate),
					Caption:      "Birth Date",
					DefaultValue: stringPtr("now()"),
				},
				{
					Type:          string(entity.FieldTypeDateTime),
					Caption:       "Updated At",
					UpdateDefault: true,
				},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect":                 true,
				"time":                                 true,
				"github.com/GoLabra/labra/entgql/date": true,
			},
		},
		{
			name: "all text field types",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeShortText), Caption: "Short"},
				{Type: string(entity.FieldTypeLongText), Caption: "Long"},
				{Type: string(entity.FieldTypeRichText), Caption: "Rich"},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect": true,
			},
		},
		{
			name: "numeric field types",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeInteger), Caption: "Int"},
				{Type: string(entity.FieldTypeDecimal), Caption: "Decimal"},
				{Type: string(entity.FieldTypeFloat), Caption: "Float"},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect": true,
			},
		},
		{
			name: "choice field types",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeSingleChoice), Caption: "Single"},
				{Type: string(entity.FieldTypeMultipleChoice), Caption: "Multiple"},
			},
			expectedImports: map[string]bool{
				"fmt": true,
			},
		},
		{
			name: "date/time field types",
			fields: []entity.Field{
				{Type: string(entity.FieldTypeDateTime), Caption: "DateTime"},
				{Type: string(entity.FieldTypeDate), Caption: "Date"},
				{Type: string(entity.FieldTypeTime), Caption: "Time"},
			},
			expectedImports: map[string]bool{
				"entgo.io/ent/dialect": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Imports(tt.fields)

			// Check that all expected imports are present
			for expectedImport := range tt.expectedImports {
				if !result[expectedImport] {
					t.Errorf("expected import %s to be present, but it was not", expectedImport)
				}
			}

			// Check that no unexpected imports are present
			for actualImport := range result {
				if !tt.expectedImports[actualImport] {
					t.Errorf("unexpected import %s found", actualImport)
				}
			}

			// Verify the count matches
			if len(result) != len(tt.expectedImports) {
				t.Errorf("expected %d imports, got %d", len(tt.expectedImports), len(result))
			}
		})
	}
}

func TestWriteEntityToSchema(t *testing.T) {
	tests := []struct {
		name           string
		setupData      func() EntityTemplateData
		expectedInFile []string // Substrings that should be in the generated file
		notInFile      []string // Substrings that should NOT be in the generated file
	}{
		{
			name: "simple entity with ShortText field",
			setupData: func() EntityTemplateData {
				required := true
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "testEntity",
						EntName:          "TestEntity",
						Caption:          "Test Entity",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "name",
					},
					Fields: []entity.Field{
						{
							Name:     "name",
							EntName:  "name",
							Caption:  "Name",
							Type:     string(entity.FieldTypeShortText),
							Required: &required,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				"package schema",
				"type TestEntity struct",
				"ent.Schema",
				`Caption: "Test Entity"`,
				`DisplayField: "name"`,
				`field.String("name")`,
				`Caption:   "Name"`,
				`Type: entity.FieldTypeShortText`,
				`SchemaType(map[string]string{`,
				`dialect.MySQL:    "VARCHAR(255)"`,
				`dialect.Postgres: "VARCHAR(255)"`,
			},
			notInFile: []string{
				"Optional()",
				"Nillable()",
			},
		},
		{
			name: "entity with required and unique Email field",
			setupData: func() EntityTemplateData {
				required := true
				unique := true
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "user",
						EntName:          "User",
						Caption:          "User",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "email",
					},
					Fields: []entity.Field{
						{
							Name:     "email",
							EntName:  "email",
							Caption:  "Email",
							Type:     string(entity.FieldTypeEmail),
							Required: &required,
							Unique:   &unique,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.String("email")`,
				`Unique()`,
				`Validate(entity.Validate`,
				`Caption:   "Email"`,
				`Type: entity.FieldTypeEmail`,
			},
			notInFile: []string{
				"Optional()",
				"Nillable()",
			},
		},
		{
			name: "entity with optional field with default value",
			setupData: func() EntityTemplateData {
				defaultVal := "active"
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "status",
						EntName:          "Status",
						Caption:          "Status",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "name",
					},
					Fields: []entity.Field{
						{
							Name:         "name",
							EntName:      "name",
							Caption:      "Name",
							Type:         string(entity.FieldTypeShortText),
							DefaultValue: &defaultVal,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`Default("active")`,
				"Optional()",
				"Nillable()",
			},
		},
		{
			name: "entity with Integer field with min and max",
			setupData: func() EntityTemplateData {
				minVal := "0"
				maxVal := "100"
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "score",
						EntName:          "Score",
						Caption:          "Score",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "value",
					},
					Fields: []entity.Field{
						{
							Name:    "value",
							EntName: "value",
							Caption: "Value",
							Type:    string(entity.FieldTypeInteger),
							Min:     &minVal,
							Max:     &maxVal,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.Int("value")`,
				`Min(0)`,
				`Max(100)`,
				`Min:     "0"`,
				`Max:     "100"`,
				`Caption: "Value"`,
				`Type:    entity.FieldTypeInteger`,
			},
		},
		{
			name: "entity with DateTime field with default and UpdateDefault",
			setupData: func() EntityTemplateData {
				defaultVal := "now()"
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "post",
						EntName:          "Post",
						Caption:          "Post",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "title",
					},
					Fields: []entity.Field{
						{
							Name:    "title",
							EntName: "title",
							Caption: "Title",
							Type:    string(entity.FieldTypeShortText),
						},
						{
							Name:          "updatedAt",
							EntName:       "updated_at",
							Caption:       "Updated At",
							Type:          string(entity.FieldTypeDateTime),
							DefaultValue:  &defaultVal,
							UpdateDefault: true,
							Nillable:      true,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.Time("updated_at")`,
				`Nillable()`,
				`UpdateDefault(time.Now)`,
				`Default(func(val string) time.Time`,
				`Type: entity.FieldTypeDateTime`,
				`DefaultValue: "now()"`,
				"time",
				"github.com/GoLabra/labra/entgql/date",
			},
		},
		{
			name: "entity with private field",
			setupData: func() EntityTemplateData {
				private := true
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "user",
						EntName:          "User",
						Caption:          "User",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "email",
					},
					Fields: []entity.Field{
						{
							Name:    "password",
							EntName: "password",
							Caption: "Password",
							Type:    string(entity.FieldTypeShortText),
							Private: &private,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`Private: true`,
			},
		},
		{
			name: "entity with SingleChoice field",
			setupData: func() EntityTemplateData {
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "status",
						EntName:          "Status",
						Caption:          "Status",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "name",
					},
					Fields: []entity.Field{
						{
							Name:           "name",
							EntName:        "name",
							Caption:        "Name",
							Type:           string(entity.FieldTypeSingleChoice),
							AcceptedValues: []string{"active", "inactive", "pending"},
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`Type: entity.FieldTypeSingleChoice`,
				`AcceptedValues: []string{`,
				`"active"`,
				`"inactive"`,
				`"pending"`,
				`fmt.Errorf`,
				"fmt",
			},
		},
		{
			name: "entity with MultipleChoice field",
			setupData: func() EntityTemplateData {
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "tag",
						EntName:          "Tag",
						Caption:          "Tag",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "tags",
					},
					Fields: []entity.Field{
						{
							Name:           "tags",
							EntName:        "tags",
							Caption:        "Tags",
							Type:           string(entity.FieldTypeMultipleChoice),
							AcceptedValues: []string{"tag1", "tag2", "tag3"},
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.Strings("tags")`,
				`Type: entity.FieldTypeMultipleChoice`,
				`AcceptedValues: []string{`,
				`"tag1"`,
				`"tag2"`,
				`"tag3"`,
				`fmt.Errorf`,
				"fmt",
			},
		},
		{
			name: "entity with LongText field",
			setupData: func() EntityTemplateData {
				required := true
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "article",
						EntName:          "Article",
						Caption:          "Article",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "title",
					},
					Fields: []entity.Field{
						{
							Name:     "title",
							EntName:  "title",
							Caption:  "Title",
							Type:     string(entity.FieldTypeShortText),
							Required: &required,
						},
						{
							Name:     "content",
							EntName:  "content",
							Caption:  "Content",
							Type:     string(entity.FieldTypeLongText),
							Required: &required,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.String("content")`,
				`Type: entity.FieldTypeLongText`,
				`dialect.MySQL:    "MEDIUMTEXT"`,
				`dialect.Postgres: "TEXT"`,
			},
			notInFile: []string{
				"Optional()",
				"Nillable()",
			},
		},
		{
			name: "entity with RichText field",
			setupData: func() EntityTemplateData {
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "page",
						EntName:          "Page",
						Caption:          "Page",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "body",
					},
					Fields: []entity.Field{
						{
							Name:    "body",
							EntName: "body",
							Caption: "Body",
							Type:    string(entity.FieldTypeRichText),
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.String("body")`,
				`Type: entity.FieldTypeRichText`,
				`dialect.MySQL:    "MEDIUMTEXT"`,
				`dialect.Postgres: "TEXT"`,
				"Optional()",
				"Nillable()",
			},
		},
		{
			name: "entity with Boolean field",
			setupData: func() EntityTemplateData {
				required := true
				defaultVal := "true"
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "settings",
						EntName:          "Settings",
						Caption:          "Settings",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "active",
					},
					Fields: []entity.Field{
						{
							Name:         "active",
							EntName:      "active",
							Caption:      "Active",
							Type:         string(entity.FieldTypeBoolean),
							Required:     &required,
							DefaultValue: &defaultVal,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.Bool("active")`,
				`Default(true)`,
				`Type: entity.FieldTypeBoolean`,
			},
			notInFile: []string{
				"Optional()",
				"Nillable()",
			},
		},
		{
			name: "entity with Float field with min and max",
			setupData: func() EntityTemplateData {
				minVal := "0.0"
				maxVal := "10.0"
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "rating",
						EntName:          "Rating",
						Caption:          "Rating",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "value",
					},
					Fields: []entity.Field{
						{
							Name:    "value",
							EntName: "value",
							Caption: "Value",
							Type:    string(entity.FieldTypeFloat),
							Min:     &minVal,
							Max:     &maxVal,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.Float("value")`,
				`Min(0.0)`,
				`Max(10.0)`,
				`Min:     "0.0"`,
				`Max:     "10.0"`,
				`Type:    entity.FieldTypeFloat`,
			},
		},
		{
			name: "entity with Decimal field",
			setupData: func() EntityTemplateData {
				required := true
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "price",
						EntName:          "Price",
						Caption:          "Price",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "amount",
					},
					Fields: []entity.Field{
						{
							Name:     "amount",
							EntName:  "amount",
							Caption:  "Amount",
							Type:     string(entity.FieldTypeDecimal),
							Required: &required,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.Float("amount")`,
				`Type:    entity.FieldTypeFloat`,
				`dialect.MySQL:    "decimal(6, 4)"`,
				`dialect.Postgres: "decimal(6, 4)"`,
			},
			notInFile: []string{
				"Optional()",
				"Nillable()",
			},
		},
		{
			name: "entity with Time field with default",
			setupData: func() EntityTemplateData {
				defaultVal := "now"
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "schedule",
						EntName:          "Schedule",
						Caption:          "Schedule",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "startTime",
					},
					Fields: []entity.Field{
						{
							Name:         "startTime",
							EntName:      "start_time",
							Caption:      "Start Time",
							Type:         string(entity.FieldTypeTime),
							DefaultValue: &defaultVal,
							Nillable:     true,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.Time("start_time")`,
				`Type: entity.FieldTypeTime`,
				`entgql.Type("TimeOnly")`,
				`dialect.MySQL:    "time"`,
				`dialect.Postgres: "time"`,
				"time",
				"github.com/GoLabra/labra/entgql/date",
			},
		},
		{
			name: "entity with Date field with default",
			setupData: func() EntityTemplateData {
				defaultVal := "now()"
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "event",
						EntName:          "Event",
						Caption:          "Event",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "name",
					},
					Fields: []entity.Field{
						{
							Name:         "eventDate",
							EntName:      "event_date",
							Caption:      "Event Date",
							Type:         string(entity.FieldTypeDate),
							DefaultValue: &defaultVal,
							Nillable:     true,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.Time("event_date")`,
				`Type: entity.FieldTypeDate`,
				`entgql.Type("DateOnly")`,
				`dialect.MySQL:    "date"`,
				`dialect.Postgres: "date"`,
				"time",
				"github.com/GoLabra/labra/entgql/date",
			},
		},
		{
			name: "entity with Json field with default",
			setupData: func() EntityTemplateData {
				defaultVal := `{"key": "value"}`
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "config",
						EntName:          "Config",
						Caption:          "Config",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "name",
					},
					Fields: []entity.Field{
						{
							Name:         "metadata",
							EntName:      "metadata",
							Caption:      "Metadata",
							Type:         string(entity.FieldTypeJson),
							DefaultValue: &defaultVal,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.JSON("metadata"`,
				`Type: entity.FieldTypeJson`,
				`entgql.Type("Map")`,
				`DefaultValue:`,
			},
		},
		{
			name: "entity with required unique field (no default)",
			setupData: func() EntityTemplateData {
				required := true
				unique := true
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "product",
						EntName:          "Product",
						Caption:          "Product",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "sku",
					},
					Fields: []entity.Field{
						{
							Name:     "sku",
							EntName:  "sku",
							Caption:  "SKU",
							Type:     string(entity.FieldTypeShortText),
							Required: &required,
							Unique:   &unique,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.String("sku")`,
				`Unique()`,
			},
			notInFile: []string{
				"Optional()",
				"Nillable()",
				"Default(",
			},
		},
		{
			name: "entity with nillable optional field",
			setupData: func() EntityTemplateData {
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "profile",
						EntName:          "Profile",
						Caption:          "Profile",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "name",
					},
					Fields: []entity.Field{
						{
							Name:    "name",
							EntName: "name",
							Caption: "Name",
							Type:    string(entity.FieldTypeShortText),
						},
						{
							Name:     "bio",
							EntName:  "bio",
							Caption:  "Bio",
							Type:     string(entity.FieldTypeLongText),
							Nillable: true,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				"Optional()",
				"Nillable()",
			},
		},
		{
			name: "entity with Integer field with default value",
			setupData: func() EntityTemplateData {
				defaultVal := "0"
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "counter",
						EntName:          "Counter",
						Caption:          "Counter",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "value",
					},
					Fields: []entity.Field{
						{
							Name:         "value",
							EntName:      "value",
							Caption:      "Value",
							Type:         string(entity.FieldTypeInteger),
							DefaultValue: &defaultVal,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.Int("value")`,
				`Default(0)`,
			},
		},
		{
			name: "entity with Float field with default value",
			setupData: func() EntityTemplateData {
				defaultVal := "0.0"
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "measurement",
						EntName:          "Measurement",
						Caption:          "Measurement",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "value",
					},
					Fields: []entity.Field{
						{
							Name:         "value",
							EntName:      "value",
							Caption:      "Value",
							Type:         string(entity.FieldTypeFloat),
							DefaultValue: &defaultVal,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.Float("value")`,
				`Default(0.0)`,
			},
		},
		{
			name: "entity with SingleChoice field with default and accepted values",
			setupData: func() EntityTemplateData {
				defaultVal := "pending"
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "order",
						EntName:          "Order",
						Caption:          "Order",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "status",
					},
					Fields: []entity.Field{
						{
							Name:           "status",
							EntName:        "status",
							Caption:        "Status",
							Type:           string(entity.FieldTypeSingleChoice),
							DefaultValue:   &defaultVal,
							AcceptedValues: []string{"pending", "processing", "completed", "cancelled"},
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.String("status")`,
				`Default("pending")`,
				`AcceptedValues: []string{`,
				`"pending"`,
				`"processing"`,
				`"completed"`,
				`"cancelled"`,
			},
		},
		{
			name: "entity with MultipleChoice field with default",
			setupData: func() EntityTemplateData {
				defaultVal := `["tag1", "tag2"]`
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "post",
						EntName:          "Post",
						Caption:          "Post",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "tags",
					},
					Fields: []entity.Field{
						{
							Name:           "tags",
							EntName:        "tags",
							Caption:        "Tags",
							Type:           string(entity.FieldTypeMultipleChoice),
							DefaultValue:   &defaultVal,
							AcceptedValues: []string{"tag1", "tag2", "tag3"},
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.Strings("tags")`,
				`Default(entity.JSONDefaultValue`,
				`AcceptedValues: []string{`,
			},
		},
		{
			name: "entity with DateTime field without default but with UpdateDefault",
			setupData: func() EntityTemplateData {
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "document",
						EntName:          "Document",
						Caption:          "Document",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "name",
					},
					Fields: []entity.Field{
						{
							Name:          "modifiedAt",
							EntName:       "modified_at",
							Caption:       "Modified At",
							Type:          string(entity.FieldTypeDateTime),
							UpdateDefault: true,
							Nillable:      true,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.Time("modified_at")`,
				`UpdateDefault(time.Now)`,
				`Nillable()`,
				"time",
			},
			notInFile: []string{
				`Default(func(val string)`,
				`github.com/GoLabra/labra/entgql/date`,
			},
		},
		{
			name: "entity with all field types comprehensive",
			setupData: func() EntityTemplateData {
				required := true
				unique := true
				private := true
				minVal := "0"
				maxVal := "100"
				defaultVal := "default"
				defaultFloat := "3.14"
				defaultBool := "true"
				defaultDate := "now()"
				defaultJson := `{"key": "value"}`
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "comprehensive",
						EntName:          "Comprehensive",
						Caption:          "Comprehensive",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "name",
					},
					Fields: []entity.Field{
						{Name: "name", EntName: "name", Caption: "Name", Type: string(entity.FieldTypeShortText), Required: &required},
						{Name: "description", EntName: "description", Caption: "Description", Type: string(entity.FieldTypeLongText)},
						{Name: "content", EntName: "content", Caption: "Content", Type: string(entity.FieldTypeRichText)},
						{Name: "email", EntName: "email", Caption: "Email", Type: string(entity.FieldTypeEmail), Unique: &unique},
						{Name: "age", EntName: "age", Caption: "Age", Type: string(entity.FieldTypeInteger), Min: &minVal, Max: &maxVal},
						{Name: "price", EntName: "price", Caption: "Price", Type: string(entity.FieldTypeDecimal)},
						{Name: "rating", EntName: "rating", Caption: "Rating", Type: string(entity.FieldTypeFloat), DefaultValue: &defaultFloat},
						{Name: "active", EntName: "active", Caption: "Active", Type: string(entity.FieldTypeBoolean), DefaultValue: &defaultBool},
						{Name: "status", EntName: "status", Caption: "Status", Type: string(entity.FieldTypeSingleChoice), AcceptedValues: []string{"a", "b", "c"}},
						{Name: "tags", EntName: "tags", Caption: "Tags", Type: string(entity.FieldTypeMultipleChoice), AcceptedValues: []string{"x", "y", "z"}},
						{Name: "createdAt", EntName: "created_at", Caption: "Created At", Type: string(entity.FieldTypeDateTime), DefaultValue: &defaultDate, Nillable: true},
						{Name: "birthDate", EntName: "birth_date", Caption: "Birth Date", Type: string(entity.FieldTypeDate), DefaultValue: &defaultDate},
						{Name: "startTime", EntName: "start_time", Caption: "Start Time", Type: string(entity.FieldTypeTime), DefaultValue: &defaultDate},
						{Name: "metadata", EntName: "metadata", Caption: "Metadata", Type: string(entity.FieldTypeJson), DefaultValue: &defaultJson},
						{Name: "secret", EntName: "secret", Caption: "Secret", Type: string(entity.FieldTypeShortText), Private: &private},
						{Name: "defaultStatus", EntName: "default_status", Caption: "Default Status", Type: string(entity.FieldTypeShortText), DefaultValue: &defaultVal},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.String("name")`,
				`field.String("description")`,
				`field.String("content")`,
				`field.String("email")`,
				`field.Int("age")`,
				`field.Float("price")`,
				`field.Float("rating")`,
				`field.Bool("active")`,
				`field.String("status")`,
				`field.Strings("tags")`,
				`field.Time("created_at")`,
				`field.Time("birth_date")`,
				`field.Time("start_time")`,
				`field.JSON("metadata"`,
				`Private: true`,
				`Default("default")`,
			},
		},
		{
			name: "entity with required field combinations",
			setupData: func() EntityTemplateData {
				required := true
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "required",
						EntName:          "Required",
						Caption:          "Required",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "field1",
					},
					Fields: []entity.Field{
						{Name: "field1", EntName: "field1", Caption: "Field1", Type: string(entity.FieldTypeShortText), Required: &required},
						{Name: "field2", EntName: "field2", Caption: "Field2", Type: string(entity.FieldTypeInteger), Required: &required},
						{Name: "field3", EntName: "field3", Caption: "Field3", Type: string(entity.FieldTypeBoolean), Required: &required},
						{Name: "field4", EntName: "field4", Caption: "Field4", Type: string(entity.FieldTypeEmail), Required: &required},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.String("field1")`,
				`field.Int("field2")`,
				`field.Bool("field3")`,
				`field.String("field4")`,
			},
			notInFile: []string{
				"Optional()",
				"Nillable()",
			},
		},
		{
			name: "entity with unique field combinations",
			setupData: func() EntityTemplateData {
				unique := true
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "unique",
						EntName:          "Unique",
						Caption:          "Unique",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "field1",
					},
					Fields: []entity.Field{
						{Name: "field1", EntName: "field1", Caption: "Field1", Type: string(entity.FieldTypeShortText), Unique: &unique},
						{Name: "field2", EntName: "field2", Caption: "Field2", Type: string(entity.FieldTypeEmail), Unique: &unique},
						{Name: "field3", EntName: "field3", Caption: "Field3", Type: string(entity.FieldTypeInteger), Unique: &unique},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`Unique()`,
			},
			notInFile: []string{
				"Default(",
			},
		},
		{
			name: "entity with default value combinations",
			setupData: func() EntityTemplateData {
				defaultText := "text"
				defaultInt := "10"
				defaultFloat := "5.5"
				defaultBool := "false"
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "defaults",
						EntName:          "Defaults",
						Caption:          "Defaults",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "text",
					},
					Fields: []entity.Field{
						{Name: "text", EntName: "text", Caption: "Text", Type: string(entity.FieldTypeShortText), DefaultValue: &defaultText},
						{Name: "number", EntName: "number", Caption: "Number", Type: string(entity.FieldTypeInteger), DefaultValue: &defaultInt},
						{Name: "decimal", EntName: "decimal", Caption: "Decimal", Type: string(entity.FieldTypeFloat), DefaultValue: &defaultFloat},
						{Name: "flag", EntName: "flag", Caption: "Flag", Type: string(entity.FieldTypeBoolean), DefaultValue: &defaultBool},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`Default("text")`,
				`Default(10)`,
				`Default(5.5)`,
				`Default(false)`,
			},
		},
		{
			name: "entity with min max combinations",
			setupData: func() EntityTemplateData {
				minInt := "1"
				maxInt := "100"
				minFloat := "0.0"
				maxFloat := "1.0"
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "ranges",
						EntName:          "Ranges",
						Caption:          "Ranges",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "intField",
					},
					Fields: []entity.Field{
						{Name: "intField", EntName: "int_field", Caption: "Int Field", Type: string(entity.FieldTypeInteger), Min: &minInt, Max: &maxInt},
						{Name: "floatField", EntName: "float_field", Caption: "Float Field", Type: string(entity.FieldTypeFloat), Min: &minFloat, Max: &maxFloat},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`Min(1)`,
				`Max(100)`,
				`Min(0.0)`,
				`Max(1.0)`,
				`Min:     "1"`,
				`Max:     "100"`,
				`Min:     "0.0"`,
				`Max:     "1.0"`,
			},
		},
		{
			name: "entity with private field combinations",
			setupData: func() EntityTemplateData {
				private := true
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "private",
						EntName:          "Private",
						Caption:          "Private",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "public",
					},
					Fields: []entity.Field{
						{Name: "public", EntName: "public", Caption: "Public", Type: string(entity.FieldTypeShortText)},
						{Name: "password", EntName: "password", Caption: "Password", Type: string(entity.FieldTypeShortText), Private: &private},
						{Name: "token", EntName: "token", Caption: "Token", Type: string(entity.FieldTypeShortText), Private: &private},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`Private: true`,
			},
			notInFile: []string{
				`Private: false`,
			},
		},
		{
			name: "entity with nillable DateTime field",
			setupData: func() EntityTemplateData {
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "nillable",
						EntName:          "Nillable",
						Caption:          "Nillable",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "name",
					},
					Fields: []entity.Field{
						{
							Name:    "name",
							EntName: "name",
							Caption: "Name",
							Type:    string(entity.FieldTypeShortText),
						},
						{
							Name:     "deletedAt",
							EntName:  "deleted_at",
							Caption:  "Deleted At",
							Type:     string(entity.FieldTypeDateTime),
							Nillable: true,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.Time("deleted_at")`,
				`Nillable()`,
				"Optional()",
			},
		},
		{
			name: "entity with multiple field types",
			setupData: func() EntityTemplateData {
				required := true
				unique := true
				private := true
				minVal := "0"
				maxVal := "100"
				defaultVal := "default"
				return EntityTemplateData{
					Entity: entity.Entity{
						Name:             "product",
						EntName:          "Product",
						Caption:          "Product",
						Owner:            entity.EntityOwnerUser,
						DisplayFieldName: "name",
					},
					Fields: []entity.Field{
						{
							Name:     "name",
							EntName:  "name",
							Caption:  "Name",
							Type:     string(entity.FieldTypeShortText),
							Required: &required,
						},
						{
							Name:    "email",
							EntName: "email",
							Caption: "Email",
							Type:    string(entity.FieldTypeEmail),
							Unique:  &unique,
						},
						{
							Name:         "status",
							EntName:      "status",
							Caption:      "Status",
							Type:         string(entity.FieldTypeShortText),
							DefaultValue: &defaultVal,
						},
						{
							Name:    "price",
							EntName: "price",
							Caption: "Price",
							Type:    string(entity.FieldTypeInteger),
							Min:     &minVal,
							Max:     &maxVal,
						},
						{
							Name:    "secret",
							EntName: "secret",
							Caption: "Secret",
							Type:    string(entity.FieldTypeShortText),
							Private: &private,
						},
					},
					Edges: []entity.Edge{},
				}
			},
			expectedInFile: []string{
				`field.String("name")`,
				`field.String("email")`,
				`Unique()`,
				`field.String("status")`,
				`Default("default")`,
				`field.Int("price")`,
				`Min(0)`,
				`Max(100)`,
				`Private: true`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := NewMockFileSystem()
			mockSubClient := &MockSubscriptionClient{}

			sm := NewSchemaManager(
				mockFS,
				"./schema",
				"./generate",
				mockSubClient,
			)

			data := tt.setupData()

			err := sm.WriteEntityToSchema(data)
			if err != nil {
				t.Fatalf("WriteEntityToSchema failed: %v", err)
			}

			// Find the created file
			if len(mockFS.createCalls) == 0 {
				t.Fatal("Create was not called")
			}

			fileName := mockFS.createCalls[0]
			content, err := mockFS.GetFileContent(fileName)
			if err != nil {
				t.Fatalf("Failed to read file content: %v", err)
			}

			if content == "" {
				t.Fatal("File content is empty")
			}

			// Debug: print content if test fails (uncomment for debugging)
			// t.Logf("Generated content:\n%s", content)

			// Verify expected content
			for _, expected := range tt.expectedInFile {
				if !strings.Contains(content, expected) {
					t.Errorf("Expected to find '%s' in generated file, but it was not found\nGenerated content preview:\n%s", expected, content[:min(500, len(content))])
				}
			}

			// Verify content that should NOT be present
			for _, notExpected := range tt.notInFile {
				if strings.Contains(content, notExpected) {
					t.Errorf("Expected NOT to find '%s' in generated file, but it was found", notExpected)
				}
			}
		})
	}
}
