package models_test

import (
	"os"
	"testing"
	"temperature/test"
)

func TestMain(m *testing.M) {
	test.Init()
	os.Exit(m.Run())
}
