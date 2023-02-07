package infrastructure

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/NatBrian/Stockbit-Golang-Challenge/application"
	"github.com/NatBrian/Stockbit-Golang-Challenge/kafka"
	"github.com/rs/zerolog/log"
)

// ConsumeMessages starting consumer to consume the messages
func ConsumeMessages(app application.App) {
	var dep = application.SetupDependency(app)
	log.Info().Msg("Consumer running")

	spoolUpConsumer(app, app.Config.KafkaConfig.Topic, dep.StockController.ConsumeRecords)
}

func spoolUpConsumer(app application.App, topic string, handler func(msg kafka.Message) error) {
	go func() { consume(app, topic, handler) }()
}

func consume(app application.App, topic string, handler func(kafka.Message) error) {
	reader := app.Kafka.Consumer
	defer reader.Close()
	log.Info().Msg(fmt.Sprintf("reading kafka message topic:%s", topic))

	for {
		msg, err := reader.Consume(topic)
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("error reading kafka message, exiting consumer, topic:%s", topic))
			continue
		}
		dataStr := fmt.Sprintf("topic:[%s], partition:[%d], offset:[%d]", topic, msg.Partition, msg.Offset)

		success := false
		attempt := 1
		for !success {
			err = handler(msg)
			if err != nil {
				if retry(attempt, app.Config.KafkaConfig.ConsumerMaxAttempt) {
					attempt++
					log.Error().Err(err).Msg("error has occurred, retrying, " + dataStr)
					continue
				}

				log.Error().Err(err).Msg("kafka consumer retry has reached max attempt, " + dataStr)
				break
			}
			success = true
			log.Info().Msg("message has been successfuly processed, " + dataStr)
		}

		if !success {
			log.Error().Err(err).Msg("error when processing kafka message, " + dataStr)
		}
	}
}

func retry(attempt int, maxAttempt int) bool {
	if attempt <= maxAttempt {
		// Add some randomness to prevent creating a Thundering Herd
		random, _ := rand.Int(rand.Reader, big.NewInt(int64(time.Millisecond)*5))
		jitter := time.Duration(random.Int64()) * time.Millisecond
		time.Sleep(jitter)
		return true
	}
	return false
}
