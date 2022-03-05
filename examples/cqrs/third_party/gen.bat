set API_PATH="pagination"

:: 生成 proto
protoc --proto_path=%API_PATH% ^
       --go_out=paths=source_relative:%API_PATH% ^
       %API_PATH%/*.proto

       pause