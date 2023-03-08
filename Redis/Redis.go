package Redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Client struct {
	Client *redis.Client
}

func (c *Client) Set(Key, values string) (e error) {
	defer func() {
		err := recover()
		if err != nil {
			e = err.(error)
		}
	}()
	cnt, _ := c.KeyExist(Key)
	if cnt > 0 {
		val, _ := c.Get(Key)
		result := c.Client.Set(Key, fmt.Sprintf("%s\n %s", val, values), 0)
		_, err := result.Result()
		return err
	}
	result := c.Client.Set(Key, values, 0)
	_, err := result.Result()

	return err
}

func (c *Client) KeyExist(Key string) (int, error) {
	cc, err := c.Client.Exists(Key).Result()

	if err != nil {
		return -1, err
	}
	return int(cc), nil
}

func (c *Client) Get(key string) (string, error) {
	Val, err := c.Client.Get(key).Result()
	if err != nil {
		return "", err
	}

	return Val, nil
}
