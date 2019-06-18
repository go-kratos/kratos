# database/hbase

## 说明
Hbase Client，进行封装加入了链路追踪和统计。

## 配置
需要指定hbase集群的zookeeper地址。
```
config := &hbase.Config{Zookeeper: &hbase.ZKConfig{Addrs: []string{"localhost"}}}
client := hbase.NewClient(config)
```

## 使用方式
```
package main

import (
	"context"
	"fmt"
    
    "github.com/bilibili/kratos/pkg/database/hbase"
)

func main() {
    config := &hbase.Config{Zookeeper: &hbase.ZKConfig{Addrs: []string{"localhost"}}}
    client := hbase.NewClient(config)

    // 
    values := map[string]map[string][]byte{"name": {"firstname": []byte("hello"), "lastname": []byte("world")}}
    ctx := context.Background()

    // 写入信息
    // table: user
    // rowkey: user1
    // values["family"] = columns
    _, err := client.PutStr(ctx, "user", "user1", values)
    if err != nil {
        panic(err)
    }

    // 读取信息
    // table: user
    // rowkey: user1
    result, err := client.GetStr(ctx, "user", "user1")
    if err != nil {
        panic(err)
    }
    fmt.Printf("%v", result)
}
```

-------------

[文档目录树](summary.md)
