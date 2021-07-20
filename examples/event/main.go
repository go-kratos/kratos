package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-kratos/kratos/examples/event/event"
	"github.com/go-kratos/kratos/examples/event/kafka"
	"github.com/go-kratos/kratos/v2/metadata"
)

func main() {
	client, err := kafka.NewKafkaClient([]string{"39.106.218.150:9092"})
	if err != nil {
		panic(err)
	}
	sender, err := kafka.NewKafkaSender(client, "test")
	if err != nil {
		panic(err)
	}
	receiver, err := kafka.NewKafkaReceiver(client, "test")
	if err != nil {
		panic(err)
	}
	receive(receiver)
	for i := 0; i < 5; i++ {
		send(sender, strconv.Itoa(i))
	}
	_ = sender.Close()
	_ = receiver.Close()
	_ = client.Close()
}

func send(sender event.Sender, num string) {
	msg := kafka.NewMessage("send", []byte("hello world"), map[string]string{
		"x-md-global-service": "service-" + num,
	})
	err := sender.Send(context.Background(), msg)
	if err != nil {
		panic(err)
	}
}

func receive(receiver event.Receiver) {
	err := receiver.Receive(context.Background(), func(ctx context.Context, message event.Message) error {
		fmt.Printf("key:%s, value:%s, header:%s\n", message.Key(), message.Value(), message.Header())
		doSomething(ctx)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func doSomething(ctx context.Context) {
	// metadata
	if md, ok := metadata.FromServerContext(ctx); ok {
		fmt.Println(md)
	}
}
