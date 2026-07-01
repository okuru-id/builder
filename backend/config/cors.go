package config

import (
	"okuru/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("cors", map[string]any{
		"paths": []string{"api/*"},

		"allowed_methods": []string{"GET", "POST", "OPTIONS"},
		"allowed_origins": []string{
			"https://okuru.id",
			"https://www.okuru.id",
			"https://shop.okuru.id",
			"http://localhost:3000",
			"http://localhost:5173",
		},
		"allowed_headers":      []string{"*"},
		"exposed_headers":      []string{},
		"max_age":              0,
		"supports_credentials": true,
	})
}
