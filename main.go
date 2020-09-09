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
	// conf := LoadConfig(configPath, configName)
	// if conf == nil {
	// 	panic("failed to load config")
	// }
	// fmt.Printf("config: %+v\n", *conf)

	a, b := 5, 8
	if a&1<<0 == 1 {
		fmt.Print("a is odd")
	} else {
		fmt.Print("a is even")
	}

	if b&1<<0 == 1 {
		fmt.Print("b is odd")
	} else {
		fmt.Print("b is even")
	}

}
