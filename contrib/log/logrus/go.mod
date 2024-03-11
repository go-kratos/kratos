module github.com/go-kratos/kratos/contrib/log/logrus/v2

go 1.19

require (
	github.com/go-kratos/kratos/v2 v2.7.2
	github.com/sirupsen/logrus v1.8.1
)

require golang.org/x/sys v0.13.0 // indirect

replace github.com/go-kratos/kratos/v2 => ../../../
