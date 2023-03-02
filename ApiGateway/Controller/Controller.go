package Controller

import (
	"github.com/mhthrh/BlueBank/Pool"
)

var (
	pool *chan Pool.Connection
)

func init() {
}
func New(t *chan Pool.Connection) {
	pool = t
}
