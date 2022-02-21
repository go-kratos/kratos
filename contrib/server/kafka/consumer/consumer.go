package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/contrib/server/kafka"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/Shopify/sarama"
)

const (
	defaultVersion     = "2.1.0"
	defaultMaxProcTime = 250 * time.Millisecond
)

type Message struct {
	key   string
	value []byte
}

func (m *Message) Key() string {
	return m.key
}

func (m *Message) Value() []byte {
	return m.value
}

type Consumer struct {
	// config setting
	brokers         []string
	topic           string
	group           string
	version         string
	balanceStrategy sarama.BalanceStrategy
	initialOffset   int64
	maxProcTime     time.Duration

	ready chan struct{}

	sarama.Client
	consumer sarama.ConsumerGroup
	handler  kafka.Handler
	logger   *log.Helper
}

// ConsumerOption is a Consumer config option.
type ConfigOption func(*Consumer)

// Version with Kafka cluster version.
func Version(version string) ConfigOption {
	return func(c *Consumer) {
		c.version = version
	}
}

// BalanceStrategy with the rebalance strategy
func BalanceStrategy(strategy sarama.BalanceStrategy) ConfigOption {
	return func(c *Consumer) {
		c.balanceStrategy = strategy
	}
}

// InitialOffset with the consumer initial offset
func InitialOffset(offset int64) ConfigOption {
	return func(c *Consumer) {
		c.initialOffset = offset
	}
}

// MaxProcessingTime with the maximum amount of time the consumer expects a message takes
func MaxProcessingTime(maxTime time.Duration) ConfigOption {
	return func(c *Consumer) {
		c.maxProcTime = maxTime
	}
}

// Logger with the specify logger
func Logger(logger log.Logger) ConfigOption {
	return func(c *Consumer) {
		c.logger = log.NewHelper(logger)
	}
}

// NewConsumer inits a consumer group consumer
func NewConsumer(brokers []string, topic, group string, opts ...ConfigOption) (*Consumer, error) {
	result := &Consumer{
		brokers:         brokers,
		topic:           topic,
		group:           group,
		version:         defaultVersion,
		balanceStrategy: sarama.BalanceStrategySticky,
		initialOffset:   sarama.OffsetNewest,
		maxProcTime:     defaultMaxProcTime,
		ready:           make(chan struct{}),
		logger:          log.NewHelper(log.GetLogger()),
	}

	for _, o := range opts {
		o(result)
	}

	// check version
	kafkaVersion, err := sarama.ParseKafkaVersion(result.version)
	if err != nil {
		return nil, fmt.Errorf("parse kafka version %s error: %v", result.version, err)
	}

	// init config
	config := sarama.NewConfig()
	config.Version = kafkaVersion
	config.Consumer.MaxProcessingTime = result.maxProcTime
	config.Consumer.Group.Rebalance.Strategy = result.balanceStrategy
	config.Consumer.Offsets.Initial = result.initialOffset
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("init sarama kafka client error: %v", err)
	}
	consumer, err := sarama.NewConsumerGroupFromClient(group, client)
	if err != nil {
		return nil, fmt.Errorf("init sarama kafka group consumer error: %v", err)
	}
	result.Client = client
	result.consumer = consumer

	return result, nil
}

// Topics returns all the topics this consumer subscribes
func (c *Consumer) Topic() string {
	return c.topic
}

// RegisterHandler registers a handler to handle the messages of a specific topic
func (c *Consumer) RegisterHandler(handler kafka.Handler) {
	c.handler = handler
}

// RegisterHandler checks whether this consumer has a handler for the specific topic
func (c *Consumer) HasHandler() bool {
	return c.handler != nil
}

// Consume receives and handles messages
func (c *Consumer) Consume(ctx context.Context) error {
	// check handlers before consuming
	for c.handler == nil {
		return fmt.Errorf("consumer has no handler")
	}

	// start consuming
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			c.logger.Infof("consumer %+v session starts", c)

			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := c.consumer.Consume(ctx, []string{c.topic}, c); err != nil {
				c.logger.Errorf("consumer %+v consumes error %+v", c, err)
				return
			}

			c.logger.Infof("consumer %+v session exits", c)

			// check if context was cancelled, signaling that the consumer should stop
			if err := ctx.Err(); err != nil {
				c.logger.Errorf("consumer %+v exits due to context is canceled %+v", c, err)
				return
			}

			c.ready = make(chan struct{})
		}
	}()

	<-c.ready // Await till the consumer has been set up

	wg.Wait()
	if err := c.consumer.Close(); err != nil {
		return fmt.Errorf("close kafka consumer error: %v", err)
	}

	return fmt.Errorf("consumer %+v exited", c.consumer)
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	c.logger.Infof("consumer of group %s setup status %+v", c.group, session.Claims())
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	c.logger.Infof("consumer of group %s exit status %+v", c.group, session.Claims())
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		if err := c.handler.Handle(&Message{key: string(message.Key), value: message.Value}); err != nil {
			c.logger.Errorf("consume message %s of topic %s partition %d error %+v", string(message.Value), message.Topic, message.Partition, err)
			continue
		}
		session.MarkMessage(message, "")
	}

	return nil
}
