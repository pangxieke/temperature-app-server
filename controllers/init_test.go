package controllers_test

import (
	"temperature/routers"
	"temperature/test"
	"github.com/julienschmidt/httprouter"
	"os"
	"testing"
)

var router *httprouter.Router

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func setup() {
	router = routers.Router()
	test.Init()
}
