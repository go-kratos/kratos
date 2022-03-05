::go install github.com/go-kratos/kratos/cmd/kratos/v2@latestgo mod download
::go mod download

@echo on

set API_PATH="../../../api/logger/job/v1"

:: 生成 proto grpc
protoc --proto_path=%API_PATH% ^
       --proto_path=../../../third_party ^
       --go_out=paths=source_relative:%API_PATH% ^
       --go-grpc_out=paths=source_relative:%API_PATH% ^
       %API_PATH%/*.proto

:: 生成 proto errors
protoc --proto_path=%API_PATH% ^
       --proto_path=../../../third_party ^
       --go_out=paths=source_relative:%API_PATH% ^
       --go-errors_out=paths=source_relative:%API_PATH% ^
       %API_PATH%/*.proto

:: 生成 swagger
protoc --proto_path=%API_PATH% ^
       --proto_path=../../../third_party ^
       --openapiv2_out %API_PATH% ^
       --openapiv2_opt logtostderr=true --openapiv2_opt json_names_for_fields=false ^
       %API_PATH%/*.proto

:: 生成配置
protoc --proto_path=. ^
       --proto_path=../../../third_party ^
       --go_out=paths=source_relative:. ^
       internal/conf/*.proto

:: 生成ent
go generate ./internal/data/ent

:: 生成wire
:: go generate wire
:: go generate ./...

:: 运行项目
:: kratos run

pause
