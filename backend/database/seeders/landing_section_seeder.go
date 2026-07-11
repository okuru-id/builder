package seeders

import (
	"okuru/app/facades"
	"okuru/app/models"
)

type LandingSectionSeeder struct{}

func (s *LandingSectionSeeder) Signature() string {
	return "LandingSectionSeeder"
}

func (s *LandingSectionSeeder) Run() error {
	sections := []models.LandingSection{
		{
			Type: "hero",
			Content: models.LandingContent{
				"greeting_en":         "Hi, I'm Kurob",
				"greeting_id":         "Halo, Saya Kurob",
				"description_en":      "Kurob is a fullstack programmer from Cilegon with 10+ years of experience building robust Laravel, React, DevOps, and automation solutions for modern businesses.",
				"description_id":      "Kurob adalah fullstack programmer dari Cilegon dengan pengalaman 10+ tahun membangun solusi Laravel, React, DevOps, dan automation untuk kebutuhan bisnis modern.",
				"profile_image":       "/images/profile.webp",
				"profile_image_mobile": "/images/profile-half.webp",
				"cta_text":            "Contact Me",
				"cta_link":            "#contact",
			},
			SortOrder: 1,
			IsActive:  true,
		},
		{
			Type: "clients",
			Content: models.LandingContent{
				"items": []any{
					map[string]any{"name": "Danareksa", "logo": "/images/danareksa.png"},
					map[string]any{"name": "Krakatau Medika", "logo": "/images/km.png"},
					map[string]any{"name": "CIMB NIAGA", "logo": "/images/cimb-niaga.png"},
					map[string]any{"name": "IASH", "logo": "/images/iash.png"},
					map[string]any{"name": "Krakatau Steel", "logo": "/images/ks.png"},
					map[string]any{"name": "Krakatau Tirta Industri", "logo": "/images/kti.png"},
					map[string]any{"name": "Pemkot Cilegon", "logo": "/images/pemkot-cilegon.png"},
					map[string]any{"name": "Yayasan Pendidikan Telkom", "logo": "/images/ypt.png"},
				},
			},
			SortOrder: 2,
			IsActive:  true,
		},
		{
			Type: "services",
			Content: models.LandingContent{
				"items": []any{
					map[string]any{"title": "Fullstack App", "description_en": "Building end-to-end web and mobile applications with modern frameworks and scalable architecture.", "description_id": "Membangun aplikasi web dan mobile end-to-end dengan framework modern dan arsitektur scalable.", "icon": "IconCode"},
					map[string]any{"title": "AI Integration Project", "description_en": "Integrating AI capabilities into existing systems to enhance automation and decision-making.", "description_id": "Mengintegrasikan kemampuan AI ke dalam sistem yang ada untuk meningkatkan otomatisasi dan pengambilan keputusan.", "icon": "IconBrain"},
					map[string]any{"title": "Automation Workflow System", "description_en": "Designing workflow automation to streamline business processes and improve efficiency.", "description_id": "Merancang otomatisasi alur kerja untuk menyederhanakan proses bisnis dan meningkatkan efisiensi.", "icon": "IconRoute"},
					map[string]any{"title": "DevOps", "description_en": "Setting up CI/CD pipelines, containerization, and infrastructure for reliable deployments.", "description_id": "Menyiapkan pipeline CI/CD, containerization, dan infrastruktur untuk deployment yang andal.", "icon": "IconBrandDocker"},
				},
			},
			SortOrder: 3,
			IsActive:  true,
		},
		{
			Type: "projects",
			Content: models.LandingContent{
				"items": []any{
					map[string]any{"title_en": "Resident Fee Management System", "title_id": "Sistem Pengelolaan Iuran Warga", "description_en": "A modern solution for managing resident fee payments designed to facilitate transparent recording, automatic calculations, and guaranteed security with modern encryption.", "description_id": "Solusi modern untuk manajemen pembayaran iuran warga yang dirancang untuk memudahkan pencatatan, perhitungan, dan pelacakan iuran secara transparan dan efisien.", "github_url": "https://github.com/kurob1993/kelola", "technologies": []any{"Laravel", "PHP", "PostgreSQL", "Tailwind CSS"}},
					map[string]any{"title_en": "SQL Server Docker Backup Script", "title_id": "Script Backup SQL Server Docker", "description_en": "A comprehensive bash script for automating SQL Server database backups in Docker environments.", "description_id": "Script bash komprehensif untuk mengotomatisasi backup database SQL Server dalam lingkungan Docker.", "github_url": "https://github.com/kurob1993/sql-server-dockerize-backup", "technologies": []any{"Docker", "Bash", "SQL Server", "DevOps"}},
					map[string]any{"title_en": "SonarQube Report Generator", "title_id": "Generator Laporan SonarQube", "description_en": "A Go-based web application for SonarQube report generation with Jenkins integration.", "description_id": "Aplikasi web berbasis Go untuk pembuatan laporan SonarQube dengan integrasi Jenkins.", "github_url": "https://github.com/okuru-id/sonar-report", "technologies": []any{"Go", "Docker", "SonarQube", "Jenkins"}},
					map[string]any{"title_en": "Tactile Match Game", "title_id": "Permainan Tactile Match", "description_en": "A fun and interactive tactile matching game built with JavaScript.", "description_id": "Permainan mencocokkan tactile yang seru dan interaktif dibangun dengan JavaScript.", "github_url": "https://github.com/okuru-id/tactile-match-game", "technologies": []any{"JavaScript", "HTML", "CSS"}},
				},
			},
			SortOrder: 4,
			IsActive:  true,
		},
		{
			Type: "cta",
			Content: models.LandingContent{
				"heading":   "Tell me about your next project",
				"email":     "kurob@okuru.id",
				"whatsapp":  "628999702143",
			},
			SortOrder: 5,
			IsActive:  true,
		},
	}

	for _, s := range sections {
		var existing models.LandingSection
		_ = facades.Orm().Query().Where("type = ?", s.Type).First(&existing)
		if existing.ID == 0 {
			if err := facades.Orm().Query().Create(&s); err != nil {
				return err
			}
		}
	}
	return nil
}
