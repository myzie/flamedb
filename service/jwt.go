package service

import (
	"crypto/rsa"
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/myzie/flamedb/models"
)

func parseJWT(key *rsa.PublicKey, tokenStr string) (*models.Principal, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, errors.New("Unknown JWT claims type")
	}
	if !token.Valid {
		return nil, errors.New("JWT claims are invalid")
	}

	return &models.Principal{
		UserID:      claims.Subject,
		Permissions: "rw",
	}, nil
}
