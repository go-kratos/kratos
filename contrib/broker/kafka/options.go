package kafka

import (
	"time"

	"github.com/segmentio/kafka-go/compress"
)

type Option func(*kafkaBroker)

func WithOffset(offset int64) Option {
	return func(o *kafkaBroker) {
		o.offset = offset
	}
}

func WithMinBytes(minBytes int) Option {
	return func(o *kafkaBroker) {
		o.minBytes = minBytes
	}
}

func WithMaxBytes(maxBytes int) Option {
	return func(o *kafkaBroker) {
		o.maxBytes = maxBytes
	}
}

func WithMaxWait(maxWait int64) Option {
	return func(o *kafkaBroker) {
		o.maxWait = time.Millisecond * time.Duration(maxWait)
	}
}

func WithCommitInterval(commitInterval int64) Option {
	return func(o *kafkaBroker) {
		o.commitInterval = time.Millisecond * time.Duration(commitInterval)
	}
}

func WithQueueCapacity(queueCapacity int) Option {
	return func(o *kafkaBroker) {
		o.queueCapacity = queueCapacity
	}
}

func WithCompressionCodec(compressionCodec compress.Codec) Option {
	return func(p *kafkaBroker) {
		p.compressionCodec = compressionCodec
	}
}

func WithCommitIgnoreError(commitIgnoreError bool) Option {
	return func(o *kafkaBroker) {
		o.commitIgnoreError = commitIgnoreError
	}
}

func WithConsumes(consumers int) Option {
	return func(o *kafkaBroker) {
		if consumers <= 0 {
			consumers = 1
		}
		o.consumers = consumers
	}
}

func WithProcessors(processors int) Option {
	return func(o *kafkaBroker) {
		if processors <= 0 {
			processors = 1
		}
		o.processors = processors
	}
}
