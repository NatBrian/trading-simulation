package kafka

import (
	"context"
	"sync"
	"time"

	kafkaGo "github.com/segmentio/kafka-go"
)

const (
	LatestOffset   int64 = -1 // The most recent offset available for a partition. set by the lib
	EarliestOffset int64 = -2 // The least recent offset available for a partition set by the lib.
)

var consumers sync.Map

type Consumer interface {
	Consume(topic string) (Message, error)
	Fetch(topic string) (Message, error)
	Commit(topic string, message Message) error
	Close()
}

type KafkaGoConsumer struct {
	ctx       context.Context
	consumers sync.Map
	config    *ConsumerConfig
}

type ConsumerConfig struct {
	// The list of broker addresses used to connect to the kafka cluster.
	Brokers []string

	// GroupID holds the optional consumer group id.  If GroupID is specified, then
	// Partition should NOT be specified e.g. 0
	GroupID string

	// The capacity of the internal message queue, defaults to 100 if none is
	// set.
	QueueCapacity int

	// Min and max number of bytes to fetch from kafka in each request.
	MinBytes int
	MaxBytes int

	// Maximum amount of time to wait for new data to come when fetching batches
	// of messages from kafka.
	MaxWait time.Duration

	// HeartbeatInterval sets the optional frequency at which the reader sends the consumer
	// group heartbeat update.
	//
	// Default: 3s
	//
	// Only used when GroupID is set
	HeartbeatInterval time.Duration

	// CommitInterval indicates the interval at which offsets are committed to
	// the broker.  If 0, commits will be handled synchronously.
	//
	// Default: 0
	//
	// Only used when GroupID is set
	CommitInterval time.Duration

	// PartitionWatchInterval indicates how often a reader checks for partition changes.
	// If a reader sees a partition change (such as a partition add) it will rebalance the group
	// picking up new partitions.
	//
	// Default: 5s
	//
	// Only used when GroupID is set and WatchPartitionChanges is set.
	PartitionWatchInterval time.Duration

	// WatchForPartitionChanges is used to inform kafka-go that a consumer group should be
	// polling the brokers and rebalancing if any partition changes happen to the topic.
	WatchPartitionChanges bool

	// SessionTimeout optionally sets the length of time that may pass without a heartbeat
	// before the coordinator considers the consumer dead and initiates a rebalance.
	//
	// Default: 30s
	//
	// Only used when GroupID is set
	SessionTimeout time.Duration

	// RebalanceTimeout optionally sets the length of time the coordinator will wait
	// for members to join as part of a rebalance.  For kafka servers under higher
	// load, it may be useful to set this value higher.
	//
	// Default: 30s
	//
	// Only used when GroupID is set
	RebalanceTimeout time.Duration

	// JoinGroupBackoff optionally sets the length of time to wait between re-joining
	// the consumer group after an error.
	//
	// Default: 5s
	JoinGroupBackoff time.Duration

	// RetentionTime optionally sets the length of time the consumer group will be saved
	// by the broker
	//
	// Default: 24h
	//
	// Only used when GroupID is set
	RetentionTime time.Duration

	// StartOffset determines from whence the consumer group should begin
	// consuming when it finds a partition without a committed offset.  If
	// non-zero, it must be set to one of FirstOffset or LastOffset.
	//
	// Default: FirstOffset
	//
	// Only used when GroupID is set
	StartOffset int64

	// BackoffDelayMin optionally sets the smallest amount of time the reader will wait before
	// polling for new messages
	//
	// Default: 100ms
	ReadBackoffMin time.Duration

	// BackoffDelayMax optionally sets the maximum amount of time the reader will wait before
	// polling for new messages
	//
	// Default: 1s
	ReadBackoffMax time.Duration

	// Limit of how many attempts will be made before delivering the error.
	//
	// The default is to try 3 times.
	MaxAttempts int
}

func NewConsumer(ctx context.Context, config *ConsumerConfig) *KafkaGoConsumer {
	k := new(KafkaGoConsumer)
	k.ctx = ctx
	k.config = config
	return k
}

// Close is
func (k *KafkaGoConsumer) Close() {
	var wg sync.WaitGroup
	defer wg.Wait()

	k.consumers.Range(
		func(topic, consumer interface{}) bool {
			if consumer != nil {
				wg.Add(1)

				go func() {
					consumer.(*kafkaGo.Reader).Close()
					wg.Done()
				}()
			}

			consumers.Delete(topic)
			return true
		})
}

// Consume will get next message and commit it immeduiately
func (k *KafkaGoConsumer) Consume(topic string) (Message, error) {
	consumer := k.getConsumer(topic)
	if msg, err := consumer.ReadMessage(k.ctx); err != nil {
		return Message{}, err
	} else {
		headers := make(map[string][]byte)
		for _, header := range msg.Headers {
			headers[header.Key] = header.Value
		}
		return Message{
			Value:     msg.Value,
			Headers:   headers,
			Partition: msg.Partition,
			Offset:    msg.Offset,
			Key:       msg.Key,
		}, nil
	}
}

// Fetch will get next nessage but will not commit it
// see also : Commit(topic string, message Message)
func (k *KafkaGoConsumer) Fetch(topic string) (Message, error) {
	consumer := k.getConsumer(topic)
	if msg, err := consumer.FetchMessage(k.ctx); err != nil {
		return Message{}, err
	} else {
		headers := make(map[string][]byte)
		for _, header := range msg.Headers {
			headers[header.Key] = header.Value
		}
		return Message{
			Value:     msg.Value,
			Headers:   headers,
			Partition: msg.Partition,
			Offset:    msg.Offset,
			Key:       msg.Key,
		}, nil
	}
}

// Commit will commit message. only need topic, partition and offset.
func (k *KafkaGoConsumer) Commit(topic string, message Message) error {
	consumer := k.getConsumer(topic)

	msg := kafkaGo.Message{ // only needs these three
		Topic:     topic,
		Partition: message.Partition,
		Offset:    message.Offset,
	}

	return consumer.CommitMessages(k.ctx, msg)
}

// GetConsumer is
func (k *KafkaGoConsumer) getConsumer(topic string) *kafkaGo.Reader {
	var consumer (*kafkaGo.Reader)

	if value, ok := consumers.Load(topic); ok {
		consumer = value.(*kafkaGo.Reader)
	} else {
		//initialize consumer for topic
		consumer = kafkaGo.NewReader(
			kafkaGo.ReaderConfig{
				Brokers:                k.config.Brokers,
				GroupID:                k.config.GroupID,
				Topic:                  topic,
				QueueCapacity:          k.config.QueueCapacity,
				MinBytes:               k.config.MinBytes,
				MaxBytes:               k.config.MaxBytes,
				MaxWait:                k.config.MaxWait,
				HeartbeatInterval:      k.config.HeartbeatInterval,
				CommitInterval:         k.config.CommitInterval,
				PartitionWatchInterval: k.config.PartitionWatchInterval,
				WatchPartitionChanges:  k.config.WatchPartitionChanges,
				SessionTimeout:         k.config.SessionTimeout,
				RebalanceTimeout:       k.config.RebalanceTimeout,
				JoinGroupBackoff:       k.config.JoinGroupBackoff,
				RetentionTime:          k.config.RetentionTime,
				StartOffset:            k.config.StartOffset,
				ReadBackoffMin:         k.config.ReadBackoffMax,
				ReadBackoffMax:         k.config.ReadBackoffMin,
				MaxAttempts:            k.config.MaxAttempts,
			})

		consumers.Store(topic, consumer)
	}

	return consumer
}
