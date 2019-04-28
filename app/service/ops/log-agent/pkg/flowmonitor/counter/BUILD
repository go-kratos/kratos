package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "collector.go",
        "counter.go",
        "desc.go",
        "fnv.go",
        "go_collector.go",
        "metric.go",
        "observer.go",
        "process_collector.go",
        "registry.go",
        "summary.go",
        "untyped.go",
        "value.go",
        "vec.go",
    ],
    importpath = "go-common/app/service/ops/log-agent/pkg/flowmonitor/counter",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/beorn7/perks/quantile:go_default_library",
        "//vendor/github.com/prometheus/client_model/go:go_default_library",
        "//vendor/github.com/prometheus/common/model:go_default_library",
        "//vendor/github.com/prometheus/procfs:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
