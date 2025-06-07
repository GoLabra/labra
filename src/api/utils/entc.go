package utils

import (
	"strings"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc/gen"
)

func ShouldSkip(n *gen.Type, skipOp string) bool {
	ant := &entgql.Annotation{}
	if err := ant.Decode(n.Annotations[ant.Name()]); err != nil {
		panic(err)
	}
	switch skipOp {
	case "mutation_create_input":
		return ant.Skip.Is(entgql.SkipMutationCreateInput)
	case "mutation_update_input":
		return ant.Skip.Is(entgql.SkipMutationUpdateInput)
	}
	return false
}

func InputEdges(m *entgql.MutationDescriptor) []*gen.Edge {
	inputEdges := make([]*gen.Edge, 0, len(m.Type.Edges))
	for _, e := range m.Type.Edges {
		if e.Type.IsEdgeSchema() || e.Immutable || e.Annotations["Skip"] == entgql.SkipMutationUpdateInput {
			continue
		}
		inputEdges = append(inputEdges, e)
	}
	return inputEdges
}

func CustomFieldName(val string) string {
	val = strings.TrimSuffix(val, "_id")
	val = strings.TrimSuffix(val, "_ids")
	return val
}
