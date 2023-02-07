package kafka

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/google/uuid"
	kafkaGo "github.com/segmentio/kafka-go"
)

const (
	keyMessageID = "MESSAGE_ID"
)

var (
	producers sync.Map
	balancer  = map[int]kafkaGo.Balancer{
		0: &kafkaGo.RoundRobin{},
		1: &kafkaGo.LeastBytes{},
		2: &kafkaGo.Hash{},
	}
)

type Producer interface {
	Produce(ctx context.Context, topic string, message Message) error
	Close()
}

type KafkaGoProducer struct {
	config  *ProducerConfig
	ctx     context.Context
	appname string
}

// For now producer only need brokers
type ProducerConfig struct {
	// AppName to be passed into Kafka header
	AppName string

	// The list of broker addresses used to connect to the kafka cluster.
	Brokers []string

	// The balancer used to distribute messages across partitions.
	// The default is to use round robin distribution
	// Selection: round-robin(0), least bytes(1), hash(2)
	Balancer int

	// Limit on how many attempts will be made to deliver a message.
	// If none is set, default to 10
	MaxAttempts int

	// A hint on the capacity of the writer's internal message queue.
	// If none is given, the default is to use a queue capacity of 100 messages.
	QueueCapacity int

	// Limit on how many messages will be buffered before being sent to a partition.
	// If none is given, the default is to use a target batch size of 100 messages.
	BatchSize int

	// Limit the maximum size of a request in bytes before being sent to a partition.
	// If none is given, the default is to use a kafka default value of 1048576.
	BatchBytes int

	// Time limit on how often incomplete message batches will be flushed to kafka.
	// If none is given, default is to flush at least every second.
	BatchTimeout time.Duration

	// Timeout for read operations performed by the Writer.
	// If none is given, defaults to 10 seconds.
	ReadTimeout time.Duration

	// Timeout for write operation performed by the Writer.
	// If none is given, defaults to 10 seconds.
	WriteTimeout time.Duration

	// This interval defines how often the list of partitions is refreshed from
	// kafka. It allows the writer to automatically handle when new partitions
	// are added to a topic.
	//
	// If none is given, the default is to refresh partitions every 15 seconds.
	RebalanceInterval time.Duration

	// Number of acknowledges from partition replicas required before receiving
	// a response to a produce request
	// If none is given, default to -1, which means to wait for all replicas
	RequiredAcks int

	// This flag indicating the message writing process to never block
	// By setting this true, it means that errors are ignored
	// Only use this if you don't care about the guarantee of whether the messages
	// were written
	Async bool
}

// NewProducer will create KafkaProducerInstance based on context and kafka config
func NewProducer(ctx context.Context, config *ProducerConfig) Producer {
	return &KafkaGoProducer{
		config:  config,
		ctx:     ctx,
		appname: config.AppName,
	}
}

// getProducer get list of producer for a specific topic
func (k *KafkaGoProducer) getProducer(topic string) *kafkaGo.Writer {
	var producer *kafkaGo.Writer

	if v, ok := producers.Load(topic); ok {
		producer = v.(*kafkaGo.Writer)
	} else {
		producer = kafkaGo.NewWriter(
			kafkaGo.WriterConfig{
				Brokers:           k.config.Brokers,
				Topic:             topic,
				Balancer:          balancer[k.config.Balancer],
				MaxAttempts:       k.config.MaxAttempts,
				QueueCapacity:     k.config.QueueCapacity,
				BatchSize:         k.config.BatchSize,
				BatchBytes:        k.config.BatchBytes,
				BatchTimeout:      k.config.BatchTimeout,
				ReadTimeout:       k.config.ReadTimeout,
				WriteTimeout:      k.config.WriteTimeout,
				RebalanceInterval: k.config.RebalanceInterval,
				RequiredAcks:      k.config.RequiredAcks,
				Async:             k.config.Async,
			})

		// Add producer to list of known producers by topic
		producers.Store(topic, producer)
	}

	return producer
}

func (k *KafkaGoProducer) transformHeaders(raw map[string][]byte) []kafkaGo.Header {
	var headers []kafkaGo.Header

	haveMessageID := false
	for k, v := range raw {
		k = strings.ToUpper(k)
		headers = append(headers, kafkaGo.Header{Key: k, Value: v})
		if k == keyMessageID {
			haveMessageID = true
		}
	}

	headers = append(headers, kafkaGo.Header{Key: "app_name", Value: []byte(k.appname)})
	if !haveMessageID {
		headers = append(headers, kafkaGo.Header{Key: keyMessageID, Value: []byte(uuid.New().String())})
	}

	return headers
}

// Produce do the producing a message to specific topic
func (k *KafkaGoProducer) Produce(ctx context.Context, topic string, message Message) error {
	messageInfo := fmt.Sprintf("Produce(%s)", topic)
	log.Info().Msg(messageInfo)

	//transform primitive map to kafkago header
	headers := k.transformHeaders(message.Headers)

	producer := k.getProducer(topic)
	// Write message
	err := producer.WriteMessages(
		ctx,

		kafkaGo.Message{
			Headers: headers,
			Value:   message.Value,
			Key:     message.Key,
		},
	)

	return err
}

// Close all known producer
func (k KafkaGoProducer) Close() {
	var wg sync.WaitGroup
	defer wg.Wait()

	producers.Range(
		func(topic, producer interface{}) bool {
			messageInfo := fmt.Sprintf("Close Producer %s", topic)
			log.Info().Msg(messageInfo)

			if producer != nil {
				wg.Add(1)

				go func() {
					// Close the writer
					_ = producer.(*kafkaGo.Writer).Close()
					wg.Done()
				}()
			}

			// Remove producer from list producer
			producers.Delete(topic)
			return true
		})
}
