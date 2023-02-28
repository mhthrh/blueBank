package Token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const minSecuritySize = 32

type JwtMaker struct {
	SecretKey string
}

func NewJwtMaker(secretKey string) (Token, error) {
	if len(secretKey) < minSecuritySize {
		return nil, fmt.Errorf("size of security key is less than %d", minSecuritySize)
	}
	return &JwtMaker{secretKey}, nil
}

func (j *JwtMaker) Create(userName string, duration time.Duration) (string, error) {
	payLoad, err := NewPayload(userName, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payLoad)
	return jwtToken.SignedString([]byte(j.SecretKey))
}

func (j *JwtMaker) Verify(token string) (*PayLoad, error) {

	t, err := jwt.ParseWithClaims(token, &PayLoad{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("validation error, %v", err)
	}
	payload, ok := t.Claims.(*PayLoad)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	return payload, nil
}
