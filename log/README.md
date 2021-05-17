# Log

## Usage

### Structured logging

```go
Logger logger = log.MultiLogger(log.NewStdLogger(os.Stdout), syslog.NewLogger())

logger = log.With(logger,
    "service.name", "hellworld",
    "service.version", "v1.0.0",
    "ts", log.DefaultTimestamp,
    "caller", log.DefaultCaller,
)

logger.Log(log.LevelInfo, "key", "value")

Helper helper = log.NewHelper(logger)
helper.Log(log.LevelInfo, "key", "value")
helper.Info("info message")
helper.Infof("info %s", "message")
helper.Infow("key", "value")
```
