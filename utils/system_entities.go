package utils

import (
	"github.com/GoLabra/labra/cache"
	"github.com/GoLabra/labra/entgql/entity"
)

// LoadSystemEntities loads hardcoded system entities into the cache.
// This avoids the need to load schema files from the labra module.
func LoadSystemEntities() {
	// Helper functions
	boolPtr := func(b bool) *bool { return &b }

	// User entity
	cache.Entity.Set("user", entity.Entity{
		Name:             "user",
		EntName:          "User",
		Caption:          "User",
		Owner:            entity.EntityOwnerUser,
		DisplayFieldName: "email",
	})

	cache.Field.Set("user", []entity.Field{
		{Caption: "Id", Name: "id", Type: string(entity.FieldTypeID)},
		{Caption: "Email", Name: "email", EntName: "Email", Type: string(entity.FieldTypeEmail), Required: boolPtr(true), Unique: boolPtr(true)},
		{Caption: "Password", Name: "password", EntName: "Password", Type: string(entity.FieldTypeShortText), Required: boolPtr(true), Private: boolPtr(true)},
	})

	cache.Edge.Set("user", []entity.Edge{
		{Name: "createdBy", EntName: "created_by", Caption: "Created By", Type: "User", RelationType: entity.RelationTypeOne},
		{Name: "updatedBy", EntName: "updated_by", Caption: "Updated By", Type: "User", RelationType: entity.RelationTypeOne},
		{Name: "adminCreatedBy", EntName: "admin_created_by", Caption: "Admin Created By", Type: "AdminUser", RelationType: entity.RelationTypeOne},
		{Name: "adminUpdatedBy", EntName: "admin_updated_by", Caption: "Admin Updated By", Type: "AdminUser", RelationType: entity.RelationTypeOne},
		{Name: "roles", EntName: "roles", Caption: "Roles", Type: "Role", RelationType: entity.RelationTypeM2M},
		{Name: "defaultRole", EntName: "default_role", Caption: "Default Role", Type: "Role", RelationType: entity.RelationTypeOne},
	})

	// AdminUser entity
	cache.Entity.Set("adminUser", entity.Entity{
		Name:             "adminUser",
		EntName:          "AdminUser",
		Caption:          "AdminUser",
		Owner:            entity.EntityOwnerAdmin,
		DisplayFieldName: "name",
	})

	defaultNow := "now()"
	cache.Field.Set("adminUser", []entity.Field{
		{Caption: "Id", Name: "id", Type: string(entity.FieldTypeID)},
		{Caption: "Name", Name: "name", EntName: "Name", Type: string(entity.FieldTypeShortText)},
		{Caption: "Email", Name: "email", EntName: "Email", Type: string(entity.FieldTypeEmail), Required: boolPtr(true), Unique: boolPtr(true)},
		{Caption: "Password", Name: "password", EntName: "Password", Type: string(entity.FieldTypeShortText), Required: boolPtr(true), Private: boolPtr(true)},
		{Caption: "First Name", Name: "firstName", EntName: "FirstName", Type: string(entity.FieldTypeShortText), Required: boolPtr(true)},
		{Caption: "Last Name", Name: "lastName", EntName: "LastName", Type: string(entity.FieldTypeShortText), Required: boolPtr(true)},
		{Caption: "Created At", Name: "createdAt", EntName: "CreatedAt", Type: string(entity.FieldTypeDateTime), DefaultValue: &defaultNow, UpdateDefault: false, Nillable: true},
		{Caption: "Updated At", Name: "updatedAt", EntName: "UpdatedAt", Type: string(entity.FieldTypeDateTime), DefaultValue: &defaultNow, UpdateDefault: true, Nillable: true},
	})

	cache.Edge.Set("adminUser", []entity.Edge{
		{Name: "adminCreatedBy", EntName: "admin_created_by", Caption: "Admin Created By", Type: "AdminUser", RelationType: entity.RelationTypeOne},
		{Name: "adminUpdatedBy", EntName: "admin_updated_by", Caption: "Admin Updated By", Type: "AdminUser", RelationType: entity.RelationTypeOne},
		{Name: "roles", EntName: "roles", Caption: "Roles", Type: "Role", RelationType: entity.RelationTypeM2M},
		{Name: "defaultRole", EntName: "default_role", Caption: "Default Role", Type: "Role", RelationType: entity.RelationTypeOne},
	})

	// Role entity
	cache.Entity.Set("role", entity.Entity{
		Name:             "role",
		EntName:          "Role",
		Caption:          "Role",
		Owner:            entity.EntityOwnerAdmin,
		DisplayFieldName: "name",
	})

	cache.Field.Set("role", []entity.Field{
		{Caption: "Id", Name: "id", Type: string(entity.FieldTypeID)},
		{Caption: "Name", Name: "name", EntName: "Name", Type: string(entity.FieldTypeShortText), Required: boolPtr(true), Unique: boolPtr(true)},
		{Caption: "Created At", Name: "createdAt", EntName: "CreatedAt", Type: string(entity.FieldTypeDateTime), DefaultValue: &defaultNow, UpdateDefault: false, Nillable: true},
		{Caption: "Updated At", Name: "updatedAt", EntName: "UpdatedAt", Type: string(entity.FieldTypeDateTime), DefaultValue: &defaultNow, UpdateDefault: true, Nillable: true},
	})

	cache.Edge.Set("role", []entity.Edge{
		{Name: "adminCreatedBy", EntName: "admin_created_by", Caption: "Admin Created By", Type: "AdminUser", RelationType: entity.RelationTypeOne},
		{Name: "adminUpdatedBy", EntName: "admin_updated_by", Caption: "Admin Updated By", Type: "AdminUser", RelationType: entity.RelationTypeOne},
		{Name: "adminUserRoles", EntName: "admin_user_roles", Caption: "Role Admin Users", Type: "AdminUser", RelationType: entity.RelationTypeM2M, Ref: "roles"},
		{Name: "userRoles", EntName: "user_roles", Caption: "Role Users", Type: "User", RelationType: entity.RelationTypeM2M, Ref: "roles"},
		{Name: "permissions", EntName: "permissions", Caption: "Permissions", Type: "Permission", RelationType: entity.RelationTypeO2M, Ref: "role"},
	})

	// Permission entity
	cache.Entity.Set("permission", entity.Entity{
		Name:             "permission",
		EntName:          "Permission",
		Caption:          "Permission",
		Owner:            entity.EntityOwnerAdmin,
		DisplayFieldName: "id",
	})

	cache.Field.Set("permission", []entity.Field{
		{Caption: "Id", Name: "id", Type: string(entity.FieldTypeID)},
		{Caption: "Created At", Name: "createdAt", EntName: "CreatedAt", Type: string(entity.FieldTypeDateTime), DefaultValue: &defaultNow, UpdateDefault: false, Nillable: true},
		{Caption: "Updated At", Name: "updatedAt", EntName: "UpdatedAt", Type: string(entity.FieldTypeDateTime), DefaultValue: &defaultNow, UpdateDefault: true, Nillable: true},
		{Caption: "Entity", Name: "entity", EntName: "Entity", Type: string(entity.FieldTypeShortText)},
		{Caption: "Operation", Name: "operation", EntName: "Operation", Type: string(entity.FieldTypeSingleChoice), AcceptedValues: []string{"Create", "Update", "Delete", "Read", ""}},
	})

	cache.Edge.Set("permission", []entity.Edge{
		{Name: "adminCreatedBy", EntName: "admin_created_by", Caption: "Admin Created By", Type: "AdminUser", RelationType: entity.RelationTypeOne},
		{Name: "adminUpdatedBy", EntName: "admin_updated_by", Caption: "Admin Updated By", Type: "AdminUser", RelationType: entity.RelationTypeOne},
		{Name: "role", EntName: "role", Caption: "Role", Type: "Role", RelationType: entity.RelationTypeM2O},
	})

	// File entity
	cache.Entity.Set("file", entity.Entity{
		Name:             "file",
		EntName:          "File",
		Caption:          "File",
		Owner:            entity.EntityOwnerAdmin,
		DisplayFieldName: "id",
	})

	cache.Field.Set("file", []entity.Field{
		{Caption: "Id", Name: "id", Type: string(entity.FieldTypeID)},
		{Caption: "Created At", Name: "createdAt", EntName: "CreatedAt", Type: string(entity.FieldTypeDateTime), DefaultValue: &defaultNow, UpdateDefault: false, Nillable: true},
		{Caption: "Updated At", Name: "updatedAt", EntName: "UpdatedAt", Type: string(entity.FieldTypeDateTime), DefaultValue: &defaultNow, UpdateDefault: true, Nillable: true},
		{Caption: "Caption", Name: "caption", EntName: "Caption", Type: string(entity.FieldTypeShortText), Unique: boolPtr(true)},
		{Caption: "Name", Name: "name", EntName: "Name", Type: string(entity.FieldTypeShortText)},
		{Caption: "MIME Type", Name: "mimeType", EntName: "MimeType", Type: string(entity.FieldTypeShortText)},
		{Caption: "Storage File Name", Name: "storageFileName", EntName: "StorageFileName", Type: string(entity.FieldTypeShortText)},
		{Caption: "Size", Name: "size", EntName: "Size", Type: string(entity.FieldTypeInteger)},
	})

	cache.Edge.Set("file", []entity.Edge{
		{Name: "adminCreatedBy", EntName: "admin_created_by", Caption: "Admin Created By", Type: "AdminUser", RelationType: entity.RelationTypeOne},
		{Name: "adminUpdatedBy", EntName: "admin_updated_by", Caption: "Admin Updated By", Type: "AdminUser", RelationType: entity.RelationTypeOne},
		{Name: "createdBy", EntName: "created_by", Caption: "Created By", Type: "User", RelationType: entity.RelationTypeOne},
		{Name: "updatedBy", EntName: "updated_by", Caption: "Updated By", Type: "User", RelationType: entity.RelationTypeOne},
	})
}
