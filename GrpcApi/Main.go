package main

import (
	"fmt"
	"github.com/mhthrh/BlueBank/Config"
	"github.com/mhthrh/BlueBank/GrpcApi/GrpcServer"
	"github.com/mhthrh/BlueBank/Pool"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoGateway"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoUser"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

type server struct {
	address string
	lis     net.Listener
	gLis    *grpc.Server
}

var (
	pool        chan Pool.Connection
	gConn       Pool.IConnection
	poolStop    chan struct{}
	osInterrupt chan os.Signal
	add         chan server
	remove      chan server
	newAdd      chan string
	servers     []server
	poolSize    int
)

func init() {
	poolStop = make(chan struct{})
	remove = make(chan server)
	add = make(chan server)
	osInterrupt = make(chan os.Signal)
	newAdd = make(chan string)
	//err := os.Chdir("..")
	//if err != nil {
	//	log.Fatal(err)
	//}
	cfg := Config.New("Coded.dat", "json", "./ConfigFiles")
	if err := cfg.Initialize(); err != nil {
		log.Fatalf("unable to fill viber, %v", err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		viper.GetString("Postgresql.Host"), viper.GetInt("Postgresql.Port"), viper.GetString("Postgresql.User"), viper.GetString("Postgresql.password"), viper.GetString("Postgresql.Dbname"))

	gConn = Pool.NewGrpcConnection(psqlInfo, Pool.RedisConnection{
		Ip:       viper.GetString("Redis.Host"),
		Port:     viper.GetInt("Redis.Port"),
		Password: viper.GetString("Redis.password"),
		Database: viper.GetInt("Redis.database"),
	})
	poolSize = viper.GetInt("ConnectionPoolCount")
	pool = make(chan Pool.Connection, poolSize)
}
func main() {
	//Connection pool
	go func() {
		defer close(poolStop)
		for {
			select {
			case <-poolStop:
				return
			default:
				c, err := gConn.Fetch()
				if err != nil {
					fmt.Println("cannot get newAdd connection, ", err)
					continue
				}
				pool <- *c
				fmt.Printf("add connection id %s to pool.\n", c.Id.String())
			}
		}
	}()

	for _, address := range viper.Get("GRPC").([]interface{}) {
		ip := address.(map[string]interface{})["ip"]
		port := address.(map[string]interface{})["port"]
		go addServer(add, remove, newAdd, fmt.Sprintf("%s:%.f", ip, port))
	}

	go func() {
		for {
			select {
			case srv := <-add:
				servers = append(servers, srv)
			case srv := <-remove:
				for i, s := range servers {
					if s.address == srv.address {
						servers = append(servers[:i], servers[i+1:]...)
					}
				}
			case newAddress := <-newAdd:
				addServer(add, remove, newAdd, newAddress)
			}

		}
	}()
	go signal.Notify(osInterrupt, os.Interrupt)
	select {
	case <-osInterrupt:
		fmt.Println(fmt.Sprintf("server will be down, because: %s", "got intrrupt"))
	}
	go func() {
		poolStop <- struct{}{}
	}()
	gConn.Release(&pool)

	for _, s := range servers {
		s.gLis.Stop()
		s.lis.Close()
	}
}

func addServer(addServer chan server, removeServer chan server, newServer chan string, address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("failed to listen: %v \n", err)
	}
	rpcServer := grpc.NewServer()
	GrpcServer.New(&pool)

	ProtoUser.RegisterServicesServer(
		rpcServer, &GrpcServer.UserServer{
			UnimplementedServicesServer: ProtoUser.UnimplementedServicesServer{},
		},
	)
	ProtoGateway.RegisterGatewayServicesServer(
		rpcServer, &GrpcServer.GatewayServer{
			UnimplementedGatewayServicesServer: ProtoGateway.UnimplementedGatewayServicesServer{},
		},
	)

	log.Printf("serving on %s\n", address)
	addServer <- server{
		address: address,
		lis:     lis,
		gLis:    rpcServer,
	}
	if err := rpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v \n", err)
	}
	removeServer <- server{
		address: address,
		lis:     lis,
		gLis:    rpcServer,
	}
	newServer <- address
}
