package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "aggr.go",
        "config.go",
        "lancer.go",
        "pool.go",
    ],
    importpath = "go-common/app/service/ops/log-agent/output/lancergrpc",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/ops/log-agent/event:go_default_library",
        "//app/service/ops/log-agent/output:go_default_library",
        "//app/service/ops/log-agent/output/cache/file:go_default_library",
        "//app/service/ops/log-agent/output/lancergrpc/lancergateway:go_default_library",
        "//app/service/ops/log-agent/pkg/bufio:go_default_library",
        "//app/service/ops/log-agent/pkg/common:go_default_library",
        "//app/service/ops/log-agent/pkg/flowmonitor:go_default_library",
        "//app/service/ops/log-agent/pkg/lancermonitor:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/BurntSushi/toml:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
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
        "//app/service/ops/log-agent/output/lancergrpc/lancergateway:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
