package main

import (
	"fmt"
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
}
