package models

import "github.com/goravel/framework/database/orm"

type Message struct {
	orm.Model
	UserID uint `gorm:"index" json:"-"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Content string `gorm:"column:message" json:"message"`
	Status  string `json:"status"`
}
