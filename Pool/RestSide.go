package Pool

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/mhthrh/BlueBank/KafkaBroker"
	"time"
)

type KafkaAddress struct {
	Ip   string
	Port int
}

type RestSide struct {
	RedisConnection RedisConnection
	Writer          KafkaAddress
}

func NewRestSide(redis RedisConnection, kafkaWriter KafkaAddress) IConnection {
	return &RestSide{
		RedisConnection: RedisConnection{
			Ip:       redis.Ip,
			Port:     redis.Port,
			Password: redis.Password,
			Database: redis.Database,
		},
		Writer: kafkaWriter,
	}
}

func (r *RestSide) Fetch() (*Connection, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.RedisConnection.Ip, r.RedisConnection.Port),
		Password: r.RedisConnection.Password,
		DB:       r.RedisConnection.Database,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("cannot connet to redis,%w", err)
	}

	writer := KafkaBroker.NewWriter(fmt.Sprintf("%s:%d", r.Writer.Ip, r.Writer.Port))

	c := Connection{
		Id:          uuid.New(),
		Redis:       client,
		KafkaWriter: *writer,
	}
	return &c, nil
}

func (r *RestSide) Release(c *chan Connection) []error {
	var errs []error
	fmt.Println()

	for i := 0; i < len(*c); i++ {
		select {
		case c := <-*c:
			err := c.Redis.Close()
			if err != nil {
				errs = append(errs, err)
			}
			err = c.KafkaWriter.Connection.Close()
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
