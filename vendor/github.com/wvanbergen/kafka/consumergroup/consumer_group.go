package consumergroup

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kazoo-go"
)

var (
	AlreadyClosing = errors.New("The consumer group is already shutting down.")
)

type Config struct {
	*sarama.Config

	Zookeeper *kazoo.Config

	Offsets struct {
		Initial           int64         // The initial offset method to use if the consumer has no previously stored offset. Must be either sarama.OffsetOldest (default) or sarama.OffsetNewest.
		ProcessingTimeout time.Duration // Time to wait for all the offsets for a partition to be processed after stopping to consume from it. Defaults to 1 minute.
		CommitInterval    time.Duration // The interval between which the processed offsets are commited.
		ResetOffsets      bool          // Resets the offsets for the consumergroup so that it won't resume from where it left off previously.
	}
}

func NewConfig() *Config {
	config := &Config{}
	config.Config = sarama.NewConfig()
	config.Zookeeper = kazoo.NewConfig()
	config.Offsets.Initial = sarama.OffsetOldest
	config.Offsets.ProcessingTimeout = 60 * time.Second
	config.Offsets.CommitInterval = 10 * time.Second

	return config
}

func (cgc *Config) Validate() error {
	if cgc.Zookeeper.Timeout <= 0 {
		return sarama.ConfigurationError("ZookeeperTimeout should have a duration > 0")
	}

	if cgc.Offsets.CommitInterval < 0 {
		return sarama.ConfigurationError("CommitInterval should have a duration >= 0")
	}

	if cgc.Offsets.Initial != sarama.OffsetOldest && cgc.Offsets.Initial != sarama.OffsetNewest {
		return errors.New("Offsets.Initial should be sarama.OffsetOldest or sarama.OffsetNewest.")
	}

	if cgc.Config != nil {
		if err := cgc.Config.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// The ConsumerGroup type holds all the information for a consumer that is part
// of a consumer group. Call JoinConsumerGroup to start a consumer.
type ConsumerGroup struct {
	config *Config

	consumer sarama.Consumer
	kazoo    *kazoo.Kazoo
	group    *kazoo.Consumergroup
	instance *kazoo.ConsumergroupInstance

	wg             sync.WaitGroup
	singleShutdown sync.Once

	messages chan *sarama.ConsumerMessage
	errors   chan error
	stopper  chan struct{}

	consumers kazoo.ConsumergroupInstanceList

	offsetManager OffsetManager
}

// Connects to a consumer group, using Zookeeper for auto-discovery
func JoinConsumerGroup(name string, topics []string, zookeeper []string, config *Config) (cg *ConsumerGroup, err error) {

	if name == "" {
		return nil, sarama.ConfigurationError("Empty consumergroup name")
	}

	if len(topics) == 0 {
		return nil, sarama.ConfigurationError("No topics provided")
	}

	if len(zookeeper) == 0 {
		return nil, errors.New("You need to provide at least one zookeeper node address!")
	}

	if config == nil {
		config = NewConfig()
	}
	config.ClientID = name

	// Validate configuration
	if err = config.Validate(); err != nil {
		return
	}

	var kz *kazoo.Kazoo
	if kz, err = kazoo.NewKazoo(zookeeper, config.Zookeeper); err != nil {
		return
	}

	brokers, err := kz.BrokerList()
	if err != nil {
		kz.Close()
		return
	}

	group := kz.Consumergroup(name)

	if config.Offsets.ResetOffsets {
		err = group.ResetOffsets()
		if err != nil {
			kz.Close()
			return
		}
	}

	instance := group.NewInstance()

	var consumer sarama.Consumer
	if consumer, err = sarama.NewConsumer(brokers, config.Config); err != nil {
		kz.Close()
		return
	}

	cg = &ConsumerGroup{
		config:   config,
		consumer: consumer,

		kazoo:    kz,
		group:    group,
		instance: instance,

		messages: make(chan *sarama.ConsumerMessage, config.ChannelBufferSize),
		errors:   make(chan error, config.ChannelBufferSize),
		stopper:  make(chan struct{}),
	}

	// Register consumer group
	if exists, err := cg.group.Exists(); err != nil {
		cg.Logf("FAILED to check for existence of consumergroup: %s!\n", err)
		_ = consumer.Close()
		_ = kz.Close()
		return nil, err
	} else if !exists {
		cg.Logf("Consumergroup `%s` does not yet exists, creating...\n", cg.group.Name)
		if err := cg.group.Create(); err != nil {
			cg.Logf("FAILED to create consumergroup in Zookeeper: %s!\n", err)
			_ = consumer.Close()
			_ = kz.Close()
			return nil, err
		}
	}

	// Register itself with zookeeper
	if err := cg.instance.Register(topics); err != nil {
		cg.Logf("FAILED to register consumer instance: %s!\n", err)
		return nil, err
	} else {
		cg.Logf("Consumer instance registered (%s).", cg.instance.ID)
	}

	offsetConfig := OffsetManagerConfig{CommitInterval: config.Offsets.CommitInterval}
	cg.offsetManager = NewZookeeperOffsetManager(cg, &offsetConfig)

	go cg.topicListConsumer(topics)

	return
}

// Returns a channel that you can read to obtain events from Kafka to process.
func (cg *ConsumerGroup) Messages() <-chan *sarama.ConsumerMessage {
	return cg.messages
}

// Returns a channel that you can read to obtain events from Kafka to process.
func (cg *ConsumerGroup) Errors() <-chan error {
	return cg.errors
}

func (cg *ConsumerGroup) Closed() bool {
	return cg.instance == nil
}

func (cg *ConsumerGroup) Close() error {
	shutdownError := AlreadyClosing
	cg.singleShutdown.Do(func() {
		defer cg.kazoo.Close()

		shutdownError = nil

		close(cg.stopper)
		cg.wg.Wait()

		if err := cg.offsetManager.Close(); err != nil {
			cg.Logf("FAILED closing the offset manager: %s!\n", err)
		}

		if shutdownError = cg.instance.Deregister(); shutdownError != nil {
			cg.Logf("FAILED deregistering consumer instance: %s!\n", shutdownError)
		} else {
			cg.Logf("Deregistered consumer instance %s.\n", cg.instance.ID)
		}

		if shutdownError = cg.consumer.Close(); shutdownError != nil {
			cg.Logf("FAILED closing the Sarama client: %s\n", shutdownError)
		}

		close(cg.messages)
		close(cg.errors)
		cg.instance = nil
	})

	return shutdownError
}

func (cg *ConsumerGroup) Logf(format string, args ...interface{}) {
	var identifier string
	if cg.instance == nil {
		identifier = "(defunct)"
	} else {
		identifier = cg.instance.ID[len(cg.instance.ID)-12:]
	}
	sarama.Logger.Printf("[%s/%s] %s", cg.group.Name, identifier, fmt.Sprintf(format, args...))
}

func (cg *ConsumerGroup) InstanceRegistered() (bool, error) {
	return cg.instance.Registered()
}

func (cg *ConsumerGroup) CommitUpto(message *sarama.ConsumerMessage) error {
	cg.offsetManager.MarkAsProcessed(message.Topic, message.Partition, message.Offset)
	return nil
}

func (cg *ConsumerGroup) FlushOffsets() error {
	return cg.offsetManager.Flush()
}

func (cg *ConsumerGroup) topicListConsumer(topics []string) {
	for {
		select {
		case <-cg.stopper:
			return
		default:
		}

		consumers, consumerChanges, err := cg.group.WatchInstances()
		if err != nil {
			cg.Logf("FAILED to get list of registered consumer instances: %s\n", err)
			return
		}

		cg.consumers = consumers
		cg.Logf("Currently registered consumers: %d\n", len(cg.consumers))

		stopper := make(chan struct{})

		for _, topic := range topics {
			cg.wg.Add(1)
			go cg.topicConsumer(topic, cg.messages, cg.errors, stopper)
		}

		select {
		case <-cg.stopper:
			close(stopper)
			return

		case <-consumerChanges:
			registered, err := cg.instance.Registered()
			if err != nil {
				cg.Logf("FAILED to get register status: %s\n", err)
			} else if !registered {
				err = cg.instance.Register(topics)
				if err != nil {
					cg.Logf("FAILED to register consumer instance: %s!\n", err)
				} else {
					cg.Logf("Consumer instance registered (%s).", cg.instance.ID)
				}
			}

			cg.Logf("Triggering rebalance due to consumer list change\n")
			close(stopper)
			cg.wg.Wait()
		}
	}
}

func (cg *ConsumerGroup) topicConsumer(topic string, messages chan<- *sarama.ConsumerMessage, errors chan<- error, stopper <-chan struct{}) {
	defer cg.wg.Done()

	select {
	case <-stopper:
		return
	default:
	}

	cg.Logf("%s :: Started topic consumer\n", topic)

	// Fetch a list of partition IDs
	partitions, err := cg.kazoo.Topic(topic).Partitions()
	if err != nil {
		cg.Logf("%s :: FAILED to get list of partitions: %s\n", topic, err)
		cg.errors <- &sarama.ConsumerError{
			Topic:     topic,
			Partition: -1,
			Err:       err,
		}
		return
	}

	partitionLeaders, err := retrievePartitionLeaders(partitions)
	if err != nil {
		cg.Logf("%s :: FAILED to get leaders of partitions: %s\n", topic, err)
		cg.errors <- &sarama.ConsumerError{
			Topic:     topic,
			Partition: -1,
			Err:       err,
		}
		return
	}

	dividedPartitions := dividePartitionsBetweenConsumers(cg.consumers, partitionLeaders)
	myPartitions := dividedPartitions[cg.instance.ID]
	cg.Logf("%s :: Claiming %d of %d partitions", topic, len(myPartitions), len(partitionLeaders))

	// Consume all the assigned partitions
	var wg sync.WaitGroup
	for _, pid := range myPartitions {

		wg.Add(1)
		go cg.partitionConsumer(topic, pid.ID, messages, errors, &wg, stopper)
	}

	wg.Wait()
	cg.Logf("%s :: Stopped topic consumer\n", topic)
}

func (cg *ConsumerGroup) consumePartition(topic string, partition int32, nextOffset int64) (sarama.PartitionConsumer, error) {
	consumer, err := cg.consumer.ConsumePartition(topic, partition, nextOffset)
	if err == sarama.ErrOffsetOutOfRange {
		cg.Logf("%s/%d :: Partition consumer offset out of Range.\n", topic, partition)
		// if the offset is out of range, simplistically decide whether to use OffsetNewest or OffsetOldest
		// if the configuration specified offsetOldest, then switch to the oldest available offset, else
		// switch to the newest available offset.
		if cg.config.Offsets.Initial == sarama.OffsetOldest {
			nextOffset = sarama.OffsetOldest
			cg.Logf("%s/%d :: Partition consumer offset reset to oldest available offset.\n", topic, partition)
		} else {
			nextOffset = sarama.OffsetNewest
			cg.Logf("%s/%d :: Partition consumer offset reset to newest available offset.\n", topic, partition)
		}
		// retry the consumePartition with the adjusted offset
		consumer, err = cg.consumer.ConsumePartition(topic, partition, nextOffset)
	}
	if err != nil {
		cg.Logf("%s/%d :: FAILED to start partition consumer: %s\n", topic, partition, err)
		return nil, err
	}
	return consumer, err
}

// Consumes a partition
func (cg *ConsumerGroup) partitionConsumer(topic string, partition int32, messages chan<- *sarama.ConsumerMessage, errors chan<- error, wg *sync.WaitGroup, stopper <-chan struct{}) {
	defer wg.Done()

	select {
	case <-stopper:
		return
	default:
	}

	// Since ProcessingTimeout is the amount of time we'll wait for the final batch
	// of messages to be processed before releasing a partition, we need to wait slightly
	// longer than that before timing out here to ensure that another consumer has had
	// enough time to release the partition. Hence, +2 seconds.
	maxRetries := int(cg.config.Offsets.ProcessingTimeout/time.Second) + 2
	for tries := 0; tries < maxRetries; tries++ {
		if err := cg.instance.ClaimPartition(topic, partition); err == nil {
			break
		} else if tries+1 < maxRetries {
			if err == kazoo.ErrPartitionClaimedByOther {
				// Another consumer still owns this partition. We should wait longer for it to release it.
				time.Sleep(1 * time.Second)
			} else {
				// An unexpected error occurred. Log it and continue trying until we hit the timeout.
				cg.Logf("%s/%d :: FAILED to claim partition on attempt %v of %v; retrying in 1 second. Error: %v", topic, partition, tries+1, maxRetries, err)
				time.Sleep(1 * time.Second)
			}
		} else {
			cg.Logf("%s/%d :: FAILED to claim the partition: %s\n", topic, partition, err)
			cg.errors <- &sarama.ConsumerError{
				Topic:     topic,
				Partition: partition,
				Err:       err,
			}
			return
		}
	}

	defer func() {
		err := cg.instance.ReleasePartition(topic, partition)
		if err != nil {
			cg.Logf("%s/%d :: FAILED to release partition: %s\n", topic, partition, err)
			cg.errors <- &sarama.ConsumerError{
				Topic:     topic,
				Partition: partition,
				Err:       err,
			}
		}
	}()

	nextOffset, err := cg.offsetManager.InitializePartition(topic, partition)
	if err != nil {
		cg.Logf("%s/%d :: FAILED to determine initial offset: %s\n", topic, partition, err)
		return
	}

	if nextOffset >= 0 {
		cg.Logf("%s/%d :: Partition consumer starting at offset %d.\n", topic, partition, nextOffset)
	} else {
		nextOffset = cg.config.Offsets.Initial
		if nextOffset == sarama.OffsetOldest {
			cg.Logf("%s/%d :: Partition consumer starting at the oldest available offset.\n", topic, partition)
		} else if nextOffset == sarama.OffsetNewest {
			cg.Logf("%s/%d :: Partition consumer listening for new messages only.\n", topic, partition)
		}
	}

	consumer, err := cg.consumePartition(topic, partition, nextOffset)

	if err != nil {
		cg.Logf("%s/%d :: FAILED to start partition consumer: %s\n", topic, partition, err)
		return
	}

	defer consumer.Close()

	err = nil
	var lastOffset int64 = -1 // aka unknown
partitionConsumerLoop:
	for {
		select {
		case <-stopper:
			break partitionConsumerLoop

		case err := <-consumer.Errors():
			if err == nil {
				cg.Logf("%s/%d :: Consumer encountered an invalid state: re-establishing consumption of partition.\n", topic, partition)

				// Errors encountered (if any) are logged in the consumerPartition function
				var cErr error
				consumer, cErr = cg.consumePartition(topic, partition, lastOffset)
				if cErr != nil {
					break partitionConsumerLoop
				}
				continue partitionConsumerLoop
			}

			for {
				select {
				case errors <- err:
					continue partitionConsumerLoop

				case <-stopper:
					break partitionConsumerLoop
				}
			}

		case message := <-consumer.Messages():
			if message == nil {
				cg.Logf("%s/%d :: Consumer encountered an invalid state: re-establishing consumption of partition.\n", topic, partition)

				// Errors encountered (if any) are logged in the consumerPartition function
				var cErr error
				consumer, cErr = cg.consumePartition(topic, partition, lastOffset)
				if cErr != nil {
					break partitionConsumerLoop
				}
				continue partitionConsumerLoop

			}

			for {
				select {
				case <-stopper:
					break partitionConsumerLoop

				case messages <- message:
					lastOffset = message.Offset
					continue partitionConsumerLoop
				}
			}
		}
	}

	cg.Logf("%s/%d :: Stopping partition consumer at offset %d\n", topic, partition, lastOffset)
	if err := cg.offsetManager.FinalizePartition(topic, partition, lastOffset, cg.config.Offsets.ProcessingTimeout); err != nil {
		cg.Logf("%s/%d :: %s\n", topic, partition, err)
	}
}
