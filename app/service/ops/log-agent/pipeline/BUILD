package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "management.go",
        "pipeline.go",
    ],
    importpath = "go-common/app/service/ops/log-agent/pipeline",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/ops/log-agent/event:go_default_library",
        "//app/service/ops/log-agent/input:go_default_library",
        "//app/service/ops/log-agent/input/file:go_default_library",
        "//app/service/ops/log-agent/input/sock:go_default_library",
        "//app/service/ops/log-agent/output:go_default_library",
        "//app/service/ops/log-agent/output/lancergrpc:go_default_library",
        "//app/service/ops/log-agent/output/lancerlogstream:go_default_library",
        "//app/service/ops/log-agent/output/stdout:go_default_library",
        "//app/service/ops/log-agent/pkg/common:go_default_library",
        "//app/service/ops/log-agent/processor:go_default_library",
        "//app/service/ops/log-agent/processor/classify:go_default_library",
        "//app/service/ops/log-agent/processor/fileLog:go_default_library",
        "//app/service/ops/log-agent/processor/grok:go_default_library",
        "//app/service/ops/log-agent/processor/httpstream:go_default_library",
        "//app/service/ops/log-agent/processor/jsonLog:go_default_library",
        "//app/service/ops/log-agent/processor/lengthCheck:go_default_library",
        "//app/service/ops/log-agent/processor/sample:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/BurntSushi/toml:go_default_library",
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
        "//app/service/ops/log-agent/pipeline/dockerlogcollector:all-srcs",
        "//app/service/ops/log-agent/pipeline/hostlogcollector:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
