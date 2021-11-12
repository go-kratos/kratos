## Apollo config centry

This module implements the `config.Source` interface in kratos based apollo config management center.

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/go-kratos/kratos/contrib/config/apollo/v2)

### Quick start

```go
import (
	"fmt"
	"log"

	"github.com/go-kratos/kratos/contrib/config/apollo/v2"
	"github.com/go-kratos/kratos/v2/config"
)

func main() {
	c := config.New(
		config.WithSource(
			apollo.NewSource(
				apollo.WithAppID("kratos"),
				apollo.WithCluster("dev"),
				apollo.WithEndpoint("http://localhost:8080"),
				apollo.WithNamespace("application,event.yaml,demo.json"),
				apollo.WithEnableBackup(),
				apollo.WithSecret("ad75b33c77ae4b9c9626d969c44f41ee"),
			),
		),
	)
	var bc bootstrap
	if err := c.Load(); err != nil {
		panic(err)
	}
	
	// use value and watch operations，help yourself. 
}
```

### Options list

> You get what you see.

```go
// specify the app id
func WithAppID(appID string) Option
// specify the cluster of application
func WithCluster(cluster string) Option

// enable backup or not, and where to back up them.
func WithBackupPath(backupPath string) Option
func WithDisableBackup() Option
func WithEnableBackup() Option

// specify apollo endpoint, such as http://localhost:8080
func WithEndpoint(endpoint string) Option

// inject a logger to debug
func WithLogger(logger log.Logger) Option

// namespaces to load, comma to separate. 
func WithNamespace(name string) Option

// secret is the apollo secret key to access application config.
func WithSecret(secret string) Option
```

### Notice

apollo config center use `Namespace` to be part of the key. For example:

***application.json***

```json
{
  "http": {
    "address": ":8080",
    "tls": {
      "enable": false,
      "cert_file": "",
      "key_file": ""
    }
  }
}
```

you got them in kratos config instance maybe look like:

```go
config := map[string]interface{}{
	// application be part of the key path.
	"application": map[string]interface{}{
        "http": map[string]interface{}{
            "address": ":8080",
            "tls": map[string]interface{}{
                "enable": false,
                "cert_file": "",
                "key_file": ""
            }
        }
    }
}
_ = config
```