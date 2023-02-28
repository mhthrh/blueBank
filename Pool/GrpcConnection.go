package Pool

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"time"
)

type RedisAddress struct {
	Host string
	Port string
}
type RedisConnection struct {
	Ip       string
	Port     int
	Password string
	Database int
}
type GrpcConnection struct {
	DbConnection    string
	RedisConnection RedisConnection
}

func NewGrpcConnection(db string, redis RedisConnection) IConnection {
	return &GrpcConnection{
		DbConnection: db,
		RedisConnection: RedisConnection{
			Ip:       redis.Ip,
			Port:     redis.Port,
			Password: redis.Password,
			Database: redis.Database,
		},
	}
}
func (g *GrpcConnection) Fetch() (*Connection, error) {
	cnn, err := sql.Open("postgres", g.DbConnection)
	if err != nil {
		return nil, fmt.Errorf("cannot connet to db,%w", err)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", g.RedisConnection.Ip, g.RedisConnection.Port),
		Password: g.RedisConnection.Password,
		DB:       g.RedisConnection.Database,
	})
	_, err = client.Ping().Result()
	if err != nil {
		errDb := cnn.Close()
		if errDb != nil {
			return nil, fmt.Errorf("cannot connet to db,%w,cannot connet to redis, %v", errDb, err)
		}
		return nil, fmt.Errorf("cannot connet to redis,%w", err)
	}
	c := Connection{
		Id:    uuid.New(),
		Sql:   cnn,
		Redis: client,
	}
	return &c, nil
}

func (g *GrpcConnection) Release(p *chan Connection) []error {
	var errs []error
	fmt.Println()

	for i := 0; i < len(*p); i++ {
		select {
		case c := <-*p:
			err := c.Sql.Close()
			if err != nil {
				errs = append(errs, err)
			}
			err = c.Redis.Close()
			if err != nil {
				errs = append(errs, err)
			}
			fmt.Printf("closed connection id %s from pool.\n", c.Id.String())
		case <-time.Tick(time.Second * 5):
			return errs
		}

	}

	return errs
}
