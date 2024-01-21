package config

import "time"

type AppConfig struct {
	Name                    string        `mapstructure:"APP_NAME"`
	Host                    string        `mapstructure:"APP_HOST"`
	Port                    string        `mapstructure:"APP_PORT"`
	GracefulShutdownTimeout time.Duration `mapstructure:"APP_GRACEFUL_SHUTDOWN_TIMEOUT"`
}
