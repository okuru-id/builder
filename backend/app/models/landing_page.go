package models

import (
	"github.com/goravel/framework/database/orm"
	"gorm.io/datatypes"
)

// LandingPage is a builder page. Tree is the source-of-truth JSON node tree.
// Publishing is multi-page: a page is reachable by a custom domain (Domain),
// by a public path (Path, e.g. "promo" → /promo), or at "/" when IsHome=true.
type LandingPage struct {
	orm.Model
	Slug          string         `gorm:"uniqueIndex" json:"slug"`
	Name          string         `json:"name"`
	Status        string         `gorm:"default:draft" json:"status"` // draft | published
	Tree          datatypes.JSON `gorm:"type:jsonb" json:"tree"`
	PublishedTree datatypes.JSON `gorm:"type:jsonb" json:"published_tree,omitempty"`
	PublishedHTML string         `gorm:"type:text" json:"published_html,omitempty"`
	Version       int            `gorm:"default:0" json:"version"`
	// Public routing. At most one page may have IsHome=true; enforced by the controller.
	Path   string `gorm:"index" json:"path,omitempty"`     // e.g. "promo" → served at /promo
	Domain string `gorm:"index" json:"domain,omitempty"`   // e.g. "client.com" → served when Host matches
	IsHome bool   `gorm:"default:false" json:"is_home"` // when true + published, served at /
}
