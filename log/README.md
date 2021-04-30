# Log

## Usage

### Structured logging

```go
logger := log.NewStdLogger(os.Stdout)
logger = log.With(logger, "caller", log.DefaultCaller, "ts", log.DefaultTimestamp)

// Levels
log.Debug(logger).Log("msg", "foo bar")
log.Info(logger).Log("msg", "foo bar")
log.Warn(logger).Log("msg", "foo bar")
log.Error(logger).Log("msg", "foo bar")

errLogger := log.Error(logger)
errLogger.Log("msg", "xxx")
errLogger.Log("msg", "yyy")
errLogger.Log("msg", "zzz")

errLogger.Log(
    "http.scheme", "https",
    "http.host", "translate.googleapis.com",
    "http.target", "/language/translate",
    "http.method", "post",
    "http.status_code", 500,
    "http.flavor", "1.1.",
    "http.user_agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
)

// Helper
logger := log.NewHelper("github.com/project/foo", log.DefaultLogger)
logger.Info("hello")
logger.Infof("foo %s", "bar")
logger.Infow("key", "value")

// Verbose
v := NewVerbose(log.DefaultLogger, 20)

v.V(10).Log("foo", "bar1")
v.V(20).Log("foo", "bar2")
v.V(30).Log("foo", "bar3")
```

