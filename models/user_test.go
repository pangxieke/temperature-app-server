package models_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"temperature/models"
	"temperature/test"
	"testing"
)

func TestFindUserByBaiDuId(t *testing.T) {
	test.Prepare(t)
	assert := assert.New(t)
	baiDuId := "11827"
	user, err := models.FindUserByBaiDuId(baiDuId)

	assert.Nil(err)
	fmt.Println(user)
	assert.NotEmpty(user)
	assert.NotEmpty(user.FaceId)

}
