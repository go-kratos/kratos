module github.com/go-kratos/kratos/contrib/encoding/cbor/v2

go 1.19

require (
	github.com/fxamacker/cbor/v2 v2.5.0
	github.com/go-kratos/kratos/v2 v2.7.0
)

require github.com/x448/float16 v0.8.4 // indirect

replace github.com/go-kratos/kratos/v2 => ../../../
