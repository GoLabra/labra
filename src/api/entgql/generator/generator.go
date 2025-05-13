package generator

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/GoLabra/labra/src/api/entgql/entity"
	"github.com/GoLabra/labra/src/api/entgql/templates"
	"github.com/GoLabra/labra/src/api/strcase"
	"github.com/GoLabra/labra/src/api/subscription"
)

type FileSystem interface {
	Open(name string) (*os.File, error)
	Create(name string) (*os.File, error)
	CopyDirectory(from, to string) error
	Remove(path string) error
	Exit(code int)
}

type SubscriptionClient interface {
	PublishAppStatusMessage(appStatus subscription.AppStatus)
}

type SchemaManager struct {
	fs                            FileSystem
	subscriptionClient            SubscriptionClient
	schemaDateImportPath          string
	generateLocation              string
	entschemaTemplateRelativePath string
	userSchemaRelativePath        string
	adminSchemaRelativePath       string
	userSchemaBackupRelativePath  string
	adminSchemaBackupRelativePath string
}

func NewSchemaManager(fs FileSystem, schemaPath string, generateLocation string, subscriptionClient SubscriptionClient) SchemaManager {
	return SchemaManager{
		fs:                            fs,
		subscriptionClient:            subscriptionClient,
		schemaDateImportPath:          "github.com/GoLabra/labra/src/api/entgql/date",
		generateLocation:              generateLocation,
		entschemaTemplateRelativePath: "entschema/*",
		userSchemaRelativePath:        fmt.Sprintf("%s/_user/", schemaPath),
		adminSchemaRelativePath:       fmt.Sprintf("%s/_admin/", schemaPath),
		userSchemaBackupRelativePath:  fmt.Sprintf("%s/_user_backup/", schemaPath),
		adminSchemaBackupRelativePath: fmt.Sprintf("%s/_admin_backup/", schemaPath),
	}
}

func (sm SchemaManager) WriteEntityToSchema(entityTemplateData EntityTemplateData) error {
	err := entityTemplateData.Validate()
	if err != nil {
		sm.RevertSchema()
		return err
	}

	tmpl, err := templates.LoadTemplate("entity.go.tmpl", sm.entschemaTemplateRelativePath, template.FuncMap{
		"Imports":   Imports,
		"CamelCase": strcase.ToLowerCamel,
		"IsTrue": func(val *bool) bool {
			return val != nil && *val
		},
		"IsFalse": func(val *bool) bool {
			return val == nil || !*val
		},
		"GetFieldRef": func(e entity.Edge, edges []entity.Edge) struct {
			Edge entity.Edge
			Ref  *entity.Edge
		} {
			return struct {
				Edge entity.Edge
				Ref  *entity.Edge
			}{
				Edge: e,
			}
		},
		"StringsJoin": func(val []string) string {
			var stringValues []string

			for _, value := range val {
				stringValues = append(stringValues, "\""+value+"\"")
			}

			return strings.Join(stringValues, ",")
		},
	})

	f, err := sm.fs.Create(fmt.Sprintf(sm.userSchemaRelativePath+"%s.go", strcase.ToSnake(entityTemplateData.Entity.Caption)))
	if err != nil {
		return fmt.Errorf("[WriteEntityToSchema]error creating schema file: %w", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, entityTemplateData)

	if err != nil {
		return fmt.Errorf("[WriteEntityToSchema] error executing template: %w", err)
	}

	return err
}

func (sm SchemaManager) RemoveEntityFromSchema(fileName string) error {
	err := sm.fs.Remove(fmt.Sprintf("%s/%s.go", sm.userSchemaRelativePath, fileName))
	if err != nil {
		return err
	}

	return nil
}

func (sm SchemaManager) Generate(ctx context.Context, entities []entity.Entity) error {
	sm.subscriptionClient.PublishAppStatusMessage(subscription.AppStatusGenerating)

	go func() {
		err := sm.RunGenerateCommand()

		if err != nil {
			sm.subscriptionClient.PublishAppStatusMessage(subscription.AppStatusReverting)

			sm.RevertSchema()

			err = sm.RunGenerateCommand()
			if err != nil {
				sm.subscriptionClient.PublishAppStatusMessage(subscription.AppStatusFatal)
				panic("FATAL ERROR app is in an unrecoverable state") // TODO log fatal
			}
		}

		// TODO fix ugly hack
		sm.subscriptionClient.PublishAppStatusMessage(subscription.AppStatusRestarting)
		time.Sleep(2 * time.Second)
		sm.fs.Exit(1)
	}()

	return nil
}

func (sm SchemaManager) RunGenerateCommand() error {
	cmd := exec.Command("go", "generate")

	cmd.Dir = sm.generateLocation

	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("[Generate] error running generate command: %w", err)
	}

	log.Println("@@@@@@@@@@@@@@@@@ ~ Generate Command Output ~ @@@@@@@@@@@@@@@@@")
	log.Println(string(out))
	log.Println("@@@@@@@@@@@@@@@@@ ~                         ~ @@@@@@@@@@@@@@@@@")

	return nil
}

func (sm SchemaManager) BackupSchema() error {
	err := sm.fs.CopyDirectory(sm.userSchemaRelativePath, sm.userSchemaBackupRelativePath)
	if err != nil {
		return fmt.Errorf("[BackupSchema] error backing up user schema: %w", err)
	}

	err = sm.fs.CopyDirectory(sm.adminSchemaRelativePath, sm.adminSchemaBackupRelativePath)
	if err != nil {
		return fmt.Errorf("[BackupSchema] error backing up admin schema: %w", err)
	}

	return nil
}

func (sm SchemaManager) RevertSchema() error {
	err := sm.fs.CopyDirectory(sm.userSchemaBackupRelativePath, sm.userSchemaRelativePath)
	if err != nil {
		return fmt.Errorf("[RevertSchema] error reverting user schema: %w", err)
	}

	err = sm.fs.CopyDirectory(sm.adminSchemaBackupRelativePath, sm.adminSchemaRelativePath)
	if err != nil {
		return fmt.Errorf("[RevertSchema] error reverting admin schema: %w", err)
	}

	return nil
}

func Imports(fields []entity.Field) map[string]bool {
	imports := map[string]bool{}

	for _, field := range fields {
		switch entity.FieldType(field.Type) {
		case entity.FieldTypeShortText:
			imports["entgo.io/ent/dialect"] = true
		case entity.FieldTypeLongText:
			imports["entgo.io/ent/dialect"] = true
		case entity.FieldTypeRichText:
			imports["entgo.io/ent/dialect"] = true
		case entity.FieldTypeEmail:
		case entity.FieldTypeInteger:
		case entity.FieldTypeDecimal:
			imports["entgo.io/ent/dialect"] = true
		case entity.FieldTypeFloat:
		case entity.FieldTypeBoolean:
		case entity.FieldTypeSingleChoice:
			imports["fmt"] = true
		case entity.FieldTypeMultipleChoice:
			imports["fmt"] = true
		case entity.FieldTypeDateTime:
			imports["entgo.io/ent/dialect"] = true
			if field.DefaultValue != nil {
				imports["time"] = true
				imports["github.com/GoLabra/labra/src/api/entgql/date"] = true
			}
			if field.UpdateDefault {
				imports["time"] = true
			}
		case entity.FieldTypeDate:
			imports["entgo.io/ent/dialect"] = true
			if field.DefaultValue != nil {
				imports["time"] = true
				imports["github.com/GoLabra/labra/src/api/entgql/date"] = true
			}
		case entity.FieldTypeTime:
			imports["entgo.io/ent/dialect"] = true
			if field.DefaultValue != nil {
				imports["time"] = true
				imports["github.com/GoLabra/labra/src/api/entgql/date"] = true
			}
		case entity.FieldTypeJson:
		case entity.FieldTypeEnum:
		case entity.FieldTypeEnums:
		}
	}

	return imports
}
