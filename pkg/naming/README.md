# naming

## 项目简介

服务发现、服务注册相关的SDK集合

## 现状

目前默认实现了B站开源的[Discovery](https://github.com/bilibili/discovery)服务注册与发现SDK。
但在使用之前，请确认discovery服务部署完成，并将该discovery.go内`fixConfig`方法的默认配置进行完善。

## 使用

可实现`naming`内的`Builder`&`Resolver`&`Registry`接口用于服务注册与发现，比如B站内部还实现了zk的。
