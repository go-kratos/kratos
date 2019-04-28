load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_binary",
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    data = ["dapper-query-example.toml"],
    importpath = "go-common/app/service/main/dapper-query/cmd",
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
    deps = [
        "//app/service/main/dapper-query/api/v1:go_default_library",
        "//app/service/main/dapper-query/conf:go_default_library",
        "//app/service/main/dapper-query/service:go_default_library",
        "//app/service/main/dapper-query/util:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
    ],
)

go_binary(
    name = "cmd",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
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
