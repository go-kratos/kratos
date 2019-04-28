package service

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/interface/main/player/conf"
)

// test func Player
func BenchmarkPolicy(b *testing.B) {
	if err := conf.Init(); err != nil {
		fmt.Println(err)
	}
	ser := New(conf.Conf)
	c := context.Background()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ser.Policy(c, 1, 6698028)
		}
	})
}
