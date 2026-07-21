package models

import "github.com/goravel/framework/database/orm"

type User struct {
	orm.Model
	Email        string  `gorm:"uniqueIndex" json:"email"`
	Name         string  `json:"name"`
	Password     string  `json:"-"`
	TotpSecret   *string `json:"-"`
	TotpVerified bool    `json:"-"`
	IsActive     bool    `gorm:"default:true" json:"is_active"`
	IsAdmin      bool    `gorm:"default:false" json:"is_admin"`
}
