### database/hbase

### 项目简介

Hbase Client，进行封装加入了链路追踪和统计。

### usage
```go
package main

import (
	"context"
	"fmt"

	"github.com/bilibili/kratos/pkg/database/hbase"
)

func main() {
	config := &hbase.Config{Zookeeper: &hbase.ZKConfig{Addrs: []string{"localhost"}}}
	client := hbase.NewClient(config)

	values := map[string]map[string][]byte{"name": {"firstname": []byte("hello"), "lastname": []byte("world")}}
	ctx := context.Background()

	_, err := client.PutStr(ctx, "user", "user1", values)
	if err != nil {
		panic(err)
	}

	result, err := client.GetStr(ctx, "user", "user1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", result)
}
```

##### 依赖包

1.[gohbase](https://github.com/tsuna/gohbase)
