package service

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type authFunction func(jwt.Claims) (interface{}, error)

type jwtAuth struct {
	Key          []byte
	AuthFunction authFunction
}

func (a *jwtAuth) Authenticate(header string) (interface{}, error) {
	parts := strings.Split(header, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, fmt.Errorf("Invalid Authorization header: %s", header)
	}
	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return a.Key, nil
	})
	if err != nil {
		return nil, err
	}
	return a.AuthFunction(token.Claims)
}
