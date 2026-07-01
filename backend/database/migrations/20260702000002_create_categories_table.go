package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

type M20260702000002CreateCategoriesTable struct{}

// Signature The unique signature for the migration.
func (r *M20260702000002CreateCategoriesTable) Signature() string {
	return "20260702000002_create_categories_table"
}

// Up Run the migrations.
func (r *M20260702000002CreateCategoriesTable) Up() error {
	if !facades.Schema().HasTable("categories") {
		if err := facades.Schema().Create("categories", func(table schema.Blueprint) {
			table.ID()
			table.String("slug")
			table.Unique("slug")
			table.String("name_en")
			table.String("name_id")
			table.Timestamps()
		}); err != nil {
			return err
		}
	}
	return nil
}

// Down Reverse the migrations.
func (r *M20260702000002CreateCategoriesTable) Down() error {
	if err := facades.Schema().DropIfExists("categories"); err != nil {
		return err
	}
	return nil
}
