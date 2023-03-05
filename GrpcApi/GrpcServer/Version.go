package GrpcServer

import (
	"context"
	"fmt"
	"github.com/mhthrh/BlueBank/Db"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoVersion"
	"github.com/mhthrh/BlueBank/Redis"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (v *VersionServer) GetVersion(ctx context.Context, in *ProtoVersion.VersionRequest) (*ProtoVersion.VersionResponse, error) {
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
	version, err := db.GetVersion(ctx, "RestVersion")

	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "", err)
	}
	return &ProtoVersion.VersionResponse{Value: version}, status.Errorf(codes.OK, "")
}
