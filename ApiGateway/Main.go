package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mhthrh/BlueBank/ApiGateway/View"
	"github.com/mhthrh/BlueBank/Config"
	"github.com/mhthrh/BlueBank/Pool"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	readTimeOut  = 10 * time.Second
	WriteTimeOut = 10 * time.Second
	idleTimeOut  = 180 * time.Second
	poolCount    = 5
)

var (
	listenerError chan error
	osInterrupt   chan os.Signal
	poolStop      chan struct{}
	pool          chan Pool.Connection
	rConn         Pool.IConnection
	poolSize      int
)

func init() {

	listenerError = make(chan error)
	osInterrupt = make(chan os.Signal)
	poolStop = make(chan struct{})
	pool = make(chan Pool.Connection, poolCount)

	cfg := Config.New("Coded.dat", "json", "ConfigFile")
	if err := cfg.Initialize(); err != nil {
		log.Fatalf("unable to fill viber, %v", err)
	}
	rConn = Pool.NewRestSide(Pool.RedisConnection{
		Ip:       viper.GetString("Redis.Host"),
		Port:     viper.GetInt("Redis.Port"),
		Password: viper.GetString("Redis.password"),
		Database: viper.GetInt("Redis.database"),
	}, Pool.KafkaAddress{
		Ip:   viper.GetString("Kafka.Host"),
		Port: viper.GetInt("Kafka.Port"),
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
				c, err := rConn.Fetch()
				if err != nil {
					fmt.Println("cannot get newAdd connection, ", err)
					continue
				}
				pool <- *c
				fmt.Printf("add connection id %s to pool.\n", c.Id.String())
			}
		}
	}()

	View.New(&pool)
	serverSync := http.Server{
		Addr:         fmt.Sprintf("%s:%d", "localhost", 8569),
		Handler:      View.RunSync(),
		TLSConfig:    nil,
		ReadTimeout:  readTimeOut,
		WriteTimeout: WriteTimeOut,
		IdleTimeout:  idleTimeOut,
	}
	serverAsync := http.Server{
		Addr:         fmt.Sprintf("%s:%d", "localhost", 8570),
		Handler:      View.RunAsync(),
		TLSConfig:    nil,
		ReadTimeout:  readTimeOut,
		WriteTimeout: WriteTimeOut,
		IdleTimeout:  idleTimeOut,
	}
	go signal.Notify(osInterrupt, os.Interrupt)
	go func() {
		if err := serverSync.ListenAndServe(); err != http.ErrServerClosed {
			listenerError <- errors.New(fmt.Sprintf("HTTP serverSync ListenAndServe: %v", err))
			return
		}
	}()
	go func() {
		if err := serverAsync.ListenAndServe(); err != http.ErrServerClosed {
			listenerError <- errors.New(fmt.Sprintf("HTTP serverAsync ListenAndServe: %v", err))
			return
		}
	}()
	//Listener for fatal error
	select {
	case listenerErrorMessage := <-listenerError:
		fmt.Println(fmt.Sprintf("serverSync will be down, because: %s", listenerErrorMessage.Error()))
		if err := serverSync.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP serverSync Shutdown: %v", err)
		}
		if err := serverAsync.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP serverAsync Shutdown: %v", err)
		}
	case <-osInterrupt:
		fmt.Println(fmt.Sprintf("serverSync will be down, because: %s", "got intrrupt"))
		if err := serverSync.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP serverSync Shutdown: %v", err)
		}
		if err := serverAsync.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP serverAsync Shutdown: %v", err)
		}
	case e := <-listenerError:
		fmt.Println(fmt.Sprintf("serverSync will be down, because: %s", "got intrrupt"))
		if err := serverSync.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP serverSync Shutdown: %v", e)
		}
		if err := serverAsync.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP serverAsync Shutdown: %v", err)
		}
		go func() {
			poolStop <- struct{}{}
		}()

		//release connections
		rConn.Release(&pool)
	}
}
