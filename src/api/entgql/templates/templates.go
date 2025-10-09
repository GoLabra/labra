package templates

import (
	"embed"
	"log"
	"text/template"

	"entgo.io/ent/entc/gen"
)

type store struct {
	MutationInput      *gen.Template
	MutationSetEdge    *gen.Template
	MutationAddEdges   *gen.Template
	MutationOldValues  *gen.Template
	AdditionalFields   *gen.Template
	OwnerFilterMethods *gen.Template
}

var (
	//go:embed */*
	_templates embed.FS
)

func loadT(name, path string, templateFuncMap template.FuncMap) (*gen.Template, error) {
	return gen.NewTemplate(name).Funcs(templateFuncMap).ParseFS(_templates, path)
}

func Load(templateFuncMap template.FuncMap) (store, error) {
	var templateStore = store{}
	var err error
	templateStore.MutationAddEdges, err = loadT("mutation_add_edges", "entgql/mutation_add_edges.go.tmpl", templateFuncMap)
	if err != nil {
		log.Fatalf("cannot parse mutation input template: %v", err)
		return templateStore, err
	}

	templateStore.MutationOldValues, err = loadT("mutation_old_values", "entgql/mutation_old_values.go.tmpl", templateFuncMap)
	if err != nil {
		log.Fatalf("cannot parse mutation set edge template: %v", err)
		return templateStore, err
	}

	templateStore.MutationSetEdge, err = loadT("mutation_set_edge", "entgql/mutation_set_edge.go.tmpl", templateFuncMap)
	if err != nil {
		log.Fatalf("cannot parse mutation add edges template: %v", err)
		return templateStore, err
	}

	templateStore.MutationInput, err = loadT("mutation_input.go.tmpl", "entgql/mutation_input.go.tmpl", templateFuncMap)
	if err != nil {
		log.Fatalf("cannot parse mutation updated fields template: %v", err)
		return templateStore, err
	}

	templateStore.AdditionalFields, err = loadT("additional_fields.go.tmpl", "entgql/additional_fields.go.tmpl", templateFuncMap)
	if err != nil {
		log.Fatalf("cannot parse additional fields template: %v", err)
		return templateStore, err
	}

	templateStore.OwnerFilterMethods, err = loadT("owner_filter_methods.go.tmpl", "ent/owner_filter_methods.go.tmpl", templateFuncMap)
	if err != nil {
		log.Fatalf("cannot parse owner filter methods template: %v", err)
		return templateStore, err
	}

	return templateStore, nil
}

func LoadTemplate(name string, path string, templateFuncMap template.FuncMap) (*gen.Template, error) {
	return loadT(name, path, templateFuncMap)
}
