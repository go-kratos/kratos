package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/SeeMusic/kratos/examples/event/event"
	"github.com/SeeMusic/kratos/examples/event/kafka"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	receiver, err := kafka.NewKafkaReceiver([]string{"localhost:9092"}, "kratos")
	if err != nil {
		panic(err)
	}
	receive(receiver)

	<-sigs
	_ = receiver.Close()
}

func receive(receiver event.Receiver) {
	fmt.Println("start receiver")
	err := receiver.Receive(context.Background(), func(ctx context.Context, msg event.Event) error {
		fmt.Printf("key:%s, value:%s\n", msg.Key(), msg.Value())
		return nil
	})
	if err != nil {
		return
	}
}
