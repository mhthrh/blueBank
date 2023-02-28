package Token

import "time"

type Token interface {
	Create(userName string, duration time.Duration) (string, error)
	Verify(token string) (*PayLoad, error)
}
