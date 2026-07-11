package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

type M20260712000001AddHtmlToLandingTemplates struct{}

func (r *M20260712000001AddHtmlToLandingTemplates) Signature() string {
	return "20260712000001_add_html_to_landing_templates"
}

func (r *M20260712000001AddHtmlToLandingTemplates) Up() error {
	if facades.Schema().HasTable("landing_templates") && !facades.Schema().HasColumns("landing_templates", []string{"html"}) {
		return facades.Schema().Table("landing_templates", func(table schema.Blueprint) {
			table.Text("html").Nullable()
		})
	}
	return nil
}

func (r *M20260712000001AddHtmlToLandingTemplates) Down() error {
	if facades.Schema().HasTable("landing_templates") && facades.Schema().HasColumns("landing_templates", []string{"html"}) {
		return facades.Schema().Table("landing_templates", func(table schema.Blueprint) {
			table.DropColumn("html")
		})
	}
	return nil
}
