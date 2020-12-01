package models_test

import (
	. "temperature/models"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOSS(t *testing.T) {
	u, err := Resign("http://cstorage.oss-cn-shenzhen.aliyuncs.com/******", 5000)
	assert := assert.New(t)
	assert.NoError(err)
	assert.NotEmpty(u)
}
