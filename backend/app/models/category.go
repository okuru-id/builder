package models

import "github.com/goravel/framework/database/orm"

type Category struct {
	orm.Model
	UserID uint `gorm:"index" json:"-"`
	Slug   string `gorm:"uniqueIndex" json:"slug"`
	NameEn string `json:"name_en"`
	NameId string `json:"name_id"`
}
