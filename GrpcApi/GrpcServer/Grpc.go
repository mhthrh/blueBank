package GrpcServer

import (
	"github.com/mhthrh/BlueBank/Pool"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoGateway"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoUser"
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

func New(p *chan Pool.Connection) {
	pool = p
}
