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
    data = ["config-service-example.toml"],
    importpath = "go-common/app/infra/config/cmd",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/infra/config/conf:go_default_library",
        "//app/infra/config/http:go_default_library",
        "//app/infra/config/rpc/server:go_default_library",
        "//app/infra/config/service/v1:go_default_library",
        "//app/infra/config/service/v2:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/log:go_default_library",
        "//library/naming:go_default_library",
        "//library/naming/discovery:go_default_library",
        "//library/net/ip:go_default_library",
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
