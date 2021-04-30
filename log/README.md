# Log

## Usage

### Structured logging

```go
logger := log.NewLogger(os.Stdout)
logger = log.With(logger, "caller", log.DefaultCaller, "ts", log.DefaultTimestamp)

// Levels
log.Debug(logger).Print("msg", "foo bar")
log.Info(logger).Print("msg", "foo bar")
log.Warn(logger).Print("msg", "foo bar")
log.Error(logger).Print("msg", "foo bar")

errLogger := log.Error(logger)
errLogger.Print("msg", "xxx")
errLogger.Print("msg", "yyy")
errLogger.Print("msg", "zzz")

errLogger.Print(
    "http.scheme", "https",
    "http.host", "translate.googleapis.com",
    "http.target", "/language/translate",
    "http.method", "post",
    "http.status_code", 500,
    "http.flavor", "1.1.",
    "http.user_agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
)
```

