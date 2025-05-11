package entity

type Entity struct {
	Name             string
	EntName          string
	Caption          string
	Owner            EntityOwner
	DisplayFieldName string
}

func (f Entity) Validate() error {
	return nil
}

func (f Entity) Equals(uniqueInput EntityWhereUniqueInput) bool {
	if uniqueInput.Caption != nil && f.Caption == *uniqueInput.Caption {
		return true
	}

	if uniqueInput.Name != nil && f.Name == *uniqueInput.Name {
		return true
	}

	return false
}

func (e *Entity) ApplyUpdateInput(data UpdateEntityInput) {
	if data.Caption != nil {
		e.Caption = *data.Caption
	}
}
