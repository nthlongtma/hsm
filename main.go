package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"hsm/configs"
	hsm_api "hsm/pkg/hsm-api"
	http_server "hsm/pkg/http-server"
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

	// http server
	s := http_server.NewServer(conf, ctx, ss)
	go s.Start()
	defer s.Stop()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign
	log.Println("server is exiting....")
}
