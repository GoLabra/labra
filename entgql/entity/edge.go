package entity

type Edge struct {
	Name             string
	EntName          string
	Caption          string
	Type             string
	BelongsToCaption *string
	Required         *bool
	RelationType     RelationType
	Private          *bool
	Ref              string
}

func (e Edge) Validate() error {
	return nil
}

func (f Edge) Equals(uniqueInput EdgeWhereUniqueInput) bool {
	if uniqueInput.Caption != nil && f.Caption == *uniqueInput.Caption {
		return true
	}

	if uniqueInput.Name != nil && f.Name == *uniqueInput.Name {
		return true
	}

	return false
}

func (e *Edge) ApplyUpdateInput(data UpdateEdgeInput) {
	if data.Caption != nil {
		e.Caption = *data.Caption
	}

	if data.Required != nil {
		e.Required = data.Required
	}

	if data.Private != nil {
		e.Private = data.Private
	}
}
