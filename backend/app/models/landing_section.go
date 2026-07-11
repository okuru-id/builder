package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/goravel/framework/database/orm"
)

type LandingSection struct {
	orm.Model
	Type      string          `json:"type"`
	Content   LandingContent  `gorm:"type:text" json:"content"`
	SortOrder int             `gorm:"default:0" json:"sort_order"`
	IsActive  bool            `gorm:"default:true" json:"is_active"`
}

type LandingContent map[string]any

func (c LandingContent) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (c *LandingContent) Scan(value any) error {
	if value == nil {
		*c = nil
		return nil
	}
	var raw string
	switch v := value.(type) {
	case []byte:
		raw = string(v)
	case string:
		raw = v
	default:
		return nil
	}
	return json.Unmarshal([]byte(raw), c)
}
