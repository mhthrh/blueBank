package Entity

import "time"

type Account struct {
	CustomerUserName string
	AccountNumber    string
	Balance          int64
	LockAmount       int64
	CreateAt         time.Time
}
