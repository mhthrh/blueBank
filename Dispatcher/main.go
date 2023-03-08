package main

import (
	"context"
	"fmt"
	"github.com/mhthrh/BlueBank/Config"
	"github.com/mhthrh/BlueBank/KafkaBroker"
	"github.com/mhthrh/BlueBank/Pool"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var (
	osInterrupt chan os.Signal
	pool        chan Pool.Connection
	dConn       Pool.IConnection
	poolStop    chan struct{}
	poolSize    int
	methods     []string
)

func init() {
	methods = []string{"account", "transaction"}
	osInterrupt = make(chan os.Signal)
	poolStop = make(chan struct{})
	cfg := Config.New("Coded.dat", "json", "./ConfigFiles")
	if err := cfg.Initialize(); err != nil {
		log.Fatalf("unable to fill viber, %v", err)
	}
	dConn = Pool.NewDispatchSide(Pool.KafkaAddress{
		Ip:   viper.GetString("Kafka.Host"),
		Port: viper.GetInt("Kafka.Port"),
	}, viper.Get("GRPC").([]interface{}),
	)

	poolSize = viper.GetInt("ConnectionPoolCount")
	pool = make(chan Pool.Connection, poolSize)
}

func main() {
	var cancels []context.CancelFunc
	go fillPool()
c:
	if len(pool) < poolSize/2 {
		fmt.Println("pool loading in process, it might takes a while")
		<-time.Tick(time.Millisecond * 300)
		goto c
	}
	fmt.Println("connection pool fill successfully")
	for _, s := range methods {
		ctx, can := context.WithCancel(context.Background())
		dispatch(ctx, s)
		cancels = append(cancels, can)
	}
	<-osInterrupt
	fmt.Println(fmt.Sprintf("server will be down, because: %s", "got intrrupt"))
	for _, cancel := range cancels {
		cancel()
	}
	go func() {
		poolStop <- struct{}{}
	}()

	//release connections
	dConn.Release(&pool)
}

func dispatch(ctx context.Context, topic string) {
	kafkaChan := make(chan KafkaBroker.Message)
	reader := KafkaBroker.NewReader([]string{"localhost:9092"}, topic, "groupId-1")
	go reader.Read(context.Background(), &kafkaChan, nil)
	cnn := <-pool
	defer func() {
		_ = reader.CloseReader()
		_ = cnn.KafkaWriter.CloseWriter()
		close(kafkaChan)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-kafkaChan:
			fmt.Println(msg)
		}
	}
}

func fillPool() {
	var cancels []context.CancelFunc
	count := 10
	for i := 0; i < count; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("pool worker terminated")
					return
				default:
					c, err := dConn.Fetch()
					if err != nil {
						fmt.Println("cannot get new connection, ", err)
						continue
					}
					pool <- *c
					fmt.Printf("add connection id %s to pool.\n", c.Id.String())
				}
			}
		}(ctx)
		cancels = append(cancels, cancel)
	}

	<-poolStop

	for _, c := range cancels {
		c()
	}
}