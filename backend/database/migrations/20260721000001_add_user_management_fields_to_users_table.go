package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

type M20260721000001AddUserManagementFieldsToUsersTable struct{}

func (r *M20260721000001AddUserManagementFieldsToUsersTable) Signature() string {
	return "20260721000001_add_user_management_fields_to_users_table"
}

func (r *M20260721000001AddUserManagementFieldsToUsersTable) Up() error {
	if facades.Schema().HasTable("users") {
		if !facades.Schema().HasColumn("users", "name") {
			if err := facades.Schema().Table("users", func(table schema.Blueprint) {
				table.String("name").Default("")
			}); err != nil {
				return err
			}
		}
		if !facades.Schema().HasColumn("users", "is_active") {
			if err := facades.Schema().Table("users", func(table schema.Blueprint) {
				table.Boolean("is_active").Default(true)
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *M20260721000001AddUserManagementFieldsToUsersTable) Down() error {
	return nil
}
