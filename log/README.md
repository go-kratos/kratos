# Logger

## Usage

### Structured logging

```go
logger := log.NewStdLogger(os.Stdout)
// fields & valuer
logger = log.With(logger,
    "service.name", "hellworld",
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
logger := log.NewHelper(log.NewFilter(logger,
	log.FilterLevel(LevelInfo),
	log.FilterKey("foo"),
	log.FilterValue("bar"),
	log.FilterFunc(customFilter),
))
logger.Log(log.LevelDebug, "foo", "bar")
logger.Debug("debug log")
logger.Info("info log")
logger.Warn("warn log")
logger.Error("warn log")
```
