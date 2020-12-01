package models_test

import (
	"temperature/models"

	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewJWT(t *testing.T) {
	assert := assert.New(t)
	secret, _ := hex.DecodeString("dc2dc8a96e7050c54ee5267363f9cd803912ea81")

	exp := time.Now().Add(24000 * time.Hour)
	jwt, err := models.NewJWT(secret, "temperature", "471", &exp, nil)
	assert.Nil(err)
	assert.NotEqual("", jwt.String())
	parsed, err := models.ParseJWT(secret, jwt.String())
	assert.Nil(err)
	assert.Equal("temperature", parsed.Claims["sub"].(string))
	assert.Equal("471", parsed.Claims["uid"].(string))
	assert.Equal(float64(exp.Unix()), parsed.Claims["exp"].(float64))

	exp = time.Now().Add(time.Second)
	jwt, err = models.NewJWT(secret, "temperature", "9527", &exp, nil)
	assert.Nil(err)
	parsed, err = models.ParseJWT(secret, jwt.String())
	assert.Nil(err)
	time.Sleep(2 * time.Second)
	parsed, err = models.ParseJWT(secret, jwt.String())
	assert.NotNil(err, "token should expire")
}

func TestNewJWT_WithPayloads(t *testing.T) {
	assert := assert.New(t)
	secret := []byte("secret")

	exp := time.Now().Add(time.Hour)
	jwt, err := models.NewJWT(secret, "temperature", "9527", &exp, map[string]interface{}{
		"integer": 1234,
		"string":  "string",
	})
	assert.Nil(err)
	assert.NotEqual("", jwt.String())

	parsed, err := models.ParseJWT(secret, jwt.String())
	assert.Nil(err)
	assert.Equal(1234.0, parsed.Claims["integer"].(float64))
	assert.Equal("string", parsed.Claims["string"].(string))
}

func TestParseJWTInvalid(t *testing.T) {
	assert := assert.New(t)
	secret := []byte("secret")

	parsed, err := models.ParseJWT(secret, "invalid token")
	assert.Nil(parsed)
	assert.Error(err)
}
