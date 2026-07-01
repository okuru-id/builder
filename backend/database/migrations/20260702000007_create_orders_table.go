package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

type M20260702000007CreateOrdersTable struct{}

// Signature The unique signature for the migration.
func (r *M20260702000007CreateOrdersTable) Signature() string {
	return "20260702000007_create_orders_table"
}

// Up Run the migrations.
func (r *M20260702000007CreateOrdersTable) Up() error {
	if !facades.Schema().HasTable("orders") {
		if err := facades.Schema().Create("orders", func(table schema.Blueprint) {
			table.ID()
			table.String("reference")
			table.Unique("reference")
			table.UnsignedBigInteger("product_id")
			table.String("buyer_email")
			table.String("buyer_name")
			table.Integer("amount").Default(0)
			table.String("status").Default("pending")
			table.String("payment_ref").Nullable()
			table.String("download_token").Nullable()
			table.Integer("download_count").Default(0)
			table.Integer("max_downloads").Default(3)
			table.Timestamp("expires_at").Nullable()
			table.Timestamp("paid_at").Nullable()
			table.Timestamps()
		}); err != nil {
			return err
		}
	}
	return nil
}

// Down Reverse the migrations.
func (r *M20260702000007CreateOrdersTable) Down() error {
	if err := facades.Schema().DropIfExists("orders"); err != nil {
		return err
	}
	return nil
}
