# Log

## Usage

### Structured logging

```
logger := log.NewLogger(os.Stdout)
logger = With(logger, "foo", "bar")

log := log.NewHelper("github.com/project/foo", logger)
// Levels
log.Info("hello")
log.Infof("hello %s", "go")
log.Infow("key", "value")
```

