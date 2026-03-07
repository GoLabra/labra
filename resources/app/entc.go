//go:build ignore

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/GoLabra/labra/entgql/annotations"
	"github.com/GoLabra/labra/entgql/entity"
	"github.com/GoLabra/labra/entgql/templates"
	"github.com/GoLabra/labra/strcase"
	"github.com/GoLabra/labra/utils"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/mitchellh/mapstructure"
)

var templateFuncMap template.FuncMap

const templatesFolderPath = "./templates/%s"

func init() {
	pluralizeClient := pluralize.NewClient()
	templateFuncMap = entgql.TemplateFuncs
	templateFuncMap["LowerFirstLetter"] = strcase.LowerFirstLetter
	templateFuncMap["ToLower"] = strings.ToLower
	templateFuncMap["Singular"] = pluralizeClient.Singular
	templateFuncMap["Plural"] = pluralizeClient.Plural
	templateFuncMap["Camel"] = strcase.ToLowerCamel // TODO @David these camels can be improved
	templateFuncMap["ToCamel"] = strcase.ToCamel
	templateFuncMap["Pascal"] = strcase.ToPascal
	templateFuncMap["ToTitle"] = strcase.ToTitle
	templateFuncMap["CreateInputs"] = CreateInputs
	templateFuncMap["CustomFieldName"] = utils.CustomFieldName
	templateFuncMap["InputEdges"] = utils.InputEdges
	templateFuncMap["GetExtendedTypes"] = getExtendedTypes
	templateFuncMap["GraphqlInputName"] = GraphqlInputName
	templateFuncMap["GoInputName"] = GoInputName
	templateFuncMap["EntMutationFieldName"] = entMutationFieldName
	templateFuncMap["ShouldSkip"] = utils.ShouldSkip
	templateFuncMap["Ignore"] = func(t *gen.Type) bool {
		return t.Annotations["Entity"] == nil || t.Annotations["Entity"].(map[string]any)["Owner"] != "User"
	}
	templateFuncMap["hasStateAnnotation"] = func(t *gen.Type) bool {
		var entityAnnotations annotations.Entity
		err := mapstructure.Decode(t.Annotations["Entity"], &entityAnnotations)
		if err != nil {
			panic(err)
		}

		return entityAnnotations.State.Enabled
	}

	os.MkdirAll("./domain/repo", os.ModePerm)
	os.MkdirAll("./domain/resolvers", os.ModePerm)
	os.MkdirAll("./domain/svc", os.ModePerm)
	os.MkdirAll("./interfaces/svc", os.ModePerm)
	os.MkdirAll("./interfaces/repo", os.ModePerm)
	os.MkdirAll("./graphql", os.ModePerm)

}

func main() {
	templates, err := templates.Load(templateFuncMap)

	ex, err := entgql.NewExtension(
		entgql.WithTemplates(
			entgql.CollectionTemplate,
			entgql.EnumTemplate,
			entgql.NodeTemplate,
			entgql.PaginationTemplate,
			entgql.TransactionTemplate,
			entgql.EdgeTemplate,
			templates.MutationInput,
			templates.MutationSetEdge,
			templates.MutationAddEdges,
			// templates.MutationOldValues,
		),
		entgql.WithWhereInputs(true),
		entgql.WithConfigPath("./gqlgen.yml"),
		entgql.WithSchemaGenerator(),
		entgql.WithRelaySpec(true),
		entgql.WithSchemaPath("./graphql/schema.graphql"),
	)

	if err != nil {
		log.Fatalf("creating entgql extension: %w", err)
	}

	opts := []entc.Option{
		entc.Extensions(ex),
		entc.FeatureNames("sql/execquery", "sql/upsert"),
	}

	graph, err := entc.LoadGraph("./ent/schema", &gen.Config{})
	if err != nil {
		panic(err)
	}
	CreateSystemEntitiesReverseRelations(graph)

	if err := entc.Generate("./ent/schema", &gen.Config{
		Features: []gen.Feature{
			gen.FeatureVersionedMigration,
			gen.FeatureUpsert,
		},
		Hooks: []gen.Hook{
			CleanupUserFiles(),
			RunGraphTemplates(),
			RunNodeTemplates(),
			CreateLifecycleMethods(),
		},
		Target:  "./ent",
		Package: "app/ent",
	}, opts...); err != nil {
		log.Fatalf("running ent codegen for system models: %w", err)
	}
}

type TemplateData struct {
	Graph         *gen.Graph
	CreateInputs  map[string]map[string]string
	TypesTemplate map[string]string
}

type GraphqlSchemaTemplateData struct {
	Name                 string
	Owner                entity.EntityOwner
	PluralName           string
	GetManyOperationName string
}

func CleanupUserFiles() gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			DeleteUserFilesFromDirectory("./domain/resolvers", "go")
			DeleteUserFilesFromDirectory("./graphql", "graphql")
			DeleteOrphanedUserFiles(g, "./domain/repo", "go")
			DeleteOrphanedUserFiles(g, "./domain/svc", "go")
			DeleteOrphanedUserFiles(g, "./interfaces/repo", "go")
			DeleteOrphanedUserFiles(g, "./interfaces/svc", "go")
			return next.Generate(g)
		})
	}
}

func DeleteUserFilesFromDirectory(directoryName, fileExtension string) {
	files, err := filepath.Glob(directoryName + string(os.PathSeparator) + "user.*." + fileExtension)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

func DeleteOrphanedUserFiles(g *gen.Graph, directoryName, fileExtension string) {
	activeFiles := map[string]bool{}
	for _, node := range g.Nodes {
		var entityAnnotation annotations.Entity
		if err := mapstructure.Decode(node.Annotations[annotations.EntityName], &entityAnnotation); err != nil {
			continue
		}
		if entityAnnotation.Owner != entity.EntityOwnerUser {
			continue
		}
		fileName := "user." + strcase.ToSnake(node.Name) + "." + fileExtension
		activeFiles[fileName] = true
	}

	files, err := filepath.Glob(directoryName + string(os.PathSeparator) + "user.*." + fileExtension)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		baseName := filepath.Base(f)
		if !activeFiles[baseName] {
			if err := os.Remove(f); err != nil {
				panic(err)
			}
		}
	}
}

func runTemplate(errPrefix, tmplName, tmplPath, outPath string, data interface{}) error {
	tmpl, err := templates.LoadTemplate(tmplName, tmplPath, templateFuncMap)
	if err != nil {
		return fmt.Errorf("%s %w", errPrefix, err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("%s %w", errPrefix, err)
	}
	defer f.Close()
	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("%s %w", errPrefix, err)
	}
	return nil
}

func entityAnnotationForNode(node *gen.Type) (annotations.Entity, bool) {
	var entityAnnotation annotations.Entity
	if err := mapstructure.Decode(node.Annotations[annotations.EntityName], &entityAnnotation); err != nil {
		return annotations.Entity{}, false
	}
	return entityAnnotation, true
}

func outputFileName(node *gen.Type, owner entity.EntityOwner, suffix string) string {
	base := strcase.ToSnake(node.Name)
	if owner == entity.EntityOwnerUser {
		base = "user." + base
	}
	return base + suffix
}

func buildCreateInputsForNode(n *gen.Type) map[string]map[string]string {
	createInputs := map[string]map[string]string{}
	for _, e := range n.Edges {
		inputName := fmt.Sprintf("Create%sWithout%sInput", e.Type.Name, n.Name)
		if e.Ref == nil || e.Ref.Optional {
			continue
		}
		if _, ok := createInputs[inputName]; ok {
			continue
		}
		if e.Name == "ref_created_by" || e.Name == "ref_updated_by" {
			continue
		}
		oneInputName := fmt.Sprintf("CreateOne%sWithout%sInput", e.Type.Name, n.Name)
		manyInputName := fmt.Sprintf("CreateMany%sWithout%sInput", e.Type.Name, n.Name)
		createInputs[oneInputName] = map[string]string{
			"connect": fmt.Sprintf("%sWhereUniqueInput", e.Type.Name),
			"create":  inputName,
		}
		createInputs[manyInputName] = map[string]string{
			"connect": fmt.Sprintf("[%sWhereUniqueInput!]", e.Type.Name),
			"create":  fmt.Sprintf("[%s!]", inputName),
		}
		createInputs[inputName] = map[string]string{}
		for _, f := range e.Type.Fields {
			scalar := mapScalar(f)
			if !f.Optional {
				scalar += "!"
			}
			createInputs[inputName][strcase.ToLowerCamel(f.Name)] = scalar
		}
		for _, ee := range e.Type.Edges {
			if ee.Unique {
				if ee.Ref == nil || ee.Ref.Optional {
					createInputs[inputName][strcase.ToLowerCamel(ee.Name)] = fmt.Sprintf("CreateOne%sInput", ee.Type.Name)
				} else {
					createInputs[inputName][strcase.ToLowerCamel(ee.Name)] = fmt.Sprintf("CreateOne%sWithout%sInput", ee.Type.Name, e.Type.Name)
				}
			} else {
				if ee.Ref == nil || ee.Ref.Optional {
					createInputs[inputName][strcase.ToLowerCamel(ee.Name)] = fmt.Sprintf("CreateMany%sInput", ee.Type.Name)
				} else {
					createInputs[inputName][strcase.ToLowerCamel(ee.Name)] = fmt.Sprintf("CreateMany%sWithout%sInput", ee.Type.Name, e.Type.Name)
				}
			}
		}
	}
	return createInputs
}

func RunGraphTemplates() gen.Hook {
	errPrefix := "[RunGraphTemplates]"
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			t := TemplateData{
				Graph: g,
				TypesTemplate: map[string]string{
					"string":                  "String",
					"bool":                    "Boolean",
					"time.Time":               "Time",
					"int":                     "Int",
					"float64":                 "Float",
					"map[string]interface {}": "Map",
				},
			}
			if err := runTemplate(errPrefix, "unique_inputs.graphql.tmpl", "graphql/unique_inputs.graphql.tmpl", "./graphql/unique_inputs.graphql", t); err != nil {
				return err
			}
			if err := runTemplate(errPrefix, "gql_where_unique_input.go.tmpl", "ent/gql_where_unique_input.go.tmpl", "./ent/gql_where_unique_input.go", g); err != nil {
				return err
			}
			if err := runTemplate(errPrefix, "repository.go.tmpl", "repo/repository.go.tmpl", "./domain/repo/repository.go", g); err != nil {
				return err
			}
			if err := runTemplate(errPrefix, "service.go.tmpl", "svc/service.go.tmpl", "./domain/svc/service.go", g); err != nil {
				return err
			}
			if err := runTemplate(errPrefix, "resolver.go.tmpl", "resolver/resolver.go.tmpl", "./domain/resolvers/resolver.go", nil); err != nil {
				return err
			}
			if err := runTemplate(errPrefix, "tx.go.tmpl", "repo/tx.go.tmpl", "./domain/repo/tx.go", g); err != nil {
				return err
			}
			return next.Generate(g)
		})
	}
}

func RunNodeTemplates() gen.Hook {
	errPrefix := "[RunNodeTemplates]"
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			for _, node := range g.Nodes {
				entityAnnotation, ok := entityAnnotationForNode(node)
				if !ok || entityAnnotation.Owner != entity.EntityOwnerUser {
					continue
				}
				baseName := outputFileName(node, entityAnnotation.Owner, "")

				createInputs := buildCreateInputsForNode(node)
				if err := runTemplate(errPrefix, "entity.graphql.tmpl", "graphql/entity.graphql.tmpl", "./graphql/"+baseName+".graphql", struct {
					Node         *gen.Type
					CreateInputs map[string]map[string]string
				}{node, createInputs}); err != nil {
					return err
				}
				if _, err := os.Stat("./domain/repo/" + baseName + ".go"); errors.Is(err, os.ErrNotExist) {
					if err := runTemplate(errPrefix, "entity.go.tmpl", "repo/*", "./domain/repo/"+baseName+".go", node); err != nil {
						return err
					}
				}
				if _, err := os.Stat("./domain/svc/" + baseName + ".go"); errors.Is(err, os.ErrNotExist) {
					if err := runTemplate(errPrefix, "entity.go.tmpl", "svc/entity.go.tmpl", "./domain/svc/"+baseName+".go", node); err != nil {
						return err
					}
				}
				if _, err := os.Stat("./interfaces/svc/" + baseName + ".go"); errors.Is(err, os.ErrNotExist) {
					if err := runTemplate(errPrefix, "interface.go.tmpl", "svc/interface.go.tmpl", "./interfaces/svc/"+baseName+".go", node); err != nil {
						return err
					}
				}
				if _, err := os.Stat("./interfaces/repo/" + baseName + ".go"); errors.Is(err, os.ErrNotExist) {
					if err := runTemplate(errPrefix, "interface.go.tmpl", "repo/interface.go.tmpl", "./interfaces/repo/"+baseName+".go", node); err != nil {
						return err
					}
				}
				if err := runTemplate(errPrefix, "entity.resolver.go.tmpl", "resolver/entity.resolver.go.tmpl", "./domain/resolvers/"+baseName+".resolvers.go", node); err != nil {
					return err
				}
			}
			return next.Generate(g)
		})
	}
}

func CreateSystemEntitiesReverseRelations(g *gen.Graph) error {
	errFormat := "[CreateSystemEntitiesReverseRelations] %w"
	f, _ := os.Create("./ent/schema/additional_edges.go")

	tmpl, err := templates.LoadTemplate("additional_edges.go.tmpl", "entschema/additional_edges.go.tmpl", templateFuncMap)
	if err != nil {
		return fmt.Errorf(errFormat, fmt.Errorf("error parsing template file: %w", err))
	}
	additionalEdges := []*gen.Edge{}

	for _, node := range g.Nodes {
		if node.Annotations["Entity"] == nil || node.Annotations["Entity"].(map[string]any)["Owner"] != "User" {
			continue
		}

		for _, edge := range node.Edges {
			if edge.Unique {
				continue
			}
			edge.Owner.ClientName()

			if owner, ok := node.Annotations["Entity"].(map[string]any)["Owner"].(string); !ok || owner != "User" {
				continue
			}

			additionalEdges = append(additionalEdges, edge)
		}
	}

	err = tmpl.Execute(f, additionalEdges)
	if err != nil {
		f.Close()
		return fmt.Errorf(errFormat, fmt.Errorf("error executing template: %w", err))
	}

	f.Close()
	return nil
}

func getExtendedTypes(fields []*gen.Field) []*gen.Field {
	var extendedFields []*gen.Field
	for _, field := range fields {
		if _, ok := field.Annotations["CustomType"]; ok {
			extendedFields = append(extendedFields, field)
		}
	}
	return extendedFields
}

func entMutationFieldName(name string) string {
	var fieldTypeMap = map[string]string{
		"path": "_path",
		"type": "_type",
	}
	if _, ok := fieldTypeMap[name]; !ok {
		return name
	}
	return fieldTypeMap[name]
}

// mapScalar provides maps an ent.Schema type into GraphQL scalar type.
func mapScalar(f *gen.Field) string {
	if f.Annotations["EntGQL"] != nil && f.Annotations["EntGQL"].(map[string]any)["Type"] != nil && f.Annotations["EntGQL"].(map[string]any)["Type"] != nil && f.Annotations["EntGQL"].(map[string]any)["Type"].(string) != "" {
		return f.Annotations["EntGQL"].(map[string]any)["Type"].(string)
	}

	scalar := f.Type.String()
	switch t := f.Type.Type; {
	case f.Name == "id":
		return "ID"
	case f.IsEdgeField():
		scalar = "ID"
	case t.Float():
		scalar = "Float"
	case t.Integer():
		scalar = "Int"
	case t == field.TypeString:
		scalar = "String"
	case t == field.TypeBool:
		scalar = "Boolean"
	case strings.ContainsRune(scalar, '.'): // Time, Enum or Other.
		scalar = scalar[strings.LastIndexByte(scalar, '.')+1:]
		if f.IsEnum() {
			// Use the GQL type as enum prefix. e.g. Todo.status
			// will generate an enum named "TodoStatus".
			//scalar = gqlType + scalar
		}
		if f.Type.RType != nil && f.Type.RType.Name == "" {
			switch f.Type.RType.Kind {
			case reflect.Slice, reflect.Array:
				if strings.HasPrefix(f.Type.RType.Ident, "[]*") {
					scalar = "[" + scalar + "]"
				} else {
					scalar = "[" + scalar + "!]"
				}
			}
		}
	case t == field.TypeJSON:
		scalar = ""
		if f.Type.RType != nil {
			switch f.Type.RType.Kind {
			case reflect.Slice, reflect.Array:
				switch f.Type.RType.Ident {
				case "[]float64":
					scalar = "[Float!]"
				case "[]int":
					scalar = "[Int!]"
				case "[]string":
					scalar = "[String!]"
				}
			case reflect.Map:
				if f.Type.RType.Ident == "map[string]interface {}" {
					scalar = "Map"
					if !f.Optional {
						scalar += "!"
					}
				}
			}
		}
	}
	return scalar
}

func CreateInputs(nodes []*gen.Type) map[string]map[string]string {
	var createInputs = map[string]map[string]string{}
	for _, n := range nodes {
		if n.Annotations["Entity"] != nil && n.Annotations["Entity"].(map[string]any)["Owner"] != "User" {
			continue
		}
		for _, e := range n.Edges {
			if e.Type.Annotations["Entity"] == nil || e.Type.Annotations["Entity"].(map[string]any)["Owner"] != "User" {
				continue
			}
			var inputName = fmt.Sprintf("Create%sWithout%sInput", e.Type.Name, n.Name)
			if e.Ref == nil || e.Ref.Optional {
				continue
				inputName = fmt.Sprintf("Create%sInput", e.Type.Name)
			}
			if _, ok := createInputs[inputName]; ok {
				continue
			}
			if e.Name == "ref_created_by" || e.Name == "ref_updated_by" {
				continue
			}

			var oneInputName = fmt.Sprintf("CreateOne%sWithout%sInput", e.Type.Name, n.Name)
			var manyInputName = fmt.Sprintf("CreateMany%sWithout%sInput", e.Type.Name, n.Name)
			createInputs[oneInputName] = map[string]string{
				"Connect": fmt.Sprintf("*%sWhereUniqueInput", e.Type.Name),
				"Create":  "*" + inputName,
			}
			createInputs[manyInputName] = map[string]string{
				"Connect": fmt.Sprintf("[]%sWhereUniqueInput", e.Type.Name),
				"Create":  fmt.Sprintf("[]%s", inputName),
			}

			createInputs[inputName] = map[string]string{}
			for _, f := range e.Type.Fields {
				createInputs[inputName][f.StructField()] = f.Type.String()
				if IsPointer(f) {
					createInputs[inputName][f.StructField()] = fmt.Sprintf("*%s", f.Type.String())
				}
			}
			for _, ee := range e.Type.Edges {
				if ee.Unique {
					if ee.Ref == nil || ee.Ref.Optional {
						createInputs[inputName][strcase.ToLowerCamel(ee.Name)] = fmt.Sprintf("CreateOne%sInput", ee.Type.Name)
					} else {
						createInputs[inputName][strcase.ToLowerCamel(ee.Name)] = fmt.Sprintf("CreateOne%sWithout%sInput", ee.Type.Name, e.Type.Name)
					}
				} else {
					if ee.Ref == nil || ee.Ref.Optional {
						createInputs[inputName][strcase.ToLowerCamel(ee.Name)] = fmt.Sprintf("CreateMany%sInput", ee.Type.Name)
					} else {
						createInputs[inputName][strcase.ToLowerCamel(ee.Name)] = fmt.Sprintf("CreateMany%sWithout%sInput", ee.Type.Name, e.Type.Name)
					}
				}
			}
		}
	}
	return createInputs
}

func CreateLifecycleMethods() gen.Hook {
	errPrefix := "[CreateLifecycleMethods]"
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			if err := next.Generate(g); err != nil {
				panic(err)
			}
			if err := runTemplate(errPrefix, "lifecycle.go.tmpl", "ent/lifecycle.go.tmpl", "./ent/lifecycle.go", g); err != nil {
				panic(err)
			}
			return runTemplate(errPrefix, "lifecycle_filter_methods.go.tmpl", "ent/lifecycle_filter_methods.go.tmpl", "./ent/lifecycle_filter_methods.go", g)
		})
	}
}

func GoInputName(isCreate bool, node *gen.Type, edge *gen.Edge) string {
	input := InputName(isCreate, node, edge)
	if edge.Type.Annotations["Entity"] == nil || edge.Type.Annotations["Entity"].(map[string]any)["Owner"] != "User" {
		input = "admin." + input
	}
	if !isCreate || edge.Optional {
		input = "*" + input
	}
	return input
}

func GraphqlInputName(isCreate bool, node *gen.Type, edge *gen.Edge) string {
	input := InputName(isCreate, node, edge)
	if isCreate && !edge.Optional {
		input += "!"
	}
	return input
}

func InputName(isCreate bool, node *gen.Type, edge *gen.Edge) string {
	input := "Update"
	if isCreate {
		input = "Create"
	}
	if edge.Unique {
		input += fmt.Sprintf("One%s", edge.Type.Name)
	} else {
		input += fmt.Sprintf("Many%s", edge.Type.Name)
	}

	if edge.Ref != nil && !edge.Ref.Optional {
		input += fmt.Sprintf("Without%s", node.Name)
	}

	input += "Input"

	return input
}

func IsPointer(f *gen.Field) bool {
	if f.Type.Nillable || f.Type.RType.IsPtr() {
		return false
	}
	return f.Optional || f.Default || f.DefaultFunc()
}
