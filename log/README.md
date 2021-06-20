# Log

## Usage

### Structured logging

```go
// log
logger := log.MultiLogger(log.NewStdLogger(os.Stdout), log.DefaultLogger)
logger.Log(log.LevelInfo, "key", "value")

// helper
helper := log.NewHelper(log.DefaultLogger)
helper.Log(log.LevelInfo, "key", "value")
helper.Info("info message")
helper.Infof("info %s", "message")
helper.Infow("key", "value")

// filter
logger := log.NewHelper(log.NewFilter(log.DefaultLogger,
	log.FilterLevel(LevelDebug),
	log.FilterKey("username"),
	log.FilterValue("hello"),
	log.FilterFunc(testFilterFunc),
))
logger.Log(LevelDebug, "msg", "test debug")
logger.Debug("test debug")
logger.Debugf("test %s", "debug")
logger.Debugw("log", "test debug")
logger.Warn("warn log")

// valuer
logger := log.With(log.DefaultLogger, "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)
logger.Log(log.LevelInfo, "msg", "helloworld")
```
