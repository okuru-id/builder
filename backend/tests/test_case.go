package tests

import (
	"github.com/goravel/framework/testing"

	"okuru/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
