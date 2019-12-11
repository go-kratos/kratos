#### paladin

##### 项目简介

paladin 是一个config SDK客户端，包括了file、mock几个抽象功能，方便使用本地文件或者sven\apollo配置中心，并且集成了对象自动reload功能。  

local files:
```
demo -conf=/data/conf/app/msm-servie.toml
// or dir
demo -conf=/data/conf/app/
```

*注：使用远程配置中心的用户在执行应用，如这里的`demo`时务必**不要**带上`-conf`参数，具体见下文远程配置中心的例子*

local file example:
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

remote config center example:
```
type exampleConf struct {
	Bool   bool
	Int    int64
	Float  float64
	String string
}

func (e *exampleConf) Set(text string) error {
	var ec exampleConf
	if err := yaml.Unmarshal([]byte(text), &ec); err != nil {
		return err
	}
	*e = ec
	return nil
}

func ExampleApolloClient() {
	/*
		pass flags or set envs that apollo needs, for example:

		```
		export APOLLO_APP_ID=SampleApp
		export APOLLO_CLUSTER=default
		export APOLLO_CACHE_DIR=/tmp
		export APOLLO_META_ADDR=localhost:8080
		export APOLLO_NAMESPACES=example.yml
		```
	*/

	if err := paladin.Init(apollo.PaladinDriverApollo); err != nil {
		panic(err)
	}
	var (
		ec   exampleConf
		eo   exampleConf
		m    paladin.Map
		strs []string
	)
	// config unmarshal
	if err := paladin.Get("example.yml").UnmarshalYAML(&ec); err != nil {
		panic(err)
	}
	// config setter
	if err := paladin.Watch("example.yml", &ec); err != nil {
        panic(err)
    }
	// paladin map
	if err := paladin.Watch("example.yml", &m); err != nil {
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
