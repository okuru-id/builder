package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

type M20260702000009CreateSettingsTable struct{}

// Signature The unique signature for the migration.
func (r *M20260702000009CreateSettingsTable) Signature() string {
	return "20260702000009_create_settings_table"
}

// Up Run the migrations.
func (r *M20260702000009CreateSettingsTable) Up() error {
	if !facades.Schema().HasTable("settings") {
		if err := facades.Schema().Create("settings", func(table schema.Blueprint) {
			table.ID()
			table.String("key")
			table.Unique("key")
			table.Text("value").Nullable()
			table.Timestamps()
		}); err != nil {
			return err
		}
	}
	return nil
}

// Down Reverse the migrations.
func (r *M20260702000009CreateSettingsTable) Down() error {
	if err := facades.Schema().DropIfExists("settings"); err != nil {
		return err
	}
	return nil
}
