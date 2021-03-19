# Kratos

This document defines the roadmap for Kratos development.

## Features
- [x] Config
    - [x] Local Files
    - [x] K8s ConfigMap
    - [ ] Consul
    - [ ] Etcd
    - [ ] Nacos
- [ ] Registry
    - [ ] Consul
    - [ ] Etcd
    - [ ] K8s
    - [ ] Nacos
- [x] Encoding
    - [x] JSON
    - [x] Protobuf
- [x] Transport
    - [x] HTTP
    - [x] gRPC
- [x] Middleware
    - [x] Logging
    - [x] metrics
    - [x] recovery
    - [x] gRPC status
    - [x] transport tracing
    - [x] validator
    - [ ] authentication
    - [ ] ratelimit
    - [ ] circuitbreaker
- [x] Metrics
    - [x] Prometheus
    - [x] DataDog
- [x] Tracing
    - [x] HTTP
        - [ ] TLS
        - [x] Client
        - [x] Service Registrar
        - [ ] javascript/typescript clients
    - [x] gRPC
        - [ ] TLS
        - [x] Uarry Handler
        - [ ] Streaming Handler
- [ ] Cache
    - [ ] go-redis
- [x] Pubsub
    - [x] Absctraction
    - [x] Kafka
    - [ ] Nats
- [ ] Database
    - [ ] Ent
    - [ ] Gorm

## Tools
- [x] Kratos
- [x] HTTP Generator
    - [ ] API YAML
- [x] Errors Generator
