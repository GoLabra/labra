package annotations

import (
	"github.com/GoLabra/labrago/src/api/entgql/entity"

	"entgo.io/ent/entc"
)

type Edge struct {
	Caption      string
	RelationType entity.RelationType
}

func (Edge) Name() string {
	return "Edge"
}

type EdgeExtension struct {
	entc.DefaultExtension
	word Edge
}

func (s *EdgeExtension) Annotations() []entc.Annotation {
	return []entc.Annotation{
		s.word,
	}
}
