# Polaris Config

```go
import (
 "log"

 "github.com/polarismesh/polaris-go"
 "github.com/go-kratos/kratos/contrib/config/polaris/v2"
)

func main() {

 configApi, err := polaris.NewConfigAPI()
 if err != nil {
  log.Fatalln(err)
 }

 source, err := New(&configApi, WithNamespace("default"), WithFileGroup("default"), WithFileName("default.yaml"))

 if err != nil {
  log.Fatalln(err)
 }
 source.Load()
}
```
