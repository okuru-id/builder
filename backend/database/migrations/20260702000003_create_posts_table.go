package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/app/facades"
)

type M20260702000003CreatePostsTable struct{}

// Signature The unique signature for the migration.
func (r *M20260702000003CreatePostsTable) Signature() string {
	return "20260702000003_create_posts_table"
}

// Up Run the migrations.
func (r *M20260702000003CreatePostsTable) Up() error {
	if !facades.Schema().HasTable("posts") {
		if err := facades.Schema().Create("posts", func(table schema.Blueprint) {
			table.ID()
			table.String("slug")
			table.Unique("slug")
			table.String("title_en")
			table.String("title_id")
			table.Text("excerpt_en").Nullable()
			table.Text("excerpt_id").Nullable()
			table.LongText("content_en").Nullable()
			table.LongText("content_id").Nullable()
			table.String("category").Nullable()
			table.Text("tags").Nullable()
			table.String("thumbnail").Nullable()
			table.String("status").Default("draft")
			table.Timestamp("published_at").Nullable()
			table.Integer("read_time").Default(0)
			table.Timestamps()
		}); err != nil {
			return err
		}
	}
	return nil
}

// Down Reverse the migrations.
func (r *M20260702000003CreatePostsTable) Down() error {
	if err := facades.Schema().DropIfExists("posts"); err != nil {
		return err
	}
	return nil
}
