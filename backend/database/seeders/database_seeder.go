package seeders

import (
	"okuru/app/facades"
	"okuru/app/models"
)

type DatabaseSeeder struct{}

func (s *DatabaseSeeder) Signature() string {
	return "DatabaseSeeder"
}

func (s *DatabaseSeeder) Run() error {
	if err := s.seedAdmin(); err != nil {
		return err
	}
	if err := s.seedCategories(); err != nil {
		return err
	}
	if err := s.seedSettings(); err != nil {
		return err
	}
	if err := seedLandingComponents(); err != nil {
		return err
	}
	return nil
}

func (s *DatabaseSeeder) seedAdmin() error {
	hashed, err := facades.Hash().Make("initial.1")
	if err != nil {
		return err
	}
	var user models.User
	return facades.Orm().Query().UpdateOrCreate(
		&user,
		models.User{Email: "admin@okuru.id"},
		models.User{Password: hashed},
	)
}

func (s *DatabaseSeeder) seedCategories() error {
	categories := []models.Category{
		{Slug: "web-development", NameEn: "Web Development", NameId: "Pengembangan Web"},
		{Slug: "devops", NameEn: "DevOps", NameId: "DevOps"},
		{Slug: "javascript", NameEn: "JavaScript", NameId: "JavaScript"},
		{Slug: "laravel", NameEn: "Laravel", NameId: "Laravel"},
		{Slug: "react", NameEn: "React", NameId: "React"},
		{Slug: "docker", NameEn: "Docker", NameId: "Docker"},
	}
	for _, c := range categories {
		var cat models.Category
		if err := facades.Orm().Query().UpdateOrCreate(
			&cat,
			models.Category{Slug: c.Slug},
			c,
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *DatabaseSeeder) seedSettings() error {
	settings := []models.Setting{
		{Key: "hero_title_en", Value: "Hi, I'm Kurob"},
		{Key: "hero_title_id", Value: "Halo, Saya Kurob"},
		{Key: "hero_desc_en", Value: "I'm from Indonesia with 8+ years of experience as a full-stack web developer."},
		{Key: "hero_desc_id", Value: "Saya dari Indonesia dengan pengalaman 8+ tahun sebagai full-stack web developer."},
		{Key: "start_year", Value: "2016"},
	}
	for _, st := range settings {
		var setting models.Setting
		if err := facades.Orm().Query().UpdateOrCreate(
			&setting,
			models.Setting{Key: st.Key},
			st,
		); err != nil {
			return err
		}
	}
	return nil
}

