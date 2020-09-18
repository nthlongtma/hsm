package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	configPath = "./configs"
	configName = "dev"
)

func main() {
	// read config
	conf := LoadConfig(configPath, configName)
	if conf == nil {
		panic("failed to load config")
	}
	fmt.Printf("config: %+v\n", *conf)

	// init hsm
	ctx, err := GetContext(conf.ModulePath)
	if err != nil {
		panic(err)
	}
	defer FinishContext(ctx)

	ss, err := GetSession(ctx, conf.HSM.SlotID, conf.HSM.Pin)
	if err != nil {
		panic(err)
	}
	defer FinishSession(ctx, ss)

	hsm := NewHSM(ctx, ss)

	// server
	s := NewServer(conf, hsm)

	http.HandleFunc("/api/v1/encrypt", s.HandleEncrypt)
	http.HandleFunc("/api/v1/decrypt", s.HandleDecrypt)

	if err := http.ListenAndServe(":8888", nil); err != nil {
		panic(err)
	}
}

func ToJson(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
