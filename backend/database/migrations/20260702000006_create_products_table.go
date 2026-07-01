package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

type M20260702000006CreateProductsTable struct{}

// Signature The unique signature for the migration.
func (r *M20260702000006CreateProductsTable) Signature() string {
	return "20260702000006_create_products_table"
}

// Up Run the migrations.
func (r *M20260702000006CreateProductsTable) Up() error {
	if !facades.Schema().HasTable("products") {
		if err := facades.Schema().Create("products", func(table schema.Blueprint) {
			table.ID()
			table.String("slug")
			table.Unique("slug")
			table.String("title")
			table.Text("description").Nullable()
			table.Integer("price").Default(0)
			table.String("type")
			table.String("file_path").Nullable()
			table.String("thumbnail").Nullable()
			table.String("status").Default("active")
			table.Timestamps()
		}); err != nil {
			return err
		}
	}
	return nil
}

// Down Reverse the migrations.
func (r *M20260702000006CreateProductsTable) Down() error {
	if err := facades.Schema().DropIfExists("products"); err != nil {
		return err
	}
	return nil
}
