package Pool

import (
	"database/sql"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type Connection struct {
	Id    uuid.UUID
	Sql   *sql.DB
	Redis *redis.Client
}
type IConnection interface {
	Fetch() (*Connection, error)
	Release(*chan Connection) []error
}
