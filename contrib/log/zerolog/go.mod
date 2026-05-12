module github.com/go-kratos/kratos/contrib/log/zerolog/v2

go 1.25.0

require github.com/rs/zerolog v1.35.1

require (
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.22 // indirect
	golang.org/x/sys v0.44.0 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../..
