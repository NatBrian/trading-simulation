package config

import (
	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

type (
	// Config contains app configs
	Config struct {
		Constants   Constants
		KafkaConfig KafkaConfig
		RedisConfig RedisConfig
	}

	Constants struct {
		HTTPPort string `mapstructure:"HTTP_PORT"`
		HTTPHost string `mapstructure:"HTTP_HOST"`
		MaxZipMb int    `mapstructure:"MAX_ZIP_MB"`
	}

	KafkaConfig struct {
		Address            string
		Topic              string
		ConsumerMaxAttempt int
	}

	RedisConfig struct {
		Address string
	}
)

func LoadConfig() (Config, error) {
	log.Info().Msg("Load Configs")

	constants, err := loadConstants()
	if err != nil {
		log.Error().Err(err).Msg("error: loadConstants")
		return Config{}, err
	}

	KafkaConfig, err := loadKafkaConfig()
	if err != nil {
		log.Error().Err(err).Msg("error: loadKafkaConfig")
		return Config{}, err
	}

	redisConfig, err := loadRedisConfig()
	if err != nil {
		log.Error().Err(err).Msg("error: loadRedisConfig")
		return Config{}, err
	}

	return Config{
		Constants:   constants,
		KafkaConfig: KafkaConfig,
		RedisConfig: redisConfig,
	}, nil
}

func loadConstants() (Constants, error) {
	viper.SetConfigType("env")
	viper.AddConfigPath("config/")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Error().Err(err).Msg("error: viper ReadInConfig")
		return Constants{}, err
	}

	return Constants{
		HTTPPort: viper.GetString("HTTP_PORT"),
		HTTPHost: viper.GetString("HTTP_HOST"),
		MaxZipMb: viper.GetInt("MAX_ZIP_MB"),
	}, nil
}

func loadKafkaConfig() (KafkaConfig, error) {
	viper.SetConfigType("env")
	viper.AddConfigPath("config/")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Error().Err(err).Msg("error: viper ReadInConfig")
		return KafkaConfig{}, err
	}

	return KafkaConfig{
		Address:            viper.GetString("KAFKA_ADDRESS"),
		Topic:              viper.GetString("KAFKA_TOPIC"),
		ConsumerMaxAttempt: viper.GetInt("KAFKA_CONSUMER_MAX_ATTEMPT"),
	}, nil
}

func loadRedisConfig() (RedisConfig, error) {
	viper.SetConfigType("env")
	viper.AddConfigPath("config/")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Error().Err(err).Msg("error: viper ReadInConfig")
		return RedisConfig{}, err
	}

	return RedisConfig{
		Address: viper.GetString("REDIS_ADDRESS"),
	}, nil
}
