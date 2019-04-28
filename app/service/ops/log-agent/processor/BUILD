package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["processor.go"],
    importpath = "go-common/app/service/ops/log-agent/processor",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/ops/log-agent/event:go_default_library",
        "//app/service/ops/log-agent/output:go_default_library",
        "//library/log:go_default_library",
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
        "//app/service/ops/log-agent/processor/classify:all-srcs",
        "//app/service/ops/log-agent/processor/fileLog:all-srcs",
        "//app/service/ops/log-agent/processor/grok:all-srcs",
        "//app/service/ops/log-agent/processor/httpstream:all-srcs",
        "//app/service/ops/log-agent/processor/jsonLog:all-srcs",
        "//app/service/ops/log-agent/processor/lengthCheck:all-srcs",
        "//app/service/ops/log-agent/processor/sample:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
