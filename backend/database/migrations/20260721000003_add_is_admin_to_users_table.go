package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

type M20260721000003AddIsAdminToUsersTable struct{}

func (r *M20260721000003AddIsAdminToUsersTable) Signature() string {
	return "20260721000003_add_is_admin_to_users_table"
}

func (r *M20260721000003AddIsAdminToUsersTable) Up() error {
	if facades.Schema().HasTable("users") && !facades.Schema().HasColumn("users", "is_admin") {
		if err := facades.Schema().Table("users", func(table schema.Blueprint) {
			table.Boolean("is_admin").Default(false)
		}); err != nil {
			return err
		}
	}
	// Backfill: lowest-ID user (previous heuristic) becomes admin.
	_, err := facades.Orm().Query().Exec(
		`UPDATE users SET is_admin = (id = (SELECT min_id FROM (SELECT MIN(id) AS min_id FROM users) m))`,
	)
	return err
}

func (r *M20260721000003AddIsAdminToUsersTable) Down() error {
	return nil
}
