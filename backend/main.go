package main

import (
	"okuru/bootstrap"
)

func main() {
	app := bootstrap.Boot()

	app.Start()
}
