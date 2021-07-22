package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kratos/kratos/examples/event/event"
	"github.com/go-kratos/kratos/examples/event/kafka"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	receiver, err := kafka.NewKafkaReceiver([]string{"localhost:9092"}, "kratos")
	if err != nil {
		panic(err)
	}
	receive(receiver)
	select {
	case <-sigs:
		_ = receiver.Close()
	}
}

func receive(receiver event.Receiver) {
	fmt.Println("start receiver")
	err := receiver.Receive(context.Background(), func(ctx context.Context, message event.Message) error {
		fmt.Printf("key:%s, value:%s, header:%s\n", message.Key(), message.Value(), message.Header())
		return nil
	})
	if err != nil {
		return
	}
}
