package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

type M20260702000008CreateMessagesTable struct{}

// Signature The unique signature for the migration.
func (r *M20260702000008CreateMessagesTable) Signature() string {
	return "20260702000008_create_messages_table"
}

// Up Run the migrations.
func (r *M20260702000008CreateMessagesTable) Up() error {
	if !facades.Schema().HasTable("messages") {
		if err := facades.Schema().Create("messages", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.String("email")
			table.Text("message")
			table.String("status").Default("unread")
			table.Timestamps()
		}); err != nil {
			return err
		}
	}
	return nil
}

// Down Reverse the migrations.
func (r *M20260702000008CreateMessagesTable) Down() error {
	if err := facades.Schema().DropIfExists("messages"); err != nil {
		return err
	}
	return nil
}
