package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"hsm/configs"
	grpc_server "hsm/pkg/crypto/v1"
)

const (
	configPath = "./configs"
	configName = "dev"
)

func main() {
	// read config
	conf := configs.LoadConfig(configPath, configName)
	if conf == nil {
		panic("failed to load config")
	}
	fmt.Printf("config: %+v\n", *conf)

	// grpc server
	g := grpc_server.NewServer(conf)
	g.Start()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign
	log.Println("server is exiting....")
	g.Stop()
}
