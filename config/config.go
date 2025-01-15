package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Http HttpConfig `mapstructure:"server"`
	DB   DbConfig   `mapstructure:"database"`
	Auth AuthConfig `mapstructure:"auth"`
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Error reading config file, %s", err))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("Unable to decode into struct, %v", err))
	}

	return &config
}
