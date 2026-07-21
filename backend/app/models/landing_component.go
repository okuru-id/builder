package models

import (
	"github.com/goravel/framework/database/orm"
	"gorm.io/datatypes"
)

// LandingComponent is a reusable component master. Instances reference it by ID
// and may carry per-instance overrides.
type LandingComponent struct {
	orm.Model
	UserID uint `gorm:"index" json:"-"`
	Key      string         `gorm:"uniqueIndex" json:"key,omitempty"`
	Name     string         `json:"name"`
	IsSystem bool           `gorm:"default:false" json:"is_system"`
	Tree     datatypes.JSON `gorm:"type:jsonb" json:"tree"`
}
