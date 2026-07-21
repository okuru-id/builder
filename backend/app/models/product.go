package models

import "github.com/goravel/framework/database/orm"

type Product struct {
	orm.Model
	UserID uint `gorm:"index" json:"-"`
	Slug        string  `gorm:"uniqueIndex" json:"slug"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Price       int     `json:"price"`
	Type        string  `json:"type"`
	FilePath    *string `json:"file_path"`
	Thumbnail   *string `json:"thumbnail"`
	Status      string  `json:"status"`
}
