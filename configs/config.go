package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

var config *viper.Viper

func init() {
	config = viper.New()
	config.AddConfigPath("./configs")
	config.SetConfigFile("./configs/config.yaml")
	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func Get() *viper.Viper {
	return config
}
