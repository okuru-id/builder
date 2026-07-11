package models

import (
	"github.com/goravel/framework/database/orm"
	"gorm.io/datatypes"
)

// LandingComponent is a reusable component master. Instances reference it by ID
// and may carry per-instance overrides.
type LandingComponent struct {
	orm.Model
	Name string         `json:"name"`
	Tree datatypes.JSON `gorm:"type:jsonb" json:"tree"`
}
