package Entity

import "time"

type Account struct {
	Name          string
	CustomerId    string
	AccountNumber int64
	Balance       int64
	LockAmount    int64
	CreateAt      time.Time
}
