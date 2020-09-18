package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig(path, name string) *Config {
	conf := &Config{}
	vp := viper.New()
	vp.AddConfigPath(path)
	vp.SetConfigName(name)
	if err := vp.ReadInConfig(); err != nil {
		fmt.Printf("failed to read config: %v", err)
		return nil
	}

	if err := vp.Unmarshal(conf); err != nil {
		fmt.Printf("failed to unmarshal config: %v", err)
		return nil
	}

	return conf
}
