package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

type M20260705000001CreateLandingSectionsTable struct{}

func (r *M20260705000001CreateLandingSectionsTable) Signature() string {
	return "20260705000001_create_landing_sections_table"
}

func (r *M20260705000001CreateLandingSectionsTable) Up() error {
	if !facades.Schema().HasTable("landing_sections") {
		if err := facades.Schema().Create("landing_sections", func(table schema.Blueprint) {
			table.ID()
			table.String("type")
			table.Unique("type")
			table.Text("content").Nullable()
			table.Integer("sort_order").Default(0)
			table.Boolean("is_active").Default(true)
			table.Timestamps()
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260705000001CreateLandingSectionsTable) Down() error {
	return facades.Schema().DropIfExists("landing_sections")
}
