package models

import "github.com/goravel/framework/database/orm"

type Setting struct {
	orm.Model
	Key   string `gorm:"uniqueIndex" json:"key"`
	Value string `json:"value"`
}
