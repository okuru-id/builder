package bootstrap

import (
	contractsseeder "github.com/goravel/framework/contracts/database/seeder"

	"okuru/database/seeders"
)

func Seeders() []contractsseeder.Seeder {
	return []contractsseeder.Seeder{
		&seeders.DatabaseSeeder{},
	}
}
