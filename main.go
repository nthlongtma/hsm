package main

import (
	"fmt"
	"hsm/configs"
	hsm_client "hsm/pkg/hsm-client"
	"hsm/pkg/server"
	"net/http"
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
	ctx, err := hsm_client.GetContext(conf.ModulePath)
	if err != nil {
		panic(err)
	}
	defer hsm_client.FinishContext(ctx)

	ss, err := hsm_client.GetSession(ctx, conf.HSM.SlotID, conf.HSM.Pin)
	if err != nil {
		panic(err)
	}
	defer hsm_client.FinishSession(ctx, ss)

	hsm := hsm_client.NewHSM(ctx, ss)

	// server
	s := server.NewServer(conf, hsm)

	http.HandleFunc("/api/v1/encrypt", s.HandleEncrypt)
	http.HandleFunc("/api/v1/decrypt", s.HandleDecrypt)

	if err := http.ListenAndServe(":8888", nil); err != nil {
		panic(err)
	}
}
