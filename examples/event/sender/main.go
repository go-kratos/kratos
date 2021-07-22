package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/examples/event/event"
	"github.com/go-kratos/kratos/examples/event/kafka"
)

func main() {
	sender, err := kafka.NewKafkaSender([]string{"localhost:9092"}, "kratos")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 50; i++ {
		send(sender)
	}

	_ = sender.Close()
}

func send(sender event.Sender) {
	msg := kafka.NewMessage("kratos", []byte("hello world"), map[string]string{
		"user":  "kratos",
		"phone": "123456",
	})
	err := sender.Send(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("key:%s, value:%s, header:%s\n", msg.Key(), msg.Value(), msg.Header())
}
