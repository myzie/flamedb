package service

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	certBegin = "-----BEGIN CERTIFICATE-----"
	certEnd   = "-----END CERTIFICATE-----"
)

// JSONWebKey is a single JSON Web Key definition
type JSONWebKey struct {
	Alg string   `json:"alg"`
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// GetCertificate returns the certificate string for this JSONWebKey
func (key JSONWebKey) GetCertificate() string {
	return fmt.Sprintf("%s\n%s\n%s", certBegin, key.X5c[0], certEnd)
}

// JSONWebKeySet contains a list of JSON web keys
type JSONWebKeySet struct {
	Keys []JSONWebKey `json:"keys"`
}

// Find looks for a key with the given key ID
// token.Header["kid"]
func (set JSONWebKeySet) Find(keyID string) (JSONWebKey, bool) {
	for i := range set.Keys {
		if keyID == set.Keys[i].Kid {
			return set.Keys[i], true
		}
	}
	return JSONWebKey{}, false
}

// GetRSAPublicKeys returns a slice of RSA public keys generated from this JWKS
func (set JSONWebKeySet) GetRSAPublicKeys() ([]*rsa.PublicKey, error) {
	var result []*rsa.PublicKey
	for _, jwk := range set.Keys {
		if jwk.Kty != "RSA" || jwk.Alg != "RS256" {
			continue
		}
		cert := jwk.GetCertificate()
		rsaKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			return nil, err
		}
		result = append(result, rsaKey)
	}
	return result, nil
}

// GetKeySet retrieves public keys from a remote JWKS endpoint
// "https://" + os.Getenv("AUTH0_DOMAIN") + "/.well-known/jwks.json"
// https://auth0.com/docs/quickstart/backend/golang/01-authorization
func GetKeySet(url string) (JSONWebKeySet, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return JSONWebKeySet{}, err
	}

	ctx, cancel := context.WithTimeout(req.Context(), time.Second)
	defer cancel()

	req = req.WithContext(ctx)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return JSONWebKeySet{}, err
	}
	defer resp.Body.Close()

	var jwks JSONWebKeySet
	if json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return JSONWebKeySet{}, err
	}
	return jwks, nil
}
