package models

import (
	"fmt"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWT struct {
	Claims jwt.MapClaims
	token  string
}

func (j JWT) String() string {
	return j.token
}

func NewJWT(secret []byte, subject, uid string, expiredAt *time.Time, payloads map[string]interface{}) (*JWT, error) {
	claims := jwt.MapClaims{
		"sub": subject,
		"uid": uid,
	}
	if expiredAt != nil {
		claims["exp"] = expiredAt.Unix()
	}
	for k, v := range payloads {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret)
	return &JWT{claims, ss}, err
}

func ParseJWT(secret []byte, tokenString string) (*JWT, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, fmt.Errorf("token parsing failed")
	}

	if !token.Valid {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				err = fmt.Errorf("That's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				err = fmt.Errorf("Token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				err = fmt.Errorf("Token is not active yet")
			} else {
				err = fmt.Errorf("Could not handle this token, err: %v", err)
			}
			return nil, fmt.Errorf("token invalid, err=%+v", err)
		} else {
			log.Printf("token invalid, err = %+v\n", err)
			return nil, fmt.Errorf("token invalid, unknown error")
		}
	}

	result := &JWT{token: tokenString}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		result.Claims = claims
	}
	return result, nil
}
