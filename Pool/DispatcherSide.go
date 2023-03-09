package Pool

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mhthrh/BlueBank/KafkaBroker"
	"google.golang.org/grpc"
	"time"
)

type DispatchSide struct {
	GrpcAddress []interface{}
	Writer      KafkaAddress
}

func NewDispatchSide(kafkaWriter KafkaAddress, grpcAddress []interface{}) IConnection {
	return &DispatchSide{
		GrpcAddress: grpcAddress,
		Writer:      kafkaWriter,
	}
}

func (d *DispatchSide) Fetch() (*Connection, error) {
	address := fmt.Sprintf("%s:%.f", d.GrpcAddress[counter%len(d.GrpcAddress)].(map[string]interface{})["ip"], d.GrpcAddress[counter%len(d.GrpcAddress)].(map[string]interface{})["port"])
	counter++
	gConn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(500*time.Millisecond))
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("did not connect: %v", err))
	}

	writer := KafkaBroker.NewWriter(fmt.Sprintf("%s:%d", d.Writer.Ip, d.Writer.Port))

	c := Connection{
		Id:             uuid.New(),
		GrpcConnection: gConn,
		KafkaWriter:    *writer,
	}
	return &c, nil
}

func (d *DispatchSide) Release(c *chan Connection) []error {
	var errs []error
	fmt.Println()

	for i := 0; i < len(*c); i++ {
		select {
		case c := <-*c:
			err := c.KafkaWriter.Connection.Close()
			if err != nil {
				errs = append(errs, err)
			}
			fmt.Printf("closed connection id %s from pool.\n", c.Id.String())
		case <-time.Tick(time.Second * 5):
			return errs
		}
	}
	return errs
}
