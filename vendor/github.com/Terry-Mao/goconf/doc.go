/*
Package goconf provides configuraton read and write implementations.

Examples:
    package main

    import (
        "fmt"
        "github.com/Terry-Mao/goconf"
		"time"
    )

    type TestConfig struct {
        ID     int      `goconf:"core:id"`
        Col    string   `goconf:"core:col"`
        Ignore int      `goconf:"-"`
        Arr    []string `goconf:"core:arr:,"`
        Test   time.Duration `goconf:"core:t_1:time"`
        Buf    int           `goconf:"core:buf:memory"`
        Arr1   []int         `goconf:"core:arr1:,"`
        M      map[int]string`goconf:"core:m:,"`
    }

    func main() {
        conf := goconf.New()
        if err := conf.Parse("./examples/conf_test.txt"); err != nil {
            panic(err)
        }
        core := conf.Get("core")
        if core == nil {
            panic("no core section")
        }
        id, err := core.Int("id")
        if err != nil {
            panic(err)
        }
        fmt.Println(id)
        tf := &TestConfig{}
        if err := conf.Unmarshal(tf); err != nil {
            panic(err)
        }
        fmt.Println(tf.ID)
    }
*/
package goconf
