package Test

import (
	"github.com/mhthrh/BlueBank/Token"
	"github.com/mhthrh/BlueBank/Utils/RandomUtil"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_jwtMaker(t *testing.T) {
	token, err := Token.NewJwtMaker(RandomUtil.RandomString(32))
	require.NoError(t, err)

	userName := RandomUtil.RandomString(10)
	duration := time.Minute
	payload, err := token.Create(userName, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	p, err := token.Verify(payload)

	require.NoError(t, err)
	require.NotEmpty(t, p)

	require.NotZero(t, p.ID)

	require.Equal(t, p.UserName, userName)
}
