package main

import (
	"github.com/mhthrh/BlueBank/Config"
	"github.com/mhthrh/BlueBank/Pool"
	"log"
	"os"
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

}
func main() {

}
