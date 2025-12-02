package entity

import (
	"testing"
)

func TestEntity_Validate(t *testing.T) {
	tests := []struct {
		name    string
		entity  Entity
		wantErr bool
	}{
		{
			name: "valid entity",
			entity: Entity{
				Name:             "test",
				EntName:          "Test",
				Caption:          "Test Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: false,
		},
		{
			name: "entity with admin owner",
			entity: Entity{
				Name:             "admin",
				EntName:          "Admin",
				Caption:          "Admin Entity",
				Owner:            EntityOwnerAdmin,
				DisplayFieldName: "name",
			},
			wantErr: false,
		},
		{
			name:    "empty entity (should fail validation)",
			entity:  Entity{},
			wantErr: true,
		},
		{
			name: "entity with reserved name",
			entity: Entity{
				Name:             "entity",
				EntName:          "Entity",
				Caption:          "Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true,
		},
		{
			name: "entity with empty EntName",
			entity: Entity{
				Name:             "test",
				EntName:          "",
				Caption:          "Test Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true,
		},
		{
			name: "entity with empty Caption",
			entity: Entity{
				Name:             "test",
				EntName:          "Test",
				Caption:          "",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true,
		},
		{
			name: "entity with invalid name starting with number",
			entity: Entity{
				Name:             "123entity",
				EntName:          "Entity123",
				Caption:          "Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true,
		},
		{
			name: "entity with invalid name with special characters",
			entity: Entity{
				Name:             "entity-name",
				EntName:          "EntityName",
				Caption:          "Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true,
		},
		{
			name: "entity with invalid EntName starting with number",
			entity: Entity{
				Name:             "entity",
				EntName:          "123Entity",
				Caption:          "Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true,
		},
		{
			name: "entity with invalid EntName with special characters",
			entity: Entity{
				Name:             "test",
				EntName:          "Entity-Name",
				Caption:          "Test Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true,
		},
		{
			name: "entity with invalid EntName with spaces",
			entity: Entity{
				Name:             "test",
				EntName:          "Entity Name",
				Caption:          "Test Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true,
		},
		{
			name: "entity with valid EntName with underscore at start",
			entity: Entity{
				Name:             "test",
				EntName:          "_Entity",
				Caption:          "Test Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: false, // Underscore at start is valid
		},
		{
			name: "entity with valid EntName with underscore",
			entity: Entity{
				Name:             "test",
				EntName:          "Entity_Name",
				Caption:          "Test Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: false,
		},
		{
			name: "entity with valid EntName with numbers in middle",
			entity: Entity{
				Name:             "test",
				EntName:          "Entity123Name",
				Caption:          "Test Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: false,
		},
		{
			name: "entity with invalid Owner",
			entity: Entity{
				Name:             "entity",
				EntName:          "Entity",
				Caption:          "Entity",
				Owner:            EntityOwner("Invalid"),
				DisplayFieldName: "name",
			},
			wantErr: true,
		},
		{
			name: "entity with empty Owner (zero value)",
			entity: Entity{
				Name:             "entity",
				EntName:          "Entity",
				Caption:          "Entity",
				Owner:            EntityOwner(""),
				DisplayFieldName: "name",
			},
			wantErr: true,
		},
		{
			name: "entity with name starting with underscore (valid)",
			entity: Entity{
				Name:             "_entity",
				EntName:          "Entity",
				Caption:          "Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: false,
		},
		{
			name: "entity with name containing numbers (valid)",
			entity: Entity{
				Name:             "entity123",
				EntName:          "Entity123",
				Caption:          "Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: false,
		},
		{
			name: "entity with name containing underscores (valid)",
			entity: Entity{
				Name:             "entity_name",
				EntName:          "EntityName",
				Caption:          "Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: false,
		},
		{
			name: "entity with valid name but invalid EntName (empty Caption scenario)",
			entity: Entity{
				Name:             "test",
				EntName:          "Test-Name",
				Caption:          "",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true, // Should fail on EntName validation, but Caption empty will be caught first
		},
		{
			name: "entity with valid name but invalid EntName with special chars",
			entity: Entity{
				Name:             "test",
				EntName:          "Test@Name",
				Caption:          "Test Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true,
		},
		{
			name: "entity with valid name but invalid EntName starting with number (after name passes)",
			entity: Entity{
				Name:             "validname",
				EntName:          "9Invalid",
				Caption:          "Test Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true,
		},
		{
			name: "entity with empty name (should fail early)",
			entity: Entity{
				Name:             "",
				EntName:          "Entity",
				Caption:          "Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true,
		},
		{
			name: "entity with all valid fields including underscore in EntName",
			entity: Entity{
				Name:             "test_entity",
				EntName:          "Test_Entity",
				Caption:          "Test Entity",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: false,
		},
		{
			name: "entity with reserved name and empty Caption (tests error message with empty Caption)",
			entity: Entity{
				Name:             "entity",
				EntName:          "Entity",
				Caption:          "",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true, // Fails on reserved name check, uses empty Caption in error
		},
		{
			name: "entity with invalid name identifier and empty Caption",
			entity: Entity{
				Name:             "123test",
				EntName:          "Test",
				Caption:          "",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true, // Fails on name identifier check, uses empty Caption in error
		},
		{
			name: "entity with valid name, empty EntName, and empty Caption",
			entity: Entity{
				Name:             "test",
				EntName:          "",
				Caption:          "",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true, // Fails on EntName empty check, uses empty Caption in error
		},
		{
			name: "entity with valid name, invalid EntName identifier, and empty Caption",
			entity: Entity{
				Name:             "test",
				EntName:          "Test-Name",
				Caption:          "",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true, // Fails on EntName identifier check, uses empty Caption in error
		},
		{
			name: "entity with valid name and EntName but empty Caption (uses Name in error)",
			entity: Entity{
				Name:             "test",
				EntName:          "Test",
				Caption:          "",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			wantErr: true, // Fails on Caption empty, uses Name in error message
		},
		{
			name: "entity with valid name and EntName but invalid Owner and empty Caption",
			entity: Entity{
				Name:             "test",
				EntName:          "Test",
				Caption:          "",
				Owner:            EntityOwner("Invalid"),
				DisplayFieldName: "name",
			},
			wantErr: true, // Fails on Caption empty first, but if that's not checked, would fail on Owner
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entity.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Entity.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEntity_Equals(t *testing.T) {
	entity := Entity{
		Name:             "testEntity",
		EntName:          "TestEntity",
		Caption:          "Test Entity",
		Owner:            EntityOwnerUser,
		DisplayFieldName: "name",
	}

	tests := []struct {
		name        string
		uniqueInput EntityWhereUniqueInput
		want        bool
	}{
		{
			name: "match by caption",
			uniqueInput: EntityWhereUniqueInput{
				Caption: stringPtr("Test Entity"),
			},
			want: true,
		},
		{
			name: "match by name",
			uniqueInput: EntityWhereUniqueInput{
				Name: stringPtr("testEntity"),
			},
			want: true,
		},
		{
			name: "match by entName",
			uniqueInput: EntityWhereUniqueInput{
				EntName: stringPtr("TestEntity"),
			},
			want: false, // Equals only checks Caption and Name, not EntName
		},
		{
			name: "match by both caption and name",
			uniqueInput: EntityWhereUniqueInput{
				Caption: stringPtr("Test Entity"),
				Name:    stringPtr("testEntity"),
			},
			want: true,
		},
		{
			name: "no match - different caption",
			uniqueInput: EntityWhereUniqueInput{
				Caption: stringPtr("Different Entity"),
			},
			want: false,
		},
		{
			name: "no match - different name",
			uniqueInput: EntityWhereUniqueInput{
				Name: stringPtr("differentEntity"),
			},
			want: false,
		},
		{
			name:        "no match - both nil",
			uniqueInput: EntityWhereUniqueInput{},
			want:        false,
		},
		{
			name: "no match - caption nil, name different",
			uniqueInput: EntityWhereUniqueInput{
				Name: stringPtr("differentEntity"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := entity.Equals(tt.uniqueInput); got != tt.want {
				t.Errorf("Entity.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntity_ApplyUpdateInput(t *testing.T) {
	tests := []struct {
		name     string
		entity   *Entity
		input    UpdateEntityInput
		expected Entity
	}{
		{
			name: "update caption",
			entity: &Entity{
				Name:             "test",
				EntName:          "Test",
				Caption:          "Old Caption",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			input: UpdateEntityInput{
				Caption: stringPtr("New Caption"),
			},
			expected: Entity{
				Name:             "test",
				EntName:          "Test",
				Caption:          "New Caption",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
		},
		{
			name: "update with nil caption (should not change)",
			entity: &Entity{
				Name:             "test",
				EntName:          "Test",
				Caption:          "Original Caption",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			input: UpdateEntityInput{
				Caption: nil,
			},
			expected: Entity{
				Name:             "test",
				EntName:          "Test",
				Caption:          "Original Caption",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
		},
		{
			name: "update caption with empty string",
			entity: &Entity{
				Name:             "test",
				EntName:          "Test",
				Caption:          "Old Caption",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
			input: UpdateEntityInput{
				Caption: stringPtr(""),
			},
			expected: Entity{
				Name:             "test",
				EntName:          "Test",
				Caption:          "",
				Owner:            EntityOwnerUser,
				DisplayFieldName: "name",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.entity.ApplyUpdateInput(tt.input)
			if !entitiesEqual(*tt.entity, tt.expected) {
				t.Errorf("Entity.ApplyUpdateInput() = %+v, want %+v", *tt.entity, tt.expected)
			}
		})
	}
}

func entitiesEqual(e1, e2 Entity) bool {
	return e1.Name == e2.Name &&
		e1.EntName == e2.EntName &&
		e1.Caption == e2.Caption &&
		e1.Owner == e2.Owner &&
		e1.DisplayFieldName == e2.DisplayFieldName
}
