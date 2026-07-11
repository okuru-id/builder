package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/goravel/framework/database/orm"
)

type LandingTemplate struct {
	orm.Model
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Preview     string               `json:"preview"`
	Sections    LandingTemplateItems `gorm:"type:text" json:"sections"`
	HTML        string               `gorm:"type:text" json:"html,omitempty"`
}

type LandingTemplateItems []any

func (i LandingTemplateItems) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (i *LandingTemplateItems) Scan(value any) error {
	if value == nil {
		*i = nil
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
	return json.Unmarshal([]byte(raw), i)
}
