package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

// M20260717000001AddPathDomainToLandingPages adds multi-page publishing fields:
//   - path:  public URL segment, e.g. "promo" → served at /promo (nullable; home uses is_home)
//   - domain: custom domain, e.g. "client.com" → served when Host matches (nullable)
//   - is_home: marks the page served at "/" (at most one)
//
// Uniqueness is enforced at the app layer (see landing_page_controller) because
// nullable unique indexes behave differently across SQLite/Postgres and we only
// want uniqueness among non-empty values.
type M20260717000001AddPathDomainToLandingPages struct{}

func (r *M20260717000001AddPathDomainToLandingPages) Signature() string {
	return "20260717000001_add_path_domain_to_landing_pages"
}

func (r *M20260717000001AddPathDomainToLandingPages) Up() error {
	if !facades.Schema().HasColumn("landing_pages", "path") {
		if err := facades.Schema().Table("landing_pages", func(table schema.Blueprint) {
			table.String("path").Nullable()
		}); err != nil {
			return err
		}
	}
	if !facades.Schema().HasColumn("landing_pages", "domain") {
		if err := facades.Schema().Table("landing_pages", func(table schema.Blueprint) {
			table.String("domain").Nullable()
		}); err != nil {
			return err
		}
	}
	if !facades.Schema().HasColumn("landing_pages", "is_home") {
		if err := facades.Schema().Table("landing_pages", func(table schema.Blueprint) {
			table.Boolean("is_home").Default(false)
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260717000001AddPathDomainToLandingPages) Down() error {
	if facades.Schema().HasColumn("landing_pages", "is_home") {
		facades.Schema().Table("landing_pages", func(table schema.Blueprint) {
			table.DropColumn("is_home")
		})
	}
	if facades.Schema().HasColumn("landing_pages", "domain") {
		facades.Schema().Table("landing_pages", func(table schema.Blueprint) {
			table.DropColumn("domain")
		})
	}
	if facades.Schema().HasColumn("landing_pages", "path") {
		facades.Schema().Table("landing_pages", func(table schema.Blueprint) {
			table.DropColumn("path")
		})
	}
	return nil
}
