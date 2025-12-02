package entity

import (
	"bytes"
	"testing"
)

func TestEntityOwner_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		owner EntityOwner
		want  bool
	}{
		{
			name:  "valid Admin",
			owner: EntityOwnerAdmin,
			want:  true,
		},
		{
			name:  "valid User",
			owner: EntityOwnerUser,
			want:  true,
		},
		{
			name:  "invalid empty",
			owner: EntityOwner(""),
			want:  false,
		},
		{
			name:  "invalid value",
			owner: EntityOwner("Invalid"),
			want:  false,
		},
		{
			name:  "invalid lowercase admin",
			owner: EntityOwner("admin"),
			want:  false,
		},
		{
			name:  "invalid lowercase user",
			owner: EntityOwner("user"),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.owner.IsValid(); got != tt.want {
				t.Errorf("EntityOwner.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntityOwner_String(t *testing.T) {
	tests := []struct {
		name  string
		owner EntityOwner
		want  string
	}{
		{
			name:  "Admin",
			owner: EntityOwnerAdmin,
			want:  "Admin",
		},
		{
			name:  "User",
			owner: EntityOwnerUser,
			want:  "User",
		},
		{
			name:  "empty",
			owner: EntityOwner(""),
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.owner.String(); got != tt.want {
				t.Errorf("EntityOwner.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntityOwner_UnmarshalGQL(t *testing.T) {
	tests := []struct {
		name    string
		v       interface{}
		wantErr bool
		want    EntityOwner
	}{
		{
			name:    "valid Admin string",
			v:       "Admin",
			wantErr: false,
			want:    EntityOwnerAdmin,
		},
		{
			name:    "valid User string",
			v:       "User",
			wantErr: false,
			want:    EntityOwnerUser,
		},
		{
			name:    "invalid type - int",
			v:       123,
			wantErr: true,
		},
		{
			name:    "invalid type - bool",
			v:       true,
			wantErr: true,
		},
		{
			name:    "invalid value - empty string",
			v:       "",
			wantErr: true,
		},
		{
			name:    "invalid value - invalid string",
			v:       "Invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var owner EntityOwner
			err := owner.UnmarshalGQL(tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("EntityOwner.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && owner != tt.want {
				t.Errorf("EntityOwner.UnmarshalGQL() = %v, want %v", owner, tt.want)
			}
		})
	}
}

func TestEntityOwner_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		owner EntityOwner
		want  string
	}{
		{
			name:  "Admin",
			owner: EntityOwnerAdmin,
			want:  `"Admin"`,
		},
		{
			name:  "User",
			owner: EntityOwnerUser,
			want:  `"User"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.owner.MarshalGQL(&buf)
			if got := buf.String(); got != tt.want {
				t.Errorf("EntityOwner.MarshalGQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelationType_IsValid(t *testing.T) {
	tests := []struct {
		name         string
		relationType RelationType
		want         bool
	}{
		{
			name:         "valid One",
			relationType: RelationTypeOne,
			want:         true,
		},
		{
			name:         "valid Many",
			relationType: RelationTypeMany,
			want:         true,
		},
		{
			name:         "valid OneToOne",
			relationType: RelationTypeO2O,
			want:         true,
		},
		{
			name:         "valid OneToMany",
			relationType: RelationTypeO2M,
			want:         true,
		},
		{
			name:         "valid ManyToOne",
			relationType: RelationTypeM2O,
			want:         true,
		},
		{
			name:         "valid ManyToMany",
			relationType: RelationTypeM2M,
			want:         true,
		},
		{
			name:         "invalid empty",
			relationType: RelationType(""),
			want:         false,
		},
		{
			name:         "invalid value",
			relationType: RelationType("Invalid"),
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.relationType.IsValid(); got != tt.want {
				t.Errorf("RelationType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelationType_String(t *testing.T) {
	tests := []struct {
		name         string
		relationType RelationType
		want         string
	}{
		{
			name:         "One",
			relationType: RelationTypeOne,
			want:         "One",
		},
		{
			name:         "Many",
			relationType: RelationTypeMany,
			want:         "Many",
		},
		{
			name:         "OneToOne",
			relationType: RelationTypeO2O,
			want:         "OneToOne",
		},
		{
			name:         "OneToMany",
			relationType: RelationTypeO2M,
			want:         "OneToMany",
		},
		{
			name:         "ManyToOne",
			relationType: RelationTypeM2O,
			want:         "ManyToOne",
		},
		{
			name:         "ManyToMany",
			relationType: RelationTypeM2M,
			want:         "ManyToMany",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.relationType.String(); got != tt.want {
				t.Errorf("RelationType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelationType_VariableName(t *testing.T) {
	tests := []struct {
		name         string
		relationType RelationType
		want         string
	}{
		{
			name:         "One",
			relationType: RelationTypeOne,
			want:         "One",
		},
		{
			name:         "Many",
			relationType: RelationTypeMany,
			want:         "Many",
		},
		{
			name:         "OneToOne",
			relationType: RelationTypeO2O,
			want:         "O2O",
		},
		{
			name:         "OneToMany",
			relationType: RelationTypeO2M,
			want:         "O2M",
		},
		{
			name:         "ManyToOne",
			relationType: RelationTypeM2O,
			want:         "M2O",
		},
		{
			name:         "ManyToMany",
			relationType: RelationTypeM2M,
			want:         "M2M",
		},
		{
			name:         "invalid (defaults to Many)",
			relationType: RelationType("Invalid"),
			want:         "Many",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.relationType.VariableName(); got != tt.want {
				t.Errorf("RelationType.VariableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelationType_UnmarshalGQL(t *testing.T) {
	tests := []struct {
		name    string
		v       interface{}
		wantErr bool
		want    RelationType
	}{
		{
			name:    "valid One string",
			v:       "One",
			wantErr: false,
			want:    RelationTypeOne,
		},
		{
			name:    "valid Many string",
			v:       "Many",
			wantErr: false,
			want:    RelationTypeMany,
		},
		{
			name:    "valid OneToOne string",
			v:       "OneToOne",
			wantErr: false,
			want:    RelationTypeO2O,
		},
		{
			name:    "valid OneToMany string",
			v:       "OneToMany",
			wantErr: false,
			want:    RelationTypeO2M,
		},
		{
			name:    "valid ManyToOne string",
			v:       "ManyToOne",
			wantErr: false,
			want:    RelationTypeM2O,
		},
		{
			name:    "valid ManyToMany string",
			v:       "ManyToMany",
			wantErr: false,
			want:    RelationTypeM2M,
		},
		{
			name:    "invalid type - int",
			v:       123,
			wantErr: true,
		},
		{
			name:    "invalid type - bool",
			v:       true,
			wantErr: true,
		},
		{
			name:    "invalid value - empty string",
			v:       "",
			wantErr: true,
		},
		{
			name:    "invalid value - invalid string",
			v:       "Invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var relationType RelationType
			err := relationType.UnmarshalGQL(tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("RelationType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && relationType != tt.want {
				t.Errorf("RelationType.UnmarshalGQL() = %v, want %v", relationType, tt.want)
			}
		})
	}
}

func TestRelationType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name         string
		relationType RelationType
		want         string
	}{
		{
			name:         "One",
			relationType: RelationTypeOne,
			want:         `"One"`,
		},
		{
			name:         "Many",
			relationType: RelationTypeMany,
			want:         `"Many"`,
		},
		{
			name:         "OneToOne",
			relationType: RelationTypeO2O,
			want:         `"OneToOne"`,
		},
		{
			name:         "OneToMany",
			relationType: RelationTypeO2M,
			want:         `"OneToMany"`,
		},
		{
			name:         "ManyToOne",
			relationType: RelationTypeM2O,
			want:         `"ManyToOne"`,
		},
		{
			name:         "ManyToMany",
			relationType: RelationTypeM2M,
			want:         `"ManyToMany"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.relationType.MarshalGQL(&buf)
			if got := buf.String(); got != tt.want {
				t.Errorf("RelationType.MarshalGQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions (reused from field_test.go)
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func boolPtrEqual(p1, p2 *bool) bool {
	if p1 == nil && p2 == nil {
		return true
	}
	if p1 == nil || p2 == nil {
		return false
	}
	return *p1 == *p2
}

func stringPtrEqual(p1, p2 *string) bool {
	if p1 == nil && p2 == nil {
		return true
	}
	if p1 == nil || p2 == nil {
		return false
	}
	return *p1 == *p2
}

