package GrpcServer

import (
	"github.com/mhthrh/BlueBank/Pool"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoGateway"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoUser"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoVersion"
)

var (
	pool *chan Pool.Connection
)

type UserServer struct {
	ProtoUser.UnimplementedServicesServer
}
type GatewayServer struct {
	ProtoGateway.UnimplementedGatewayServicesServer
}

type VersionServer struct {
	ProtoVersion.UnimplementedVersionServicesServer
}

func New(p *chan Pool.Connection) {
	pool = p
}
