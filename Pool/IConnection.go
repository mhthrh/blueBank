package Pool

import (
	"database/sql"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/mhthrh/BlueBank/KafkaBroker"
	"google.golang.org/grpc"
)

type Connection struct {
	Id             uuid.UUID
	Sql            *sql.DB
	Redis          *redis.Client
	GrpcConnection *grpc.ClientConn
	KafkaReader    KafkaBroker.Reader
	KafkaWriter    KafkaBroker.Writer
}
type IConnection interface {
	Fetch() (*Connection, error)
	Release(*chan Connection) []error
}
