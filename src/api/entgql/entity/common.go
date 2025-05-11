package entity

import (
	"fmt"
	"io"
	"strconv"
)

type CreateEdgeInput struct {
	Caption          string
	RelatedEntity    *EntityConnectInput
	Required         *bool
	RelationType     RelationType
	Private          *bool
	BelongsToCaption *string
	Type             *string
	Ref              string
}

type CreateEntityInput struct {
	Caption      string
	DisplayField FieldWhereUniqueInput
	Fields       *CreateManyFieldsInput
	Edges        *CreateManyEdgesInput
}

type CreateFieldInput struct {
	Caption        string
	Type           string
	Required       *bool
	Unique         *bool
	DefaultValue   *string
	Min            *string
	Max            *string
	Private        *bool
	AcceptedValues []string
}

type CreateManyEdgesInput struct {
	Create []*CreateEdgeInput
}

type CreateManyFieldsInput struct {
	Create []*CreateFieldInput
}

type EdgeWhereUniqueInput struct {
	Name    *string
	Caption *string
}

type EntityConnectInput struct {
	Connect EntityWhereUniqueInput
}

type EntityWhereUniqueInput struct {
	Name    *string
	Caption *string
	EntName *string
}

type FieldWhereUniqueInput struct {
	Name    *string
	Caption *string
}

type UpdateEdgeInput struct {
	Caption  *string
	Required *bool
	Private  *bool
}

type UpdateEntityInput struct {
	Caption      *string
	DisplayField *FieldWhereUniqueInput
	Fields       *UpdateManyFieldsInput
	Edges        *UpdateManyEdgesInput
}

type UpdateFieldInput struct {
	Caption        *string  `json:"caption"`
	Required       *bool    `json:"required"`
	Unique         *bool    `json:"unique"`
	DefaultValue   *string  `json:"defaultValue"`
	Min            *string  `json:"min"`
	Max            *string  `json:"max"`
	Private        *bool    `json:"private"`
	AcceptedValues []string `json:"acceptedValues"`
}

type UpdateManyEdgesInput struct {
	Create []*CreateEdgeInput      `json:"create"`
	Update []*UpdateOneEdgeInput   `json:"update"`
	Delete []*EdgeWhereUniqueInput `json:"delete"`
}

type UpdateOneEdgeInput struct {
	Where EdgeWhereUniqueInput
	Data  UpdateEdgeInput
}

type UpdateManyFieldsInput struct {
	Create []*CreateFieldInput      `json:"create"`
	Update []*UpdateOneFieldInput   `json:"update"`
	Delete []*FieldWhereUniqueInput `json:"delete"`
}

type UpdateOneFieldInput struct {
	Where FieldWhereUniqueInput
	Data  UpdateFieldInput
}

type EntityOwner string

const (
	EntityOwnerAdmin EntityOwner = "Admin"
	EntityOwnerUser  EntityOwner = "User"
)

var AllEntityOwner = []EntityOwner{
	EntityOwnerAdmin,
	EntityOwnerUser,
}

func (e EntityOwner) IsValid() bool {
	switch e {
	case EntityOwnerAdmin, EntityOwnerUser:
		return true
	}
	return false
}

func (e EntityOwner) String() string {
	return string(e)
}

func (e *EntityOwner) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EntityOwner(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EntityOwner", str)
	}
	return nil
}

func (e EntityOwner) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type RelationType string

const (
	RelationTypeOne  RelationType = "One"
	RelationTypeMany RelationType = "Many"
	RelationTypeO2O  RelationType = "OneToOne"
	RelationTypeO2M  RelationType = "OneToMany"
	RelationTypeM2O  RelationType = "ManyToOne"
	RelationTypeM2M  RelationType = "ManyToMany"
)

var AllRelationType = []RelationType{
	RelationTypeOne,
	RelationTypeMany,
	RelationTypeO2O,
	RelationTypeO2M,
	RelationTypeM2O,
	RelationTypeM2M,
}

func (e RelationType) VariableName() string {
	switch e {
	case RelationTypeOne:
		return "One"
	case RelationTypeMany:
		return "Many"
	case RelationTypeO2O:
		return "O2O"
	case RelationTypeO2M:
		return "O2M"
	case RelationTypeM2O:
		return "M2O"
	case RelationTypeM2M:
		return "M2M"
	}
	return "Many"
}

func (e RelationType) IsValid() bool {
	switch e {
	case RelationTypeO2O, RelationTypeO2M, RelationTypeM2O, RelationTypeM2M, RelationTypeOne, RelationTypeMany:
		return true
	}
	return false
}

func (e RelationType) String() string {
	return string(e)
}

func (e *RelationType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RelationType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RelationType", str)
	}
	return nil
}

func (e RelationType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
