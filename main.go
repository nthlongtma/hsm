package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"hsm/configs"
	grpc_server "hsm/pkg/crypto/v1"
	hsm_api "hsm/pkg/hsm-api"
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
	// init hsm
	ctx, err := hsm_api.GetContext(conf.ModulePath)
	if err != nil {
		panic(err)
	}

	ss, err := hsm_api.GetSession(ctx, conf.HSM.SlotID, conf.HSM.Pin)
	if err != nil {
		panic(err)
	}

	// // http server
	// s := http_server.NewServer(conf, ctx, ss)
	// go s.Start()

	// grpc server
	g := grpc_server.NewServer(conf, ctx, ss)
	g.Start()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign
	log.Println("server is exiting....")
	g.Stop()
}
