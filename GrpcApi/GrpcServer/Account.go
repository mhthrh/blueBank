package GrpcServer

import (
	"context"
	"fmt"
	"github.com/mhthrh/BlueBank/Db"
	"github.com/mhthrh/BlueBank/Entity"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoAccount"
	"github.com/mhthrh/BlueBank/Redis"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (v *AccountServer) Create(ctx context.Context, in *ProtoAccount.CreateRequest) (*ProtoAccount.CreateResponse, error) {
	p := <-*pool
	defer func() {
		_ = p.Sql.Close()
		_ = p.Redis.Close()
	}()

	client := Redis.Client{Client: p.Redis}
	err := client.Set(p.Id.String(), "test")
	if err != nil {
		fmt.Printf("canot insert to redis, %v\n", err)
	}
	db := Db.NewDb(p.Sql)
	err = db.CreateAccount(ctx, Entity.Account{
		CustomerUserName: in.UserName,
		AccountNumber:    "",
		Balance:          0,
		LockAmount:       0,
		CreateAt:         time.Now(),
	})
	if err != nil {
		return &ProtoAccount.CreateResponse{Error: err.Error()}, status.Errorf(codes.FailedPrecondition, err.Error())
	}
	return &ProtoAccount.CreateResponse{Error: ""}, status.Errorf(codes.OK, "")
}

func (v *AccountServer) Balance(ctx context.Context, in *ProtoAccount.BalanceRequest) (*ProtoAccount.BalanceResponse, error) {
	p := <-*pool
	defer func() {
		_ = p.Sql.Close()
		_ = p.Redis.Close()
	}()

	client := Redis.Client{Client: p.Redis}
	err := client.Set(p.Id.String(), "test")
	if err != nil {
		fmt.Printf("canot insert to redis, %v\n", err)
	}
	db := Db.NewDb(p.Sql)
	amount, err := db.BalanceAccount(ctx, Entity.Account{
		CustomerUserName: in.UserName,
		AccountNumber:    in.AccountNumber,
		Balance:          0,
		LockAmount:       0,
		CreateAt:         time.Now(),
	})

	if err != nil {
		return &ProtoAccount.BalanceResponse{Balance: 0}, status.Errorf(codes.FailedPrecondition, err.Error())
	}
	return &ProtoAccount.BalanceResponse{Balance: float32(amount)}, status.Errorf(codes.OK, "")
}
