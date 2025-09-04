# Logger

## Usage

### Structured logging

```go
logger := log.NewStdLogger(os.Stdout)
// fields & valuer
logger = log.With(logger,
    "service.name", "helloworld",
    "service.version", "v1.0.0",
    "ts", log.DefaultTimestamp,
    "caller", log.DefaultCaller,
)
logger.Log(log.LevelInfo, "key", "value")

// helper
helper := log.NewHelper(logger)
helper.Log(log.LevelInfo, "key", "value")
helper.Info("info message")
helper.Infof("info %s", "message")
helper.Infow("key", "value")

// filter
log := log.NewHelper(log.NewFilter(logger,
	log.FilterLevel(log.LevelInfo),
	log.FilterKey("foo"),
	log.FilterValue("bar"),
	log.FilterFunc(customFilter),
))
log.Debug("debug log")
log.Info("info log")
log.Warn("warn log")
log.Error("warn log")
```

## Third party log library

### zap

```shell
go get -u github.com/go-kratos/kratos/contrib/log/zap/v2
```
### logrus

```shell
go get -u github.com/go-kratos/kratos/contrib/log/logrus/v2
```

### fluent

```shell
go get -u github.com/go-kratos/kratos/contrib/log/fluent/v2
```

### aliyun

```shell
go get -u github.com/go-kratos/kratos/contrib/log/aliyun/v2
```