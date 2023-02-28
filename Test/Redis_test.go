package Test

import (
	"github.com/mhthrh/BlueBank/Redis"
	"github.com/mhthrh/BlueBank/Utils/RandomUtil"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Redis(t *testing.T) {
	rndKey := RandomUtil.RandomString(10)
	rndValue := RandomUtil.RandomString(30)

	client := Redis.Client{Client: RedisClient}

	err := client.Set(rndKey, rndValue)
	require.NoError(t, err)

	value, err := client.Get(rndKey)
	require.NoError(t, err)
	require.Equal(t, rndValue, value)

	cnt, err := client.KeyExist(rndKey)
	require.NoError(t, err)
	require.Equal(t, cnt, 1)

}
