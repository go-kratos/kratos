<p align="center"><a href="https://go-kratos.dev/" target="_blank"><img src="https://github.com/go-kratos/kratos/blob/main/docs/images/kratos-large.png?raw=true"></a></p>

<p align="center">
<a href="https://github.com/go-kratos/kratos/actions"><img src="https://github.com/go-kratos/kratos/workflows/Go/badge.svg" alt="Build Status"></a>
<a href="https://pkg.go.dev/github.com/go-kratos/kratos/v2"><img src="https://pkg.go.dev/badge/github.com/go-kratos/kratos/v2" alt="GoDoc"></a>
<a href="https://codecov.io/gh/go-kratos/kratos"><img src="https://codecov.io/gh/go-kratos/kratos/master/graph/badge.svg" alt="codeCov"></a>
<a href="https://goreportcard.com/report/github.com/go-kratos/kratos"><img src="https://goreportcard.com/badge/github.com/go-kratos/kratos" alt="Go Report Card"></a>
<a href="https://github.com/go-kratos/kratos/blob/main/LICENSE"><img src="https://img.shields.io/github/license/go-kratos/kratos" alt="License"></a>
<a href="https://github.com/avelino/awesome-go"><img src="https://awesome.re/mentioned-badge.svg" alt="Awesome Go"></a>
<a href="https://discord.gg/BWzJsUJ"><img src="https://img.shields.io/discord/766619759214854164?label=chat&logo=discord" alt="Discord"></a>
</p>
<p align="center">
<a href="https://www.producthunt.com/posts/go-kratos?utm_source=badge-featured&utm_medium=badge&utm_souce=badge-go-kratos" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=306565&theme=light" alt="Go Kratos - A Go framework for microservices. | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" /></a>
</p>
  
##### Translate to: [简体中文](README_zh.md)

## About Kratos
  
> The name is inspired by the game God of War which is based on Greek myths, tells the adventures of Kratos who became a god of war from a mortal and launched a god-killing slaughter.

Kratos is a microservice-oriented governance framework implemented by golang, which offers convenient capabilities to help you quickly build a bulletproof application from scratch, such as:

- The [communication protocol](https://go-kratos.dev/en/docs/component/api) is based on the HTTP/gRPC through the definition of Protobuf.
- Abstract [transport](https://go-kratos.dev/en/docs/component/transport/overview) layer are support: [HTTP](https://go-kratos.dev/en/docs/component/transport/http) / [gRPC](https://go-kratos.dev/en/docs/component/transport/grpc).
- Powerful [middleware](https://go-kratos.dev/en/docs/component/middleware/overview) design, support: [Tracing (OpenTelemetry)](https://go-kratos.dev/en/docs/component/middleware/tracing)、[Metrics (Prometheus is default)](https://go-kratos.dev/en/docs/component/middleware/metrics)、[Recovery](https://go-kratos.dev/en/docs/component/middleware/recovery) and more.
- [Registry](https://go-kratos.dev/en/docs/component/registry) interface able to be connected with various other centralized registries through plug-ins.
- The [standard log interfaces](https://go-kratos.dev/en/docs/component/log) ease the integration of the third-party log libs and logs are collected through the *Fluentd*.
- The selection of the content [encoding](https://go-kratos.dev/en/docs/component/encoding) is automatically supported by Accept and Content-Type.
- Multiple data sources are supported for [configurations](https://go-kratos.dev/en/docs/component/config) and dynamic configurations (use atomic operations).
- In the protocol of HTTP/gRPC, use uniform [metadata](https://go-kratos.dev/en/docs/component/metadata) transfer method.
- You can define [errors](https://go-kratos.dev/en/docs/component/errors/) in protos and generate enums with protoc-gen-go.
- You can define [verification rules](https://go-kratos.dev/en/docs/component/middleware/validate) in Protobuf which is supported by HTTP/gRPC service.
- [Swagger API](https://go-kratos.dev/en/docs/guide/openapi) generated Automatically and embed Swagger UI endpoint can be started by adding [Swagger plugin](https://github.com/go-kratos/swagger-api).

Kratos is accessible, powerful, and provides tools required for large, robust applications.

## Learning Kratos

Kratos has the most extensive and thorough [documentation](https://go-kratos.dev/en/docs/getting-started/start) and [example](https://github.com/go-kratos/examples) library of all modern web application frameworks, making it a breeze to get started with the framework.

We also provide a [modern template](https://github.com/go-kratos/kratos-layout), This template should help reduce the work required to setup up a modern project.

### Goals

Kratos boosts your productivity. With the integration of excellent resources and further support, programmers can get rid of most issues might encounter in the field of distributed systems and software engineering such that they are allowed to focus on the release of businesses only. Additionally, for each programmer, Kratos is also an ideal one learning warehouse for many aspects of microservices to enrich their experiences and skills.

### Principles

* **Simple**: Appropriate design, plain and easy code.
* **General**: Cover the various utilities for business development.
* **Highly efficient**: Speeding up the efficiency of businesses upgrading.
* **Stable**: The base libs validated in the production environment which have the characters of the high testability, high coverage as well as high security and reliability.
* **Robust**: Eliminating misusing through high quality of the base libs.
* **High-performance**: Optimal performance excluding the optimization of hacking in case of *unsafe*. 
* **Expandability**: Properly designed interfaces, you can expand utilities such as base libs to meet your further requirements.
* **Fault-tolerance**: Designed against failure, enhance the understanding and exercising of SRE within Kratos to achieve more robustness.
* **Toolchain**: Includes an extensive toolchain, such as the code generation of cache, the lint tool, and so forth.

## Getting Started

Create a kratos playground through [docker](https://www.docker.com/products/docker-desktop):
  
```shell
docker run -it --rm -p 8000:8000 --workdir /workspace golang
```
  
```shell
apt-get update && apt-get -y install protobuf-compiler
export GOPROXY=https://goproxy.io,direct
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest && kratos upgrade
```
  
```shell
kratos new helloworld
cd helloworld/ && go mod tidy
kratos run
```
  
Use a browser to open and visit: `http://localhost:8000/helloworld/kratos`, The kratos program is running!

If you need more, please visit the kratos [documentation](https://go-kratos.dev/en/docs/getting-started/start).

## Security Vulnerabilities

If you discover a security vulnerability within Kratos, please send an e-mail to tonybase via go-kratos@googlegroups.com. All security vulnerabilities will be promptly addressed.

## Community

- [Wechat Group](https://github.com/go-kratos/kratos/issues/682)
- [Discord Group](https://discord.gg/BWzJsUJ)
- [go-kratos.dev](https://go-kratos.dev/en)
- QQ Group: 716486124

## Contributors

Thank you for considering contributing to the Kratos framework! The contribution guide can be found in the [Kratos documentation](https://go-kratos.dev/en/docs/community/contribution).

<a href="https://github.com/go-kratos/kratos/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=go-kratos/kratos" />
</a>

## License

The Kratos framework is open-sourced software licensed under the [MIT license](./LICENSE).
  
## Acknowledgments
  
The following project had particular influence on kratos's design.
  
- [go-kit/kit](https://github.com/go-kit/kit) is a programming toolkit for building microservices in go.
- [asim/go-micro](https://github.com/asim/go-micro) a distributed systems development framework.
- [google/go-cloud](https://github.com/google/go-cloud) is go cloud development kit.
- [tal-tech/go-zero](https://github.com/tal-tech/go-zero) is a web and rpc framework with lots of builtin engineering practices.
- [beego/beego](https://github.com/beego/beego) is a web framework including RESTful APIs, web apps and backend services.
