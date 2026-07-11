package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

type M20260713000001CreateLandingPagesTables struct{}

func (r *M20260713000001CreateLandingPagesTables) Signature() string {
	return "20260713000001_create_landing_pages_tables"
}

func (r *M20260713000001CreateLandingPagesTables) Up() error {
	if !facades.Schema().HasTable("landing_pages") {
		if err := facades.Schema().Create("landing_pages", func(table schema.Blueprint) {
			table.ID()
			table.String("slug")
			table.Unique("slug")
			table.String("name")
			table.String("status").Default("draft") // draft | published
			table.Json("tree")
			table.Json("published_tree").Nullable()
			table.Text("published_html").Nullable()
			table.Integer("version").Default(0)
			table.Timestamps()
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("landing_page_revisions") {
		if err := facades.Schema().Create("landing_page_revisions", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("landing_page_id")
			table.Json("tree")
			table.String("message").Nullable()
			table.Timestamps()
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("landing_components") {
		if err := facades.Schema().Create("landing_components", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.Json("tree")
			table.Timestamps()
		}); err != nil {
			return err
		}
	}

	return nil
}

func (r *M20260713000001CreateLandingPagesTables) Down() error {
	facades.Schema().DropIfExists("landing_components")
	facades.Schema().DropIfExists("landing_page_revisions")
	facades.Schema().DropIfExists("landing_pages")
	return nil
}
