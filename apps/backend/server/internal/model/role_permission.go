package model

type RoleSpacePermissions struct {
	RoleID       uint `gorm:"type:bigint;primaryKey"`
	PermissionID uint `gorm:"type:bigint;primaryKey"`
}

func (RoleSpacePermissions) TableName() string { return "role_space_permissions" }

type RoleDocumentPermissions struct {
	RoleID       uint `gorm:"type:bigint;primaryKey"`
	PermissionID uint `gorm:"type:bigint;primaryKey"`
}

func (RoleDocumentPermissions) TableName() string { return "role_document_permissions" }
