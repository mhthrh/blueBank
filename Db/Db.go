package Db

import (
	"database/sql"
	"github.com/mhthrh/BlueBank/Utils/CryptoUtil"
)

var (
	c *CryptoUtil.Crypto
)

func init() {
	c = CryptoUtil.NewKey()
}

type dataBase struct {
	db *sql.DB
}

func NewDb(db *sql.DB) *dataBase {
	return &dataBase{
		db: db,
	}
}
