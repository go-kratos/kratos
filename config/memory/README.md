# memory

## usage

```gotemplate
var yamlConf = `
server:
  http:
    addr: 0.0.0.0:8000
    timeout: 1s
  grpc:
    addr: 0.0.0.0:9000
    timeout: 1s
  name: hizhc
data:
  database:
    driver: mysql
    source: root:root@tcp(127.0.0.1:3306)/hizhc?charset=utf8mb4&parseTime=True&loc=Local
  redis:
    addr: 127.0.0.1:6379
    dial_timeout: 0.1s
    read_timeout: 0.2s
    write_timeout: 0.2s
`


func main() {
	flag.Parse()
	logger := log.NewStdLogger(os.Stdout)

	cfg := config.New(
		config.WithSource(
			memory.NewSource(memory.WithYAML([]byte(yamlConf))),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)
	if err := cfg.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := cfg.Scan(&bc); err != nil {
		panic(err)
	}

	app, err := initApp(&bc, bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
```