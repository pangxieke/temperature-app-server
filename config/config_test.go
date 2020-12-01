package config_test

import (
	"temperature/config"

	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("MYSQL_HOST", "new host")
	err := config.Init("../test")
	assert.NoError(err)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal("new host", config.MySQL.Host, "should be override by environment")
	//assert.Equal("temperature_fixtures", config.MySQL.DB)
}
