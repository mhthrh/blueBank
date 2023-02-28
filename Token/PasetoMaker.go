package Token

import (
	"fmt"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

type PasetoMaker struct {
	Paseto *paseto.V2
	key    []byte
}

func NewPasetoMaker(key string) (Token, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("input key size not valid, it should be %d", chacha20poly1305.KeySize)
	}
	paseto := PasetoMaker{
		Paseto: paseto.NewV2(),
		key:    []byte(key),
	}
	return &paseto, nil
}

func (p *PasetoMaker) Create(userName string, duration time.Duration) (string, error) {
	payload, err := NewPayload(userName, duration)
	if err != nil {
		return "", err
	}
	return p.Paseto.Encrypt(p.key, payload, nil)
}

func (p *PasetoMaker) Verify(token string) (*PayLoad, error) {
	payload := &PayLoad{}
	err := p.Paseto.Decrypt(token, p.key, payload, nil)
	if err != nil {
		return nil, err
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
