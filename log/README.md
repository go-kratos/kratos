# Log

## Usage

### Structured logging

```
logger := stdlog.NewLogger(stdlog.Writer(os.Stdout))
log := log.NewHelper("module_name", logger)
// Levels
log.Info("some log")
log.Infof("format %s", "some log")
log.Infow("field", "some log")
// Verbose
log.V(10).Print("field", "some log")
```

