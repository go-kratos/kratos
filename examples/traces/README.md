# traces example

- There are two services in the example

- After receiving the HTTP request, the user service calls the grpc server provided by the message service through the grpc client

- The example shows how to use tracing middleware

- You can experience the effect by following the command

### run example
```shell
# Windows start-windows.bat
# mac/linux
sh start.sh
curl http://localhost:8000/v1/user/get/message/10
# Open with browser http://localhost:16686
# You can see the effect
```


