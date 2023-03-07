package GrpcServer

import (
	"context"
	"fmt"
	"github.com/mhthrh/BlueBank/Db"
	"github.com/mhthrh/BlueBank/Entity"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoUser"
	"github.com/mhthrh/BlueBank/Redis"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServer) CreateUser(ctx context.Context, in *ProtoUser.UserRequest) (*ProtoUser.Error, error) {
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
	errG := db.Create(ctx, Entity.Customer{
		FullName: in.FullName,
		UserName: in.UserName,
		PassWord: in.PassWord,
		Email:    in.Email,
	})

	if errG != nil {
		return &ProtoUser.Error{Message: errG.Error()}, status.Errorf(codes.FailedPrecondition, errG.Error())
	}
	return &ProtoUser.Error{Message: ""}, status.Errorf(codes.OK, "")

}

func (s *UserServer) LoginUser(ctx context.Context, in *ProtoUser.LoginRequest) (*ProtoUser.User, error) {
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
	user, err := db.Login(ctx, Entity.CustomerLogin{
		UserName: in.UserName,
		PassWord: in.PassWord,
	})
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "", err)
	}

	return &ProtoUser.User{
		ID:       1,
		FullName: user.FullName,
		UserName: user.UserName,
		PassWord: user.PassWord,
		Email:    user.Email,
		CreateAt: nil,
		ExpireAt: nil,
	}, status.Errorf(codes.OK, "")
}

func (s *UserServer) ExistUser(ctx context.Context, in *ProtoUser.ExistRequest) (*ProtoUser.ExistResponse, error) {
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
	i, err := db.Exist(ctx, in.UserName)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "", err)
	}
	return &ProtoUser.ExistResponse{
		Error: err.Error(),
		Count: int32(i),
	}, status.Errorf(codes.OK, "")
}
