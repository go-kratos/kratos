load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_binary",
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    data = ["secure-service-test.toml"],
    importpath = "go-common/app/service/main/secure/cmd",
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
    deps = [
        "//app/service/main/secure/conf:go_default_library",
        "//app/service/main/secure/http:go_default_library",
        "//app/service/main/secure/rpc/server:go_default_library",
        "//app/service/main/secure/service:go_default_library",
        "//library/log:go_default_library",
        "//library/net/rpc:go_default_library",
        "//library/net/trace:go_default_library",
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
