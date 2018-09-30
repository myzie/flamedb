package service

import (
	"crypto/rsa"
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/myzie/flamedb/models"
)

// CustomClaims contained within a JWT
type CustomClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.StandardClaims
}

func parseJWT(key *rsa.PublicKey, tokenStr string) (*models.Principal, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("Unknown JWT claims type")
	}
	if !token.Valid {
		return nil, errors.New("JWT claims are invalid")
	}

	return &models.Principal{
		Email:  claims.Email,
		Name:   claims.Name,
		UserID: claims.Subject,
		Token:  tokenStr,
	}, nil
}
