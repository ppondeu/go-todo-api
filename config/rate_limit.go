package config

const (
	defaultRateLimitRequestsPerSecond = 5.0
	defaultRateLimitBurst             = 10
	defaultRateLimitExpiresInMinutes  = 10
)

type RateLimitConfig struct {
	RequestsPerSecond float64 `mapstructure:"requests_per_second"`
	Burst             int     `mapstructure:"burst"`
	ExpiresInMinutes  int     `mapstructure:"expires_in_minutes"`
}

func (c RateLimitConfig) WithDefaults() RateLimitConfig {
	if c.RequestsPerSecond <= 0 {
		c.RequestsPerSecond = defaultRateLimitRequestsPerSecond
	}
	if c.Burst <= 0 {
		c.Burst = defaultRateLimitBurst
	}
	if c.ExpiresInMinutes <= 0 {
		c.ExpiresInMinutes = defaultRateLimitExpiresInMinutes
	}

	return c
}
