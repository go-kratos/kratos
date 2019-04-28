package archive

import (
	"fmt"
	"testing"
	"time"
)

func Test_VideoEditor(t *testing.T) {
	edit := NewEditor(3)

	edit.Add(10,
		func() { fmt.Println("WOW!10 SUCCESS") },
		func() { fmt.Println("WOW!10 FAIL") },
		time.Duration(20*time.Second),
		[]func() (int64, int, int, error){
			func() (int64, int, int, error) { return 10, 0, 1, nil },
			func() (int64, int, int, error) { return 10, 2, 3, fmt.Errorf("boom") },
		}...)

	edit.Add(11,
		func() { fmt.Println("WOW!11 SUCCESS") },
		func() { fmt.Println("WOW!11 FAIL") },
		time.Duration(5*time.Second),
		[]func() (int64, int, int, error){
			func() (int64, int, int, error) { return 11, 0, 1, nil },
			func() (int64, int, int, error) { return 11, 0, 1, nil },
		}...)

	time.Sleep(time.Second)
	edit.Close()
	t.Fail()
}
