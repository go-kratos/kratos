#### business/warden/resolver/direct

##### 项目简介

warden 的直连服务模块，用于通过IP地址列表直接连接后端服务
连接字符串格式：　direct://default/192.168.1.1:8080,192.168.1.2:8081

##### 编译环境

- **请只用 Golang v1.9.x 以上版本编译执行**

##### 依赖包

- [grpc](google.golang.org/grpc)