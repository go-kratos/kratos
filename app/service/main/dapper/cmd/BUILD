package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "cmd",
    embed = [":go_default_library"],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    data = ["dapper-service-example.toml"],
    importpath = "go-common/app/service/main/dapper/cmd",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/dapper/collector:go_default_library",
        "//app/service/main/dapper/conf:go_default_library",
        "//app/service/main/dapper/pkg/util:go_default_library",
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
        "//app/service/main/dapper/cmd/agent:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
