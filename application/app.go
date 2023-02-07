package application

import (
	"context"

	"github.com/NatBrian/Stockbit-Golang-Challenge/config"
	"github.com/NatBrian/Stockbit-Golang-Challenge/kafka"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// App contains application instances
type App struct {
	Context context.Context
	Config  config.Config
	Kafka   KafkaDriver
	Redis   *redis.Client
}

type KafkaDriver struct {
	Producer kafka.Producer
	Consumer kafka.Consumer
}

func SetupApp(context context.Context) (App, error) {
	var app App

	log.Info().Msg("Setup App")

	app.Context = context

	loadConfig, err := config.LoadConfig()
	if err != nil {
		log.Error().Err(err).Msg("error: LoadConfig")
		return App{}, err
	}
	app.Config = loadConfig

	producer := kafka.NewProducer(context, &kafka.ProducerConfig{
		AppName:      "stock-bit",
		Brokers:      []string{app.Config.KafkaConfig.Address},
		RequiredAcks: 1,
		Async:        true,
	})

	consumer := kafka.NewConsumer(context, &kafka.ConsumerConfig{
		Brokers:     []string{app.Config.KafkaConfig.Address},
		GroupID:     "stock-bit",
		MaxAttempts: app.Config.KafkaConfig.ConsumerMaxAttempt,
	})

	app.Kafka = KafkaDriver{
		Producer: producer,
		Consumer: consumer,
	}

	app.Redis = redis.NewClient(&redis.Options{
		Addr:     app.Config.RedisConfig.Address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return app, nil
}
