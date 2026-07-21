package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

// M20260721000002AddUserIdToContentTables scopes every admin-managed content
// row to its owner. Default 1 = first admin (site owner), so pre-existing rows
// stay readable by the original admin after migration.
type M20260721000002AddUserIdToContentTables struct{}

func (r *M20260721000002AddUserIdToContentTables) Signature() string {
	return "20260721000002_add_user_id_to_content_tables"
}

func (r *M20260721000002AddUserIdToContentTables) Up() error {
	tables := []string{
		"posts", "projects", "open_source_projects", "products",
		"categories", "messages", "settings",
		"landing_pages", "landing_components",
	}
	for _, t := range tables {
		if !facades.Schema().HasTable(t) {
			continue
		}
		if facades.Schema().HasColumn(t, "user_id") {
			continue
		}
		if err := facades.Schema().Table(t, func(table schema.Blueprint) {
			table.UnsignedBigInteger("user_id").Default(1)
			table.Index("user_id")
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260721000002AddUserIdToContentTables) Down() error {
	return nil
}
