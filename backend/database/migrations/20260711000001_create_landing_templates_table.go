package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

type M20260711000001CreateLandingTemplatesTable struct{}

func (r *M20260711000001CreateLandingTemplatesTable) Signature() string {
	return "20260711000001_create_landing_templates_table"
}

func (r *M20260711000001CreateLandingTemplatesTable) Up() error {
	if !facades.Schema().HasTable("landing_templates") {
		return facades.Schema().Create("landing_templates", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.String("description").Default("")
			table.String("preview").Default("")
			table.Text("sections").Nullable()
			table.Timestamps()
		})
	}
	return nil
}

func (r *M20260711000001CreateLandingTemplatesTable) Down() error {
	return facades.Schema().DropIfExists("landing_templates")
}
