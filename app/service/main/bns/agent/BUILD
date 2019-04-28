package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "agent.go",
        "dns.go",
        "http.go",
        "ns_endpoint.go",
    ],
    importpath = "go-common/app/service/main/bns/agent",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/bns/agent/backend:go_default_library",
        "//app/service/main/bns/conf:go_default_library",
        "//app/service/main/bns/lib/resolvconf:go_default_library",
        "//app/service/main/bns/lib/shuffle:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/log:go_default_library",
        "//library/stat/prom:go_default_library",
        "//vendor/github.com/miekg/dns:go_default_library",
        "//vendor/github.com/prometheus/client_golang/prometheus/promhttp:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//app/service/main/bns/agent/backend:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
