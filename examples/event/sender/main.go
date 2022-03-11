package main

import (
	"context"
	"fmt"

	"github.com/SeeMusic/kratos/examples/event/event"
	"github.com/SeeMusic/kratos/examples/event/kafka"
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
	msg := kafka.NewMessage("kratos", []byte("hello world"))
	err := sender.Send(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("key:%s, value:%s\n", msg.Key(), msg.Value())
}
