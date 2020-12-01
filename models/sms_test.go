package models_test

import (
	"github.com/stretchr/testify/assert"
	"temperature/models"
	"testing"
)

func TestSMS_Send(t *testing.T) {
	mobile := "**"

	sms := new(models.SMS)
	sms.Mobile = mobile
	err := sms.Send()
	assert := assert.New(t)
	assert.NoError(err)
}

func TestSMS_Verify(t *testing.T) {
	mobile := "***"

	sms := new(models.SMS)
	sms.Mobile = mobile
	code := "904532"
	ok  := sms.Verify(code)
	assert := assert.New(t)
	assert.True(ok)
}