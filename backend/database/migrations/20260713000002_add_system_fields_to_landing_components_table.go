package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

type M20260713000002AddSystemFieldsToLandingComponentsTable struct{}

func (r *M20260713000002AddSystemFieldsToLandingComponentsTable) Signature() string {
	return "20260713000002_add_system_fields_to_landing_components_table"
}

func (r *M20260713000002AddSystemFieldsToLandingComponentsTable) Up() error {
	if facades.Schema().HasTable("landing_components") {
		if !facades.Schema().HasColumn("landing_components", "key") {
			if err := facades.Schema().Table("landing_components", func(table schema.Blueprint) {
				table.String("key").Nullable()
			}); err != nil {
				return err
			}
		}
		if !facades.Schema().HasColumn("landing_components", "is_system") {
			if err := facades.Schema().Table("landing_components", func(table schema.Blueprint) {
				table.Boolean("is_system").Default(false)
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *M20260713000002AddSystemFieldsToLandingComponentsTable) Down() error {
	return nil
}
