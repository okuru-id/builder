package models

import "github.com/goravel/framework/database/orm"

type User struct {
	orm.Model
	Email        string  `gorm:"uniqueIndex" json:"email"`
	Password     string  `json:"-"`
	TotpSecret   *string `json:"-"`
	TotpVerified bool    `json:"-"`
}
