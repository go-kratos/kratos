# config

可以指定多个配置源，config 会进行合并成 map[string]interface{}，然后通过 Scan 或者 Value 获取值内容；

```
c := config.New(
    config.WithSource(
        file.NewSource(path),
    ),
    config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
        // kv.Key
        // kv.Value
        // kv.Metadata
        // 自定义实现对应的数据源解析，如果是配置中心数据源也可以指定metadata进行识别配置类型
        return yaml.Unmarshal(kv.Value, v)
    }),
)
// 加载配置源：
if err := c.Load(); err != nil {
    panic(err)
}
// 获取对应的值内容：
name, err := c.Value("service").String()
// 解析到结构体（由于已经合并到map[string]interface{}，所以需要指定 jsonName 进行解析）：
var v struct {
    Service string `json:"service"`
    Version string `json:"version"`
}
if err := c.Scan(&v); err != nil {
    panic(err)
}
// 监听值内容变更
c.Watch("service.name", func(key string, value config.Value) {
    // 值内容变更
})
```

