# 从 Kratos v2 迁移到 v3

Kratos v3 清理了核心框架中的历史耦合，并将部分隐式依赖改为显式选择。本文总结了 [go-kratos/kratos#3820](https://github.com/go-kratos/kratos/issues/3820) 中已经落地的主要升级点。

## 1. 更新模块路径

将 Kratos 导入路径从 `github.com/go-kratos/kratos/v2` 更新为 `github.com/go-kratos/kratos/v3`。

```go
// v2
import "github.com/go-kratos/kratos/v2"

// v3
import "github.com/go-kratos/kratos/v3"
```

然后刷新依赖：

```shell
go get github.com/go-kratos/kratos/v3@latest
go mod tidy
```

contrib 模块同样使用 `/v3` 导入路径，例如 `github.com/go-kratos/kratos/contrib/middleware/jwt/v3`。

## 2. 显式选择 JSON Codec

v2 中 `encoding/json` 同时处理普通 Go JSON 值和 `proto.Message`，其中 protobuf 消息会使用 protobuf JSON 语义。v3 将核心 JSON codec 拆分：

- `github.com/go-kratos/kratos/v3/encoding/json` 注册标准库 JSON codec，名称为 `json`。
- `github.com/go-kratos/kratos/v3/encoding/protojson` 注册 protobuf JSON codec，名称为 `protojson`。
- `github.com/go-kratos/kratos/contrib/encoding/json/v3` 保留 v2 兼容的 `json` codec 行为，用于迁移期兼容。

新的 v3 代码建议显式导入：

```go
import (
	_ "github.com/go-kratos/kratos/v3/encoding/json"
	_ "github.com/go-kratos/kratos/v3/encoding/protojson"
)
```

如果服务依赖 v2 中 `json` codec 自动处理 protobuf 消息的行为，可以在迁移期使用 contrib 兼容 codec：

```go
import _ "github.com/go-kratos/kratos/contrib/encoding/json/v3"
```

除非明确希望后初始化的 codec 覆盖前一个 `json` 注册，否则不要在同一个进程里同时注册两个 JSON codec。

## 3. 迁移日志到 slog

Kratos v3 使用标准库 `log/slog` API。依赖 `log.Logger`、`log.Helper`、`log.Valuer`、`log.NewStdLogger` 或 trace/service 辅助字段的代码，需要迁移到 `*slog.Logger`、`log.NewHandler` 和 `log.NewLogger`。

```go
import (
	"log/slog"
	"os"

	"github.com/go-kratos/kratos/v3/log"
)

logger := log.NewLogger(
	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}),
	log.WithFilter(log.FilterKey("password")),
).With(
	slog.String("service.name", "helloworld"),
	slog.String("service.version", "v1.0.0"),
)

log.SetDefault(logger)
```

Kratos 应用和 middleware 选项现在接收 `*slog.Logger`。

```go
app := kratos.New(
	kratos.Name("helloworld"),
	kratos.Logger(logger),
)
```

OpenTelemetry 日志能力由 contrib handler 提供：

```go
import (
	"github.com/go-kratos/kratos/v3/log"
	otel "github.com/go-kratos/kratos/contrib/otel/v3/log"
)

logger := log.NewLogger(otel.NewHandler("helloworld"))
```

## 4. 更新 JWT Middleware 导入

JWT middleware 已从核心模块移动到 contrib，只有使用 JWT 的项目才会引入 `github.com/golang-jwt/jwt/v5`。

```go
// v2
import "github.com/go-kratos/kratos/v2/middleware/auth/jwt"

// v3
import "github.com/go-kratos/kratos/contrib/middleware/jwt/v3"
```

修改导入后运行：

```shell
go get github.com/go-kratos/kratos/contrib/middleware/jwt/v3@latest
go mod tidy
```

## 5. 检查自定义熔断器

默认熔断器不再依赖 `github.com/go-kratos/aegis`。使用默认行为的代码无需调整：

```go
handler := circuitbreaker.Client()(next)
```

如果服务注入了自定义 Aegis 熔断器，需要显式添加 Aegis 依赖，并改为通过 `WithBreakerFactory` 注入：

```go
handler := circuitbreaker.Client(
	circuitbreaker.WithBreakerFactory(func() circuitbreaker.CircuitBreaker {
		return newBreaker()
	}),
)(next)
```

## 6. 替换直接使用的 HTTP binding 包

导出的 `transport/http/binding` 包已移除。重新生成后的 `_http.pb.go` 文件不需要手动修改。

手写代码需要按下面方式迁移：

- 将 `binding.EncodeURL` 替换为 `http.BuildPath`。
- 使用 `transport/http.Context` 的 `Bind`、`BindVars`、`BindQuery`、`BindForm`。
- 只有在需要底层 query/form 编解码时才直接使用 `encoding/form`。

```go
path := http.BuildPath("/v1/users/{id}", req)
```

## 7. 重新生成代码

完成依赖和导入更新后，重新生成 Kratos 相关代码：

```shell
go generate ./...
go mod tidy
```

如果服务使用了生成的 HTTP client 或 server，这一步尤其重要，因为 v3 生成代码已经不再导入被移除的 HTTP binding 包。

## 8. 验证升级

上线前运行项目测试和 lint：

```shell
go test ./...
go vet ./...
```

在本仓库中可以使用：

```shell
make test
make lint
```

## 迁移检查清单

- [ ] 将导入路径从 `/v2` 更新为 `/v3`。
- [ ] 选择 `encoding/json`、`encoding/protojson` 或 contrib 兼容 JSON codec。
- [ ] 将旧 Kratos 日志辅助 API 替换为基于 `log/slog` 的 API。
- [ ] 将 JWT 导入路径迁移到 `contrib/middleware/jwt/v3`。
- [ ] 检查自定义熔断器依赖。
- [ ] 替换直接使用的 `transport/http/binding`。
- [ ] 重新生成生成代码。
- [ ] 运行测试、lint 和 `go mod tidy`。
