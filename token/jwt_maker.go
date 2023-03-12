package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secret string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be atleast %d characters", minSecretKeySize)
	}

	return &JWTMaker{secret: secretKey}, nil
}

// CreteToken creates a new token for a specific username and duration
func (maker *JWTMaker) CreateToken(userName string, duration time.Duration) (string, error) {
	payload, err := NewJWTPayload(userName, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(maker.secret))
}

// VerifyToken verifies if the input token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JWTPayload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*JWTPayload)
	if !ok || !jwtToken.Valid {
		return nil, ErrInvalidToken
	}
	if err = claims.Payload.Valid(); err != nil {
		return nil, ErrExpiredToken
	}

	return &claims.Payload, nil

}
