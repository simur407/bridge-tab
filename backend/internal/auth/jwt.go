package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: dev only; change and hide key behind envs
var key = []byte("ZT32sOpCt6HzF7dBPVlPHIARsgiwIGbCJmyADp9iRoWhKNhtkQj0bGkrnkMZzDpX")
// 10 days
const EXPIRES_AT = time.Hour * 24 * 10 

func Generate(userId string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   userId,
		"exp": time.Now().Add(EXPIRES_AT).Unix(),
	}).SignedString(key)
}

func Validate(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return key, nil
	})
}
