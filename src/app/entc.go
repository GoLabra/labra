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
	"github.com/GoLabra/labra/src/api/entgql/annotations"
	"github.com/GoLabra/labra/src/api/entgql/entity"
	"github.com/GoLabra/labra/src/api/entgql/templates"
	"github.com/GoLabra/labra/src/api/strcase"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/mitchellh/mapstructure"
)

var templateFuncMap template.FuncMap

const templatesFolderPath = "./templates/%s"

func init() {
	pluralizeClient := pluralize.NewClient()
	templateFuncMap = entgql.TemplateFuncs
	templateFuncMap["ToUpper"] = strings.Title
	templateFuncMap["LowerFirstLetter"] = func(val string) string {
		return strings.ToLower(val[:1]) + val[1:]
	}
	templateFuncMap["ToLower"] = strings.ToLower
	templateFuncMap["Singular"] = pluralizeClient.Singular
	templateFuncMap["Plural"] = pluralizeClient.Plural
	templateFuncMap["Camel"] = strcase.ToLowerCamel
	templateFuncMap["Pascal"] = strcase.ToPascal
	templateFuncMap["ToTitle"] = func(val string) string {
		return strings.ToTitle(val[:1]) + val[1:]
	}
	templateFuncMap["CreateInputs"] = CreateInputs
	templateFuncMap["CustomFieldName"] = func(val string) string {
		val = strings.TrimSuffix(val, "_id")
		val = strings.TrimSuffix(val, "_ids")
		return val
	}

	templateFuncMap["InputEdges"] = func(m *entgql.MutationDescriptor) []*gen.Edge {
		inputEdges := make([]*gen.Edge, 0, len(m.Type.Edges))
		for _, e := range m.Type.Edges {
			if e.Type.IsEdgeSchema() || e.Immutable || e.Annotations["Skip"] == entgql.SkipMutationUpdateInput {
				continue
			}
			inputEdges = append(inputEdges, e)
		}
		return inputEdges
	}
	templateFuncMap["GetExtendedTypes"] = getExtendedTypes
	templateFuncMap["GraphqlInputName"] = GraphqlInputName
	templateFuncMap["GoInputName"] = GoInputName
	templateFuncMap["EntMutationFieldName"] = entMutationFieldName

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
			templates.MutationUpdatedFields,
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

	if err := entc.Generate("./ent/schema", &gen.Config{
		Features: []gen.Feature{
			gen.FeatureVersionedMigration,
			gen.FeatureUpsert,
		},
		Hooks: []gen.Hook{
			CleanupUserFiles(),
			CreateGraphqlUniqueInputs(),
			CreateEntUniqueInputs(),
			CreateGraphqlSchema(),
			CreateServiceInterface(),
			CreateRepositoryInterface(),
			CreateTxRepo(),
			CreateRepository(),
			CreateRepositories(),
			CreateService(),
			CreateServices(),
			CreateResolverFile(),
			CreateResolvers(),
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
			DeleteUserFilesFromDirectory("./domain/repo", "go")
			DeleteUserFilesFromDirectory("./domain/svc", "go")
			DeleteUserFilesFromDirectory("./graphql", "graphql")
			DeleteUserFilesFromDirectory("./interfaces/repo", "go")
			DeleteUserFilesFromDirectory("./interfaces/svc", "go")
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

func CreateEntUniqueInputs() gen.Hook {
	errFormat := "[CreateEntUniqueInputs] %w"
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			tmpl, err := templates.LoadTemplate("gql_where_unique_input.go.tmpl", "ent/gql_where_unique_input.go.tmpl", templateFuncMap)
			if err != nil {
				return fmt.Errorf(errFormat, fmt.Errorf("error parsing template file: %w", err))
			}

			f, err := os.Create("./ent/gql_where_unique_input.go")
			if err != nil {
				return fmt.Errorf(errFormat, fmt.Errorf("error creating file: %w", err))
			}

			err = tmpl.Execute(f, g)
			if err != nil {
				f.Close()
				return fmt.Errorf(errFormat, fmt.Errorf("error executing template: %w", err))
			}

			f.Close()
			return next.Generate(g)
		})
	}
}

func CreateGraphqlUniqueInputs() gen.Hook {
	errFormat := "[CreateGraphqlUniqueInputs] %w"
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

			tmpl, err := templates.LoadTemplate("unique_inputs.graphql.tmpl", "graphql/unique_inputs.graphql.tmpl", templateFuncMap)
			if err != nil {
				return fmt.Errorf(errFormat, fmt.Errorf("error parsing template file: %w", err))
			}

			if err != nil {
				return fmt.Errorf(errFormat, fmt.Errorf("error creating folder: %w", err))
			}

			f, err := os.Create("./graphql/unique_inputs.graphql")
			if err != nil {
				return fmt.Errorf(errFormat, fmt.Errorf("error creating file: %w", err))
			}

			err = tmpl.Execute(f, t)
			if err != nil {
				f.Close()
				return fmt.Errorf(errFormat, fmt.Errorf("error executing template: %w", err))
			}

			f.Close()
			return next.Generate(g)
		})
	}
}

func CreateGraphqlSchema() gen.Hook {
	errFormat := "[CreateGraphqlSchema] %w"
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			for _, n := range g.Nodes {
				var createInputs = map[string]map[string]string{}
				var entityAnnotation = annotations.Entity{}

				err := mapstructure.Decode(n.Annotations[annotations.EntityName], &entityAnnotation)

				if err != nil {
					panic(err) // TODO @David do something
				}

				for _, e := range n.Edges {
					var inputName = fmt.Sprintf("Create%sWithout%sInput", e.Type.Name, n.Name)
					if e.Ref == nil || e.Ref.Optional {
						continue
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

				fileName := strcase.ToSnake(n.Name) + ".graphql"

				if entityAnnotation.Owner == entity.EntityOwnerUser {
					fileName = "user." + fileName
				}

				if _, err := os.Stat("./graphql/" + fileName); err == nil {
					//continue
				} else if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf(errFormat, fmt.Errorf("error checking if file already exists: %w", err))
				}

				tmpl, err := templates.LoadTemplate("entity.graphql.tmpl", "graphql/entity.graphql.tmpl", templateFuncMap)
				if err != nil {
					return fmt.Errorf(errFormat, fmt.Errorf("error parsing template file: %w", err))
				}

				f, err := os.Create("./graphql/" + fileName)
				if err != nil {
					return fmt.Errorf(errFormat, fmt.Errorf("error creating graphql file: %w", err))
				}

				err = tmpl.Execute(f, struct {
					Node         *gen.Type
					CreateInputs map[string]map[string]string
				}{n, createInputs})
				if err != nil {
					f.Close()
				}

				f.Close()
			}
			return next.Generate(g)
		})
	}
}

func CreateRepository() gen.Hook {
	errFormat := "[CreateRepository] %w"
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			tmpl, err := templates.LoadTemplate("repository.go.tmpl", "repo/repository.go.tmpl", templateFuncMap)
			if err != nil {
				return fmt.Errorf(errFormat, fmt.Errorf("error parsing template file: %w", err))
			}

			f, err := os.Create("./domain/repo/repository.go")
			if err != nil {
				return fmt.Errorf(errFormat, fmt.Errorf("error creating repository file: %w", err))
			}

			err = tmpl.Execute(f, g)
			if err != nil {
				f.Close()
				return fmt.Errorf(errFormat, fmt.Errorf("error executing template: %w", err))
			}
			f.Close()
			return next.Generate(g)
		})
	}
}

func CreateRepositories() gen.Hook {
	errFormat := "[CreateRepositories] %w"
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			for _, node := range g.Nodes {
				var entityAnnotation = annotations.Entity{}

				err := mapstructure.Decode(node.Annotations[annotations.EntityName], &entityAnnotation)

				if err != nil {
					panic(err) // TODO treat errir
				}

				fileName := strcase.ToSnake(node.Name) + ".go"

				if entityAnnotation.Owner == entity.EntityOwnerUser {
					fileName = "user." + fileName
				}

				if _, err := os.Stat("./domain/repo/" + fileName); err == nil {
					// continue
				} else if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf(errFormat, fmt.Errorf("error checking if file already exists: %w", err))
				}

				tmpl, err := templates.LoadTemplate("entity.go.tmpl", "repo/*", templateFuncMap)

				if err != nil {
					return fmt.Errorf(errFormat, fmt.Errorf("error parsing template file: %w", err))
				}

				f, err := os.Create("./domain/repo/" + fileName)
				if err != nil {
					return fmt.Errorf(errFormat, fmt.Errorf("error creating graphql file: %w", err))
				}

				err = tmpl.Execute(f, node)
				if err != nil {
					f.Close()
					return fmt.Errorf(errFormat, fmt.Errorf("error executing template: %w", err))
				}

				f.Close()
			}
			return next.Generate(g)
		})
	}
}

func CreateService() gen.Hook {
	errFormat := "[CreateService] %w"
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			tmpl, err := templates.LoadTemplate("service.go.tmpl", "svc/service.go.tmpl", templateFuncMap)
			if err != nil {
				return fmt.Errorf(errFormat, fmt.Errorf("error parsing template file: %w", err))
			}

			f, err := os.Create("./domain/svc/service.go")
			if err != nil {
				return fmt.Errorf(errFormat, fmt.Errorf("error creating service file: %w", err))
			}

			err = tmpl.Execute(f, g)
			if err != nil {
				f.Close()
				return fmt.Errorf(errFormat, fmt.Errorf("error executing template: %w", err))
			}
			f.Close()
			return next.Generate(g)
		})
	}
}

func CreateServices() gen.Hook {
	errFormat := "[CreateServices] %w"
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			for _, node := range g.Nodes {

				var entityAnnotation = annotations.Entity{}

				err := mapstructure.Decode(node.Annotations[annotations.EntityName], &entityAnnotation)

				if err != nil {
					panic(err) // TODO treat errir
				}

				fileName := strcase.ToSnake(node.Name) + ".go"

				if entityAnnotation.Owner == entity.EntityOwnerUser {
					fileName = "user." + fileName
				}

				if _, err := os.Stat("./domain/svc/" + fileName); err == nil {
					// continue
				} else if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf(errFormat, fmt.Errorf("error checking if file already exists: %w", err))
				}

				tmpl, err := templates.LoadTemplate("entity.go.tmpl", "svc/entity.go.tmpl", templateFuncMap)
				if err != nil {
					return fmt.Errorf(errFormat, fmt.Errorf("error parsing template file: %w", err))
				}

				f, err := os.Create("./domain/svc/" + fileName)
				if err != nil {
					return fmt.Errorf(errFormat, fmt.Errorf("error creating graphql file: %w", err))
				}

				err = tmpl.Execute(f, node)
				if err != nil {
					f.Close()
					return fmt.Errorf(errFormat, fmt.Errorf("error executing template: %w", err))
				}
				f.Close()
			}
			return next.Generate(g)
		})
	}
}

func CreateServiceInterface() gen.Hook {
	errFormat := "[CreateServiceInterface] %w"
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			for _, node := range g.Nodes {
				var entityAnnotation = annotations.Entity{}

				err := mapstructure.Decode(node.Annotations[annotations.EntityName], &entityAnnotation)

				if err != nil {
					panic(err) // TODO treat errir
				}

				fileName := strcase.ToSnake(node.Name) + ".go"

				if entityAnnotation.Owner == entity.EntityOwnerUser {
					fileName = "user." + fileName
				}

				if _, err := os.Stat("./interfaces/svc/" + fileName); err == nil {
					// continue
				} else if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf(errFormat, fmt.Errorf("error checking if file already exists: %w", err))
				}
				tmpl, err := templates.LoadTemplate("interface.go.tmpl", "svc/interface.go.tmpl", templateFuncMap)
				if err != nil {
					return fmt.Errorf(errFormat, fmt.Errorf("error parsing template file: %w", err))
				}

				data := GraphqlSchemaTemplateData{
					Name:  node.Name,
					Owner: entityAnnotation.Owner,
				}

				f, err := os.Create("./interfaces/svc/" + fileName)
				if err != nil {
					return fmt.Errorf(errFormat, fmt.Errorf("error creating graphql file: %w", err))
				}

				err = tmpl.Execute(f, data)
				if err != nil {
					f.Close()
					return fmt.Errorf(errFormat, fmt.Errorf("error executing template: %w", err))
				}
				f.Close()
			}
			return next.Generate(g)
		})
	}
}

func CreateRepositoryInterface() gen.Hook {
	errFormat := "[CreateRepositoryInterface] %w"
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			for _, node := range g.Nodes {
				var entityAnnotation = annotations.Entity{}

				err := mapstructure.Decode(node.Annotations[annotations.EntityName], &entityAnnotation)

				if err != nil {
					panic(err) // TODO treat errir
				}

				fileName := strcase.ToSnake(node.Name) + ".go"

				if entityAnnotation.Owner == entity.EntityOwnerUser {
					fileName = "user." + fileName
				}

				if _, err := os.Stat("./interfaces/repo/" + fileName); err == nil {
					// continue
				} else if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf(errFormat, fmt.Errorf("error checking if file already exists: %w", err))
				}

				tmpl, err := templates.LoadTemplate("interface.go.tmpl", "repo/interface.go.tmpl", templateFuncMap)
				if err != nil {
					return fmt.Errorf(errFormat, fmt.Errorf("error parsing template file: %w", err))
				}

				data := GraphqlSchemaTemplateData{
					Name:  node.Name,
					Owner: entityAnnotation.Owner,
				}

				f, err := os.Create("./interfaces/repo/" + fileName)
				if err != nil {
					return fmt.Errorf(errFormat, fmt.Errorf("error creating graphql file: %w", err))
				}

				err = tmpl.Execute(f, data)
				if err != nil {
					f.Close()
					return fmt.Errorf(errFormat, fmt.Errorf("error executing template: %w", err))
				}
				f.Close()
			}
			return next.Generate(g)
		})
	}
}

func CreateResolvers() gen.Hook {
	errFormat := "[CreateResolvers] %w"
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			for _, node := range g.Nodes {
				var entityAnnotation = annotations.Entity{}

				err := mapstructure.Decode(node.Annotations[annotations.EntityName], &entityAnnotation)

				if err != nil {
					panic(err) // TODO treat errir
				}

				fileName := strcase.ToSnake(node.Name) + ".resolvers.go"

				if entityAnnotation.Owner == entity.EntityOwnerUser {
					fileName = "user." + fileName
				}

				if _, err := os.Stat("./domain/resolvers/" + fileName); err == nil {
					// continue
				} else if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf(errFormat, fmt.Errorf("error checking if file already exists: %w", err))
				}

				tmpl, err := templates.LoadTemplate("entity.resolver.go.tmpl", "resolver/entity.resolver.go.tmpl", templateFuncMap)
				if err != nil {
					return fmt.Errorf(errFormat, fmt.Errorf("error parsing template file: %w", err))
				}

				data := GraphqlSchemaTemplateData{
					Name:  node.Name,
					Owner: entityAnnotation.Owner,
				}

				f, err := os.Create("./domain/resolvers/" + fileName)
				if err != nil {
					return fmt.Errorf(errFormat, fmt.Errorf("error creating graphql file: %w", err))
				}

				err = tmpl.Execute(f, data)
				if err != nil {
					f.Close()
					return fmt.Errorf(errFormat, fmt.Errorf("error executing template: %w", err))
				}
				f.Close()
			}
			return next.Generate(g)
		})
	}
}

func CreateResolverFile() gen.Hook {
	errFormat := "[CreateResolverFile] %w"
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			f, _ := os.Create("./domain/resolvers/resolver.go")

			tmpl, err := templates.LoadTemplate("resolver.go.tmpl", "resolver/resolver.go.tmpl", templateFuncMap)
			if err != nil {
				return fmt.Errorf(errFormat, fmt.Errorf("error parsing template file: %w", err))
			}

			err = tmpl.Execute(f, nil)
			if err != nil {
				f.Close()
				return fmt.Errorf(errFormat, fmt.Errorf("error executing template: %w", err))
			}

			f.Close()
			return next.Generate(g)
		})
	}
}

func CreateTxRepo() gen.Hook {
	errFormat := "[CreateTxRepo] %w"
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			f, err := os.Create("./domain/repo/tx.go")

			tmpl, err := templates.LoadTemplate("tx.go.tmpl", "repo/tx.go.tmpl", templateFuncMap)
			if err != nil {
				return fmt.Errorf(errFormat, fmt.Errorf("error parsing template file: %w", err))
			}

			err = tmpl.Execute(f, g)
			if err != nil {
				f.Close()
				return fmt.Errorf(errFormat, fmt.Errorf("error executing template: %w", err))
			}

			f.Close()
			return next.Generate(g)
		})
	}
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
		for _, e := range n.Edges {
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

func GoInputName(isCreate bool, node *gen.Type, edge *gen.Edge) string {
	input := InputName(isCreate, node, edge)
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
