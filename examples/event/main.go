package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/examples/event/event"
	"github.com/go-kratos/kratos/examples/event/kafka"
)

func main() {
	client, err := kafka.NewKafkaClient([]string{"localhost:9092"})
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
		send(sender)
	}
	_ = sender.Close()
	_ = receiver.Close()
	_ = client.Close()
}

func send(sender event.Sender) {
	msg := kafka.NewMessage("send", []byte("hello world"), map[string]string{
		"user":  "kratos",
		"phone": "123456",
	})
	err := sender.Send(context.Background(), msg)
	if err != nil {
		panic(err)
	}
}

func receive(receiver event.Receiver) {
	err := receiver.Receive(context.Background(), func(ctx context.Context, message event.Message) error {
		fmt.Printf("key:%s, value:%s, header:%s\n", message.Key(), message.Value(), message.Header())
		return nil
	})
	if err != nil {
		panic(err)
	}
}
