package health

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"golang.org/x/net/context"
)

type A struct{}

func (A) Check(ctx context.Context) (interface{}, error) {
	fmt.Println("check A")
	if rand.Int()%5 == 0 {
		return "出错A", fmt.Errorf("错误:%s", "123")
	}
	return "正常A", nil
}

type B struct{}

func (B) Check(ctx context.Context) (interface{}, error) {
	fmt.Println("check B")
	if rand.Int()%5 == 0 {
		return "出错B", fmt.Errorf("错误:%s", "123B")
	}
	return "正常B", nil
}

func TestNew(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	cm := New(ctx)
	cm.RegisterChecker(NewChecker("A", A{}, WithInterval(time.Second), WithTimeout(time.Second)))
	cm.RegisterChecker(NewChecker("B", B{}, WithInterval(time.Second), WithTimeout(time.Second)))
	cm.Start()
	go func() {
		w := cm.NewWatcher()
		s := cm.GetStatus()
		fmt.Println("----", s)
		defer w.Close()
		for i := range w.Ch {
			fmt.Println("--->>", i, cm.GetStatus(i))
		}
	}()
	time.Sleep(time.Second * 20)
	cm.Stop()
	t.Log("-----=")
	time.Sleep(time.Second * 5)
	cancel()
}
