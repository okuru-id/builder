package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

type M20260702000001CreateUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20260702000001CreateUsersTable) Signature() string {
	return "20260702000001_create_users_table"
}

// Up Run the migrations.
func (r *M20260702000001CreateUsersTable) Up() error {
	if !facades.Schema().HasTable("users") {
		if err := facades.Schema().Create("users", func(table schema.Blueprint) {
			table.ID()
			table.String("email")
			table.Unique("email")
			table.String("password")
			table.String("totp_secret").Nullable()
			table.Boolean("totp_verified").Default(false)
			table.Timestamps()
		}); err != nil {
			return err
		}
	}
	return nil
}

// Down Reverse the migrations.
func (r *M20260702000001CreateUsersTable) Down() error {
	if err := facades.Schema().DropIfExists("users"); err != nil {
		return err
	}
	return nil
}
