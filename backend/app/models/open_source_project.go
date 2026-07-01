package models

import (
	"github.com/goravel/framework/database/orm"
	"gorm.io/datatypes"
)

type OpenSourceProject struct {
	orm.Model
	SortOrder     int               `json:"sort_order"`
	TitleEn       string            `json:"title_en"`
	TitleId       string            `json:"title_id"`
	DescriptionEn *string           `json:"description_en"`
	DescriptionId *string           `json:"description_id"`
	GithubUrl     string            `json:"github_url"`
	Technologies  datatypes.JSONMap `gorm:"type:json" json:"technologies"`
	Stars         int               `json:"stars"`
	License       *string           `json:"license"`
}
