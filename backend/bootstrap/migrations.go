package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"okuru/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
		&migrations.M20260702000001CreateUsersTable{},
		&migrations.M20260702000002CreateCategoriesTable{},
		&migrations.M20260702000003CreatePostsTable{},
		&migrations.M20260702000004CreateProjectsTable{},
		&migrations.M20260702000005CreateOpenSourceProjectsTable{},
		&migrations.M20260702000006CreateProductsTable{},
		&migrations.M20260702000007CreateOrdersTable{},
		&migrations.M20260702000008CreateMessagesTable{},
		&migrations.M20260702000009CreateSettingsTable{},
		&migrations.M20260713000001CreateLandingPagesTables{},
		&migrations.M20260713000002AddSystemFieldsToLandingComponentsTable{},
		&migrations.M20260717000001AddPathDomainToLandingPages{},
	}
}
