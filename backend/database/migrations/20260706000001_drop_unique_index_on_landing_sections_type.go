package migrations

import (
	"github.com/goravel/framework/facades"
)

type M20260706000001DropUniqueIndexOnLandingSectionsType struct{}

func (m *M20260706000001DropUniqueIndexOnLandingSectionsType) Signature() string {
	return "20260706000001_drop_unique_index_on_landing_sections_type"
}

func (m *M20260706000001DropUniqueIndexOnLandingSectionsType) Up() error {
	facades.DB().Statement("ALTER TABLE landing_sections DROP CONSTRAINT IF EXISTS landing_sections_type_unique")
	return nil
}

func (m *M20260706000001DropUniqueIndexOnLandingSectionsType) Down() error {
	return nil
}
