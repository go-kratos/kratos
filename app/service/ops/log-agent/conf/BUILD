package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["conf.go"],
    importpath = "go-common/app/service/ops/log-agent/conf",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/ops/log-agent/conf/configcenter:go_default_library",
        "//app/service/ops/log-agent/pipeline/dockerlogcollector:go_default_library",
        "//app/service/ops/log-agent/pipeline/hostlogcollector:go_default_library",
        "//app/service/ops/log-agent/pkg/flowmonitor:go_default_library",
        "//app/service/ops/log-agent/pkg/httpstream:go_default_library",
        "//app/service/ops/log-agent/pkg/lancermonitor:go_default_library",
        "//app/service/ops/log-agent/pkg/limit:go_default_library",
        "//library/log:go_default_library",
        "//library/naming/discovery:go_default_library",
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
        "//app/service/ops/log-agent/conf/configcenter:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
