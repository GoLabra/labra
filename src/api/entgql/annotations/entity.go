package annotations

import (
	"github.com/GoLabra/labrago/src/api/entgql/entity"

	"entgo.io/ent/entc"
)

const EntityName = "Entity"

type Entity struct {
	Caption      string
	DisplayField string
	Owner        entity.EntityOwner
}

func (Entity) Name() string {
	return EntityName
}

type EntityExtension struct {
	entc.DefaultExtension
	word Entity
}

func (s *EntityExtension) Annotations() []entc.Annotation {
	return []entc.Annotation{
		s.word,
	}
}
