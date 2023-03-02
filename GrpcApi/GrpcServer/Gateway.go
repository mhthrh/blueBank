package GrpcServer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mhthrh/BlueBank/Db"
	"github.com/mhthrh/BlueBank/Entity"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoGateway"
	"github.com/mhthrh/BlueBank/Redis"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g GatewayServer) GatewayLogin(ctx context.Context, in *ProtoGateway.GatewayLoginRequest) (*ProtoGateway.GatewayLoginResponse, error) {
	p := <-*pool
	defer func() {
		_ = p.Sql.Close()
		_ = p.Redis.Close()
	}()

	client := Redis.Client{Client: p.Redis}
	byt, _ := json.Marshal(in)
	err := client.Set(p.Id.String(), string(byt))
	if err != nil {
		fmt.Printf("canot insert to redis, %v\n", err)
	}
	db := Db.NewDb(p.Sql)

	gateway, err := db.GatewayLogin(ctx, Entity.GatewayLogin{
		UserName: in.UserName,
		Password: in.Password,
	})

	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "", err)
	}
	r := &ProtoGateway.GatewayLoginResponse{
		UserName:    gateway.UserName,
		Password:    "********",
		Ips:         gateway.Ips,
		GatewayName: gateway.GatewayName,
		Status:      gateway.Status,
	}
	byt, _ = json.Marshal(r)
	_ = client.Set(p.Id.String(), string(byt))
	return r, status.Errorf(codes.OK, "")

}
