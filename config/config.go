package config

import (
	"log"

	"github.com/spf13/viper"
)

type (
	// Config contains app configs
	Config struct {
		HTTPPort string `mapstructure:"HTTP_PORT"`
		HTTPHost string `mapstructure:"HTTP_HOST"`
	}
)

func LoadConfig() (Config, error) {
	log.Println("Load Configs")

	viper.SetConfigType("env")
	viper.AddConfigPath("config/")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("error: viper ReadInConfig", err)
		return Config{}, err
	}

	return Config{
		HTTPPort: viper.GetString("HTTP_PORT"),
		HTTPHost: viper.GetString("HTTP_HOST"),
	}, nil
}
