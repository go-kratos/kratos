# Log

## Usage

### Structured logging

```
log := NewHelper("module_name", stdlog.NewLogger(stdlog.Writer(os.Stdout)))
// Levels
log.Info("some log")
log.Infof("format %s", "some log")
log.Infow("field", "some log")
// Verbose
log.V(10).Print("field", "some log")
```

