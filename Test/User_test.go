package Test

import (
	_ "github.com/lib/pq"
	"github.com/mhthrh/BlueBank/Db"
	"github.com/stretchr/testify/require"
	"testing"
)

func createUser() error {
	db := Db.NewDb(cnn)
	return db.Create(ctx, user)
}
func Test_CreateUser(t *testing.T) {
	err := createUser()
	require.NoError(t, err)
}

func Test_ExistUser(t *testing.T) {
	db := Db.NewDb(cnn)
	count, err := db.Exist(ctx, user.UserName)
	require.NoError(t, err)
	require.Equal(t, count, 1)
}

func Test_LoginUser(t *testing.T) {
	db := Db.NewDb(cnn)
	usr, err := db.Login(ctx, login)
	require.NoError(t, err)
	require.Equal(t, usr.UserName, login.UserName)
	require.Equal(t, usr.Email, user.Email)
}
