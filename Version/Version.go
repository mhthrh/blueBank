package Version

import (
	"context"
	"fmt"
	"github.com/mhthrh/BlueBank/Pool"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoVersion"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

func CheckInstance(address ...string) error {
	var wErr error
	for _, s := range address {
		cnn, err := net.DialTimeout("tcp", s, time.Second*5)
		if err == nil {
			wErr = errors.Wrapf(fmt.Errorf("there is a listener on address: %s", s), "connection error")
			_ = cnn.Close()
		}

	}
	return wErr
}
func CheckVersion(cnn Pool.Connection, key string) error {
	defer func() {
		cnn.Redis.Close()
	}()
	gCnn := ProtoVersion.NewVersionServicesClient(cnn.GrpcConnection)
	result, stat := gCnn.GetVersion(context.Background(), &ProtoVersion.VersionRequest{
		Key: key,
	})

	st, ok := status.FromError(stat)
	if !ok {
		return fmt.Errorf("canot connect to grpc")
	}
	if st != nil {
		return fmt.Errorf("canot catch version from grpc")
	}
	cnfgVersion := viper.GetString(key)
	if cnfgVersion != result.Value {
		return fmt.Errorf("version mismatch, %s", fmt.Sprintf("config version is: %s and db version is %s", cnfgVersion, result.Value))
	}
	return nil
}
