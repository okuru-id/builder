package models

import (
	"github.com/goravel/framework/database/orm"
	"gorm.io/datatypes"
)

// LandingPageRevision is an append-only snapshot of a page tree for undo/redo and rollback.
type LandingPageRevision struct {
	orm.Model
	LandingPageID uint           `gorm:"index" json:"landing_page_id"`
	Tree          datatypes.JSON `gorm:"type:jsonb" json:"tree"`
	Message       string         `json:"message,omitempty"`
}
