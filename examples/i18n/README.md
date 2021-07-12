# i18n Example
This project is an demo for i18n usage.
We integrated [go-i18n](https://github.com/nicksnyder/go-i18n) into our project for localization.

We start with the project generate with `kratos-layout`, you could use `kratos new` to generate it.

The following steps are applied.
## Step 1: install 
Install the cli tool `goi18n` and the package.

```bash
go get -u github.com/nicksnyder/go-i18n/v2/goi18n
go get -u github.com/nicksnyder/go-i18n
```

## Step 2: implement and register the middleware
There's a middleware in `internal/pkg/middleware/localize`, 
we implement this middleware, and register it into our http server in `internal/server/http.go`.
We init the bundle and load our translation file(we will generate it later).
For convenience, we hard-coded the translation file path here, you should avoid it in your project :)

Also the middleware will extract `accept-language` header
from our request, then set the correct `localizer` to the context. 
In our service request handler, we can use the `FromContext` method to get this `localizer` for localization.

This middleware is so simple that we can show it as following:
```go
package localize

import (
	"context"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type localizerKey struct{}

func I18N() middleware.Middleware {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("../../active.zh.toml")

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				accept := tr.Header().Get("accept-language")
				println(accept)
				localizer := i18n.NewLocalizer(bundle, accept)
				ctx = context.WithValue(ctx, localizerKey{}, localizer)
			}
			return handler(ctx, req)
		}
	}
}

func FromContext(ctx context.Context) *i18n.Localizer {
	return ctx.Value(localizerKey{}).(*i18n.Localizer)
}
```

## Step 3:
Write our request handler code. 
In `internal/service/greeter.go`, we use `localize.FromContext(ctx)` to get the localizer, and write the string for localizing.
You should write your message with golang template syntax.
Notice that both `One` and `Other` fields must be filled or you will get an panic on translation file generation with `goi18n` tool.

```go
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	s.log.WithContext(ctx).Infof("SayHello Received: %v", in.GetName())

	if in.GetName() == "error" {
		return nil, v1.ErrorUserNotFound("user not found: %s", in.GetName())
	}
	localizer := localize.FromContext(ctx)
	helloMsg, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			Description: "sayhello",
			ID:    "sayHello",
			One:   "Hello {{.Name}}",
			Other: "Hello {{.Name}}",
		},
		TemplateData: map[string]interface{}{
			"Name": in.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: helloMsg}, nil
}
```

## Step 3: Generate and translate

In the root of this project, generate `active.en.toml` with
```bash
goi18n extract
```
The string which should be translated will be extracted from your code, and write into `active.en.toml`
You should create the empty target language file:
```bash
touch tranlate.zh.toml
```

Then fill the translate file:
```bash
goi18n merge active.en.toml translate.zh.toml
```

You should edit the `translate.zh.toml` to finish the translation work. We translate `Hello` to `你好` in Chinese.
After that, you should rename this file to `active.zh.toml`
And this file is the translation file that we load in the Step 2.
You could also embed these translation files into your binaries for easier deployment.

## Step 4: Run
Go to `cmd/i18n/`, and run `go run .` to start the service.

You could try with curl with `Accept-Language` header
```bash
curl "http://localhost:8000/helloworld/eric" \
     -H 'Accept-Language: zh-CN'
```
Will get the Chinese result `{"message":"你好 eric"}`

And if no header:
```
curl "http://localhost:8000/helloworld/eric"
```
Will get the default English result `{"message":"Hello eric"}`

## Reference
* [go-i18n](https://github.com/nicksnyder/go-i18n) You could refer to this repository for more detailed document.