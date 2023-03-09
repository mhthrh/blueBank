package Entity

import "time"

type Account struct {
	CustomerUserName string    `json:"customerUserName"`
	AccountNumber    string    `json:"accountNumber"`
	Balance          int64     `json:"balance"`
	LockAmount       int64     `json:"lockAmount"`
	CreateAt         time.Time `json:"createAt"`
}
