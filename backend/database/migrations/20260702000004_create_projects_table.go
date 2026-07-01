package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

type M20260702000004CreateProjectsTable struct{}

// Signature The unique signature for the migration.
func (r *M20260702000004CreateProjectsTable) Signature() string {
	return "20260702000004_create_projects_table"
}

// Up Run the migrations.
func (r *M20260702000004CreateProjectsTable) Up() error {
	if !facades.Schema().HasTable("projects") {
		if err := facades.Schema().Create("projects", func(table schema.Blueprint) {
			table.ID()
			table.Integer("sort_order").Default(0)
			table.String("title_en")
			table.String("title_id")
			table.Text("description_en").Nullable()
			table.Text("description_id").Nullable()
			table.Text("tech_stack").Nullable()
			table.String("url").Nullable()
			table.Boolean("featured").Default(false)
			table.Timestamps()
		}); err != nil {
			return err
		}
	}
	return nil
}

// Down Reverse the migrations.
func (r *M20260702000004CreateProjectsTable) Down() error {
	if err := facades.Schema().DropIfExists("projects"); err != nil {
		return err
	}
	return nil
}
