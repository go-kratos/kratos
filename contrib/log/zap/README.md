# Usage

- Copy `third_party/contrib/log` folder to your project under `third_party/contrib`
- Import log proto to your `conf.proto`, and generate code
```proto
import "contrib/log/zap/config.proto";
import "contrib/log/lumberjack/config.proto";

message Bootstrap {
  ...
  Logger logger = 3;
}

message Logger {
  zap.Config zap = 1;
  lumberjack.Config lumberjack = 2;
}
```
- Init zap logger with file rotation in your `main.go`
```go
package main

import (
	"flag"
	"net/url"
	"os"

	"github.com/go-kratos/kratos/contrib/log/lumberjack/v2"
	"github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "go.uber.org/automaxprocs"

	"<your project package>/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "c", "../../configs", "config path, eg: -c config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	if bc.Logger.Lumberjack != nil {
		err := zap.RegisterSink("lumberjack", func(u *url.URL) (zap.Sink, error) {
			return lumberjack.NewLoggerWithURL(bc.Logger.Lumberjack, u), nil
		})
		if err != nil {
			panic(err)
		}
	}

	logger, err := zap.NewZapLogger(bc.Logger.Zap)
	if err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

```
- Add zap log config to your `config.yaml`, then you run your app with `-c config.yaml`
```yaml
server:
  http:
    addr: 0.0.0.0:8080
    timeout: 30s
  grpc:
    addr: 0.0.0.0:9000
    timeout: 30s

...

logger:
  zap:
    level: debug # Optional, default: info
    development: true # Optional, default: false
    disableCaller: true # Optional, default: false
    disableStacktrace: true # Optional, default: false
    encoding: json # Optional, default: json
    outputPaths: ["lumberjack:", "stdout"] # Optional, default: [stdout]
    errorOutputPaths: ["lumberjack:", "stderr"] # Optional, default: [stderr]
    # encoderConfig: # Optional
    #   messageKey: "msg" # Optional, default: msg
    #   timeKey: "ts" # Optional, default: ts
    #   levelKey: "level" # Optional, default: level
    #   nameKey: "logger" # Optional, default: logger
    #   callerKey: "caller" # Optional, default: caller
    #   stacktraceKey: "stacktrace" # Optional, default: stacktrace
    #   skipLineEnding: false # Optional, default: false
    #   lineEnding: "\n" # Optional, default: \n
    #   consoleSeparator: "\t" # Optional, default: \t
    # initialFields: # Optional, default: empty map
    #   key: value
  lumberjack: # Optional, default: nil
    filename: "logs/app.log" # Optional, use zap path if empty
    maxsize: 1024 # Optional, default: 1024 (MB)
    maxage: 7 # Optional, default: 7 (day)
    maxbackups: 3 # Optional, default: 3 (day)
    localtime: true # Optional, default: false
    compress: true # Optional, default: false
```


