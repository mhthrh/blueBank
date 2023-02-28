package Token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type PayLoad struct {
	ID       uuid.UUID `json:"id"`
	UserName string    `json:"user_name"`
	CreateAt time.Time `json:"create_at"`
	ExpireAt time.Time `json:"expire_at"`
}

func NewPayload(UserName string, duration time.Duration) (*PayLoad, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := PayLoad{
		ID:       id,
		UserName: UserName,
		CreateAt: time.Now(),
		ExpireAt: time.Now().Add(duration),
	}
	return &payload, nil
}
func (p *PayLoad) Valid() error {
	if time.Now().After(p.ExpireAt) {
		return errors.New("invalid duration")
	}
	return nil
}
