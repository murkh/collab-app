package jwt

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type Signer interface {
	SignToken(claims jwt.MapClaims) (string, error)
}

type fileSigner struct {
	priv *rsa.PrivateKey
}

func (s *fileSigner) SignToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.priv)
}

func NewSignerFromFile(path string) (Signer, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		return nil, err
	}
	return &fileSigner{priv: key}, nil
}

func VerifyToken(pubPem []byte, tokenStr string) (*jwt.Token, error) {
	pub, err := jwt.ParseRSAPublicKeyFromPEM(pubPem)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return pub, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
