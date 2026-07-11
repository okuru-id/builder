package models

import (
	"github.com/goravel/framework/database/orm"
	"gorm.io/datatypes"
)

// LandingPage is a builder page. Tree is the source-of-truth JSON node tree.
type LandingPage struct {
	orm.Model
	Slug          string         `gorm:"uniqueIndex" json:"slug"`
	Name          string         `json:"name"`
	Status        string         `gorm:"default:draft" json:"status"` // draft | published
	Tree          datatypes.JSON `gorm:"type:jsonb" json:"tree"`
	PublishedTree datatypes.JSON `gorm:"type:jsonb" json:"published_tree,omitempty"`
	PublishedHTML string         `gorm:"type:text" json:"published_html,omitempty"`
	Version       int            `gorm:"default:0" json:"version"`
}
