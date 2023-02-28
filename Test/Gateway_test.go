package Test

import (
	"github.com/mhthrh/BlueBank/Db"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_GatewayLogin(t *testing.T) {
	db := Db.NewDb(cnn)
	gateway, err := db.GatewayLogin(ctx, gatewayLogin)
	require.NoError(t, err)
	require.Equal(t, gatewayLogin.UserName, gateway.UserName)
	require.Equal(t, "localhost, 127.0.0.1", gateway.Ips)
	require.Equal(t, true, gateway.Status)
}
