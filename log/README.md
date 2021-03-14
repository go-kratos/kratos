# Log

## Usage

### Structured logging

```
logger := log.NewLogger(os.Stdout)
logger = With(logger, "key", "value")

log := log.NewHelper("github.com/project/foo", logger)
// Levels
log.Info("hello")
log.Infof("hello %s", "kratos")
log.Infow("key", "value")
```

