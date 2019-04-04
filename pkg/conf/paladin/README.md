#### paladin

##### 项目简介

paladin 是一个config SDK客户端，包括了file、mock几个抽象功能，方便使用本地文件或者sven配置中心，并且集成了对象自动reload功能。  


local files:
```
demo -conf=/data/conf/app/msm-servie.toml
// or dir
demo -conf=/data/conf/app/

```
example:
```
type exampleConf struct {
	Bool   bool
	Int    int64
	Float  float64
	String string
}

func (e *exampleConf) Set(text string) error {
	var ec exampleConf
	if err := toml.Unmarshal([]byte(text), &ec); err != nil {
		return err
	}
	*e = ec
	return nil
}

func ExampleClient() {
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	var (
		ec   exampleConf
		eo   exampleConf
		m    paladin.TOML
		strs []string
	)
	// config unmarshal
	if err := paladin.Get("example.toml").UnmarshalTOML(&ec); err != nil {
		panic(err)
	}
	// config setter
	if err := paladin.Watch("example.toml", &ec); err != nil {
        panic(err)
    }
	// paladin map
	if err := paladin.Watch("example.toml", &m); err != nil {
        panic(err)
    }
	s, err := m.Value("key").String()
	b, err := m.Value("key").Bool()
	i, err := m.Value("key").Int64()
	f, err := m.Value("key").Float64()
	// value slice
	err = m.Value("strings").Slice(&strs)
	// watch key
	for event := range paladin.WatchEvent(context.TODO(), "key") {
		fmt.Println(event)
	}
}
```

##### 编译环境

- **请只用 Golang v1.12.x 以上版本编译执行**

##### 依赖包
