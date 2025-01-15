package config

type AuthConfig struct {
	AccessSecret  string `mapstructure:"access_secret"`
	RefreshSecret string `mapstructure:"refresh_secret"`
}
