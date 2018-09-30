package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	certBegin = "-----BEGIN CERTIFICATE-----"
	certEnd   = "-----END CERTIFICATE-----"
)

// JSONWebKey is a single JSON Web Key definition
type JSONWebKey struct {
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
func (keys JSONWebKeySet) Find(keyID string) (JSONWebKey, bool) {
	for k := range keys.Keys {
		if keyID == keys.Keys[k].Kid {
			return keys.Keys[k], true
		}
	}
	return JSONWebKey{}, false
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
