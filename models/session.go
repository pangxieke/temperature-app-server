package models

import (
	"fmt"
	"time"
)

type Session struct {
	Token JWT
	Exp   time.Time
	Uid   int
	FaceId uint   `json:"face_id"`
	UserID uint   `json:"user_id"`
	Mobile string `json:"mobile"`
}

var Subject = "temperature"

func NewSession(mobile string) (s *Session, err error) {
	exp := time.Now().Add(24 * time.Hour)
	jwt, err := NewJWT(secret, Subject, mobile, &exp, nil)
	if err != nil {
		return
	}
	s = &Session{Mobile: mobile, Token: *jwt, Exp: exp}
	return
}

func LoadSession(token string) (s *Session, err error) {
	jwt, err := ParseJWT(secret, token)
	if err != nil {
		return
	}
	fmt.Println(jwt)
	mobile := jwt.Claims["uid"].(string)
	if err != nil {
		return
	}
	return &Session{Mobile: mobile, Token: *jwt}, nil
}