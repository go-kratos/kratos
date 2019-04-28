package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "config.go",
        "file.go",
        "glob.go",
        "harvester.go",
        "multiline.go",
        "reader.go",
        "registrar.go",
        "state.go",
        "states.go",
    ],
    importpath = "go-common/app/service/ops/log-agent/input/file",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/ops/log-agent/event:go_default_library",
        "//app/service/ops/log-agent/input:go_default_library",
        "//app/service/ops/log-agent/pkg/lancerroute:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
