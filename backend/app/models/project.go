package models

import (
	"github.com/goravel/framework/database/orm"
	"gorm.io/datatypes"
)

type Project struct {
	orm.Model
	UserID    uint `gorm:"index" json:"-"`
	SortOrder     int               `json:"sort_order"`
	TitleEn       string            `json:"title_en"`
	TitleId       string            `json:"title_id"`
	DescriptionEn *string           `json:"description_en"`
	DescriptionId *string           `json:"description_id"`
	TechStack     datatypes.JSONMap `gorm:"type:json" json:"tech_stack"`
	Url           *string           `json:"url"`
	Featured      bool              `json:"featured"`
}
