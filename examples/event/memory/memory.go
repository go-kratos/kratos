package memory

import (
	"context"
	"log"
	"sync"

	"github.com/SeeMusic/kratos/examples/event/event"
)

var (
	_ event.Sender   = (*memorySender)(nil)
	_ event.Receiver = (*memoryReceiver)(nil)
	_ event.Event    = (*Message)(nil)
)

var (
	chanMap = struct {
		sync.RWMutex
		cm map[string]chan *Message
	}{}
	ChanSize = 256
)

func init() {
	chanMap.cm = make(map[string]chan *Message)
}

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

type memorySender struct {
	topic string
}

func (m *memorySender) Send(ctx context.Context, msg event.Event) error {
	chanMap.cm[m.topic] <- &Message{
		key:   msg.Key(),
		value: msg.Value(),
	}
	return nil
}

func (m *memorySender) Close() error {
	return nil
}

type memoryReceiver struct {
	topic string
}

func (m *memoryReceiver) Receive(ctx context.Context, handler event.Handler) error {
	go func() {
		for msg := range chanMap.cm[m.topic] {
			err := handler(context.Background(), msg)
			if err != nil {
				log.Fatal("message handling exception:", err)
			}
		}
	}()
	return nil
}

func (m *memoryReceiver) Close() error {
	return nil
}

func NewMemory(topic string) (event.Sender, event.Receiver) {
	chanMap.RLock()
	if _, ok := chanMap.cm[topic]; !ok {
		// chanMap.Lock()
		chanMap.cm[topic] = make(chan *Message, ChanSize)
		// chanMap.Unlock()
	}
	defer chanMap.RUnlock()

	return &memorySender{topic: topic}, &memoryReceiver{topic: topic}
}
