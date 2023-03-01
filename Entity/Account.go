package Entity

import "time"

type Account struct {
	CustomerUserName string
	AccountNumber    int64
	Balance          int64
	LockAmount       int64
	CreateAt         time.Time
}
