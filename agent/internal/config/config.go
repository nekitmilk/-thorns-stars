package config

import (
	"time"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	MonitoringCenterURL string        `env:"MONITORING_CENTER_URL" default:"http://localhost:8080"`
	HostID              string        `env:"HOST_ID" required:"true"`
	PollingInterval     time.Duration `env:"POLLING_INTERVAL" default:"5m"`
	RequestTimeout      time.Duration `env:"REQUEST_TIMEOUT" default:"30s"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
