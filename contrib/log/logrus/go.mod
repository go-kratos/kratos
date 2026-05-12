module github.com/go-kratos/kratos/contrib/log/logrus/v2

go 1.25.0

require github.com/sirupsen/logrus v1.9.4

require (
	github.com/go-kratos/kratos/v2 v2.9.2
	golang.org/x/sys v0.44.0 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../..
