package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Http      HttpConfig      `mapstructure:"server"`
	DB        DbConfig        `mapstructure:"database"`
	Auth      AuthConfig      `mapstructure:"auth"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetDefault("rate_limit.requests_per_second", 5.0)
	viper.SetDefault("rate_limit.burst", 10)
	viper.SetDefault("rate_limit.expires_in_minutes", 10)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Error reading config file, %s", err))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("Unable to decode into struct, %v", err))
	}
	config.RateLimit = config.RateLimit.WithDefaults()

	return &config
}
