package annotations

import (
	"github.com/GoLabra/labrago/src/api/entgql/entity"

	"entgo.io/ent/entc"
)

type Field struct {
	Caption        string
	Type           entity.FieldType
	Min            string
	Max            string
	Private        bool
	DefaultValue   string
	AcceptedValues []string
}

func (Field) Name() string {
	return "Field"
}

type FieldExtension struct {
	entc.DefaultExtension
	word Field
}

func (s *FieldExtension) Annotations() []entc.Annotation {
	return []entc.Annotation{
		s.word,
	}
}
