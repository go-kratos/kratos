# Kratos

This document defines the roadmap for Kratos development.

## Features
- [x] Config
    - [x] Local Files
    - [x] K8s ConfigMap
    - [x] Consul
    - [x] Etcd
    - [x] Nacos
- [x] Registry
    - [x] Consul
    - [x] Etcd
    - [x] K8s
    - [x] Nacos
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
    - [x] Validator
    - [x] Authentication
    - [x] Ratelimit
    - [x] CircuitBreaker
- [x] Metrics
    - [x] Prometheus
    - [x] DataDog
- [x] Tracing
    - [x] HTTP
        - [x] TLS
        - [x] Client
        - [x] Service Registrar
        - [ ] javascript/typescript clients
    - [x] gRPC
        - [x] TLS
        - [x] Unary Handler
        - [x] Streaming Handler
- [ ] Cache
    - [ ] go-redis
- [x] Event
    - [x] Pub/Sub
    - [x] Kafka
    - [ ] Nats
- [x] Database
    - [x] Ent
    - [ ] Gorm

## Platform
- [ ] Kratos API
    - [ ] Auth
    - [ ] Config
    - [ ] Registry
    - [ ] Events
- [ ] Kratos Runtime
    - [ ] Secrets
    - [ ] Service-to-Service
    - [ ] Publish and Subscribe
    - [ ] Observability
    - [ ] Controllable
- [ ] Kratos UI
    - [ ] Auth
    - [ ] Config
    - [ ] Services
    - [ ] Endpoints
    - [ ] Ratelimit
    - [ ] CircuitBreaker
    - [ ] FaultInjection
    - [ ] TrafficPolicy

## Tools
- [x] Kratos
- [x] HTTP Generator
    - [ ] API YAML
- [x] Errors Generator
