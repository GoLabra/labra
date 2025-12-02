package entity

import (
	"testing"
)

func TestEdge_Validate(t *testing.T) {
	tests := []struct {
		name    string
		edge    Edge
		wantErr bool
	}{
		{
			name: "valid edge",
			edge: Edge{
				Name:         "test",
				EntName:      "Test",
				Caption:      "Test Edge",
				Type:         "OneToMany",
				RelationType: RelationTypeO2M,
				Required:     boolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "edge with belongs to",
			edge: Edge{
				Name:             "parent",
				EntName:          "Parent",
				Caption:          "Parent",
				Type:             "ManyToOne",
				RelationType:     RelationTypeM2O,
				BelongsToCaption: stringPtr("Child"),
			},
			wantErr: false,
		},
		{
			name:    "empty edge (should fail validation)",
			edge:    Edge{},
			wantErr: true,
		},
		{
			name: "edge with reserved name",
			edge: Edge{
				Name:         "id",
				EntName:      "ID",
				Caption:      "ID",
				Type:         "OneToMany",
				RelationType: RelationTypeO2M,
			},
			wantErr: true,
		},
		{
			name: "edge with empty Caption",
			edge: Edge{
				Name:         "test",
				EntName:      "Test",
				Caption:      "",
				Type:         "OneToMany",
				RelationType: RelationTypeO2M,
			},
			wantErr: true,
		},
		{
			name: "edge with invalid name starting with number",
			edge: Edge{
				Name:         "123edge",
				EntName:      "Edge123",
				Caption:      "Edge",
				Type:         "OneToMany",
				RelationType: RelationTypeO2M,
			},
			wantErr: true,
		},
		{
			name: "edge with invalid name with special characters",
			edge: Edge{
				Name:         "edge-name",
				EntName:      "EdgeName",
				Caption:      "Edge",
				Type:         "OneToMany",
				RelationType: RelationTypeO2M,
			},
			wantErr: true,
		},
		{
			name: "edge with invalid RelationType",
			edge: Edge{
				Name:         "test",
				EntName:      "Test",
				Caption:      "Test Edge",
				Type:         "OneToMany",
				RelationType: RelationType("Invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.edge.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Edge.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEdge_Equals(t *testing.T) {
	edge := Edge{
		Name:         "testEdge",
		EntName:      "TestEdge",
		Caption:      "Test Edge",
		Type:         "OneToMany",
		RelationType: RelationTypeO2M,
	}

	tests := []struct {
		name        string
		uniqueInput EdgeWhereUniqueInput
		want        bool
	}{
		{
			name: "match by caption",
			uniqueInput: EdgeWhereUniqueInput{
				Caption: stringPtr("Test Edge"),
			},
			want: true,
		},
		{
			name: "match by name",
			uniqueInput: EdgeWhereUniqueInput{
				Name: stringPtr("testEdge"),
			},
			want: true,
		},
		{
			name: "match by both caption and name",
			uniqueInput: EdgeWhereUniqueInput{
				Caption: stringPtr("Test Edge"),
				Name:    stringPtr("testEdge"),
			},
			want: true,
		},
		{
			name: "no match - different caption",
			uniqueInput: EdgeWhereUniqueInput{
				Caption: stringPtr("Different Edge"),
			},
			want: false,
		},
		{
			name: "no match - different name",
			uniqueInput: EdgeWhereUniqueInput{
				Name: stringPtr("differentEdge"),
			},
			want: false,
		},
		{
			name:        "no match - both nil",
			uniqueInput: EdgeWhereUniqueInput{},
			want:        false,
		},
		{
			name: "no match - caption nil, name different",
			uniqueInput: EdgeWhereUniqueInput{
				Name: stringPtr("differentEdge"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := edge.Equals(tt.uniqueInput); got != tt.want {
				t.Errorf("Edge.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEdge_ApplyUpdateInput(t *testing.T) {
	tests := []struct {
		name     string
		edge     *Edge
		input    UpdateEdgeInput
		expected Edge
	}{
		{
			name: "update caption",
			edge: &Edge{
				Name:         "test",
				EntName:      "Test",
				Caption:      "Old Caption",
				Type:         "OneToMany",
				RelationType: RelationTypeO2M,
			},
			input: UpdateEdgeInput{
				Caption: stringPtr("New Caption"),
			},
			expected: Edge{
				Name:         "test",
				EntName:      "Test",
				Caption:      "New Caption",
				Type:         "OneToMany",
				RelationType: RelationTypeO2M,
			},
		},
		{
			name: "update required",
			edge: &Edge{
				Name:         "test",
				EntName:      "Test",
				Caption:      "Test",
				Type:         "OneToMany",
				RelationType: RelationTypeO2M,
				Required:     boolPtr(false),
			},
			input: UpdateEdgeInput{
				Required: boolPtr(true),
			},
			expected: Edge{
				Name:         "test",
				EntName:      "Test",
				Caption:      "Test",
				Type:         "OneToMany",
				RelationType: RelationTypeO2M,
				Required:     boolPtr(true),
			},
		},
		{
			name: "update private",
			edge: &Edge{
				Name:         "test",
				EntName:      "Test",
				Caption:      "Test",
				Type:         "OneToMany",
				RelationType: RelationTypeO2M,
				Private:      boolPtr(false),
			},
			input: UpdateEdgeInput{
				Private: boolPtr(true),
			},
			expected: Edge{
				Name:         "test",
				EntName:      "Test",
				Caption:      "Test",
				Type:         "OneToMany",
				RelationType: RelationTypeO2M,
				Private:      boolPtr(true),
			},
		},
		{
			name: "update multiple fields",
			edge: &Edge{
				Name:         "test",
				EntName:      "Test",
				Caption:      "Old Caption",
				Type:         "OneToMany",
				RelationType: RelationTypeO2M,
				Required:     boolPtr(false),
				Private:      boolPtr(false),
			},
			input: UpdateEdgeInput{
				Caption:  stringPtr("New Caption"),
				Required: boolPtr(true),
				Private:  boolPtr(true),
			},
			expected: Edge{
				Name:         "test",
				EntName:      "Test",
				Caption:      "New Caption",
				Type:         "OneToMany",
				RelationType: RelationTypeO2M,
				Required:     boolPtr(true),
				Private:      boolPtr(true),
			},
		},
		{
			name: "update with nil values (should not change)",
			edge: &Edge{
				Name:         "test",
				EntName:      "Test",
				Caption:      "Test",
				Type:         "OneToMany",
				RelationType: RelationTypeO2M,
				Required:     boolPtr(true),
			},
			input: UpdateEdgeInput{
				Caption:  nil,
				Required: nil,
			},
			expected: Edge{
				Name:         "test",
				EntName:      "Test",
				Caption:      "Test",
				Type:         "OneToMany",
				RelationType: RelationTypeO2M,
				Required:     boolPtr(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.edge.ApplyUpdateInput(tt.input)
			if !edgesEqual(*tt.edge, tt.expected) {
				t.Errorf("Edge.ApplyUpdateInput() = %+v, want %+v", *tt.edge, tt.expected)
			}
		})
	}
}

func edgesEqual(e1, e2 Edge) bool {
	if e1.Name != e2.Name || e1.EntName != e2.EntName || e1.Caption != e2.Caption ||
		e1.Type != e2.Type || e1.RelationType != e2.RelationType || e1.Ref != e2.Ref {
		return false
	}
	if !boolPtrEqual(e1.Required, e2.Required) || !boolPtrEqual(e1.Private, e2.Private) {
		return false
	}
	if !stringPtrEqual(e1.BelongsToCaption, e2.BelongsToCaption) {
		return false
	}
	return true
}
