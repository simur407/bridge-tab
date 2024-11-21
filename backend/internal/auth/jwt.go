package auth

import "github.com/golang-jwt/jwt/v5"

func Generate(id *UserId, name string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   id,
		"name": name,
	}).SignedString([]byte("secret"))
}
