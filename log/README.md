# Log

## Usage

### Structured logging

```
log := NewHelper("module_name", stdlog.NewLogger(Writer(os.Stdout)))
// Levels
log.Debug("some log")
log.Debugf("format %s", "some log")
log.Debugw("field", "some log")
// Verbose
log.V(10).Print("field", "some log")
```

