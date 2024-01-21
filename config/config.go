package config

import (
	"fmt"
	"github.com/spf13/viper"
	"order-backend/pkg/logger"
	"strings"
)

type Config struct {
	App      AppConfig      `mapstructure:",squash"`
	Postgres PostgresConfig `mapstructure:",squash"`
}

func NewConfig() (*Config, error) {
	config := Config{}
	err := loadConfig(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func loadConfig(output interface{}) error {
	log := logger.Log()
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	viper.GetViper().AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(fmt.Sprintf("Error reading config file, %s", err))
	}

	err := viper.Unmarshal(&output)
	if err != nil {
		logger.Log().Error(err.Error())
		return err
	}

	return nil
}
