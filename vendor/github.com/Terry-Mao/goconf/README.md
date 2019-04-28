## Terry-Mao/goconf

`Terry-Mao/goconf` is an configuration file parse module.

## Requeriments
* Go 1.2 or higher

## Installation

Just pull `Terry-Mao/goconf` from github using `go get`:

```sh
# download the code
$ go get -u github.com/Terry-Mao/goconf
```

## Usage

```go
package main                                                                   
                                                                               
import (                                                                       
    "fmt"                                                                      
    "github.com/Terry-Mao/goconf"                                              
    "time"
)                                                                              

type TestConfig struct {
	ID     int           `goconf:"core:id"`
	Col    string        `goconf:"core:col"`
	Ignore int           `goconf:"-"`
	Arr    []string      `goconf:"core:arr:,"`
	Test   time.Duration `goconf:"core:t_1:time"`
	Buf    int           `goconf:"core:buf:memory"`
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
```

## Examples

```sh
# configuration examples
# this is comment, goconf will ignore it.

# this is the section name
[core]

# a key-value config which key is test and value is 1
test 1

# one mb
test1 1mb

# one second
test2 1s

# boolean
test3 true

# arr
arr hello,the,world

# map
m 1=hello,2=the,3=world
```

## Documentation

Read the `Terry-Mao/goconf` documentation from a terminal

```go
$ godoc github.com/Terry-Mao/goconf -http=:6060
```

Alternatively, you can [goconf](http://godoc.org/github.com/Terry-Mao/goconf) online.
