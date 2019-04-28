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
    data = [
        "archive-service-example.toml",
        "archive-service-test.toml",
    ],
    importpath = "go-common/app/service/main/archive/cmd",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/archive/conf:go_default_library",
        "//app/service/main/archive/server/gorpc:go_default_library",
        "//app/service/main/archive/server/grpc:go_default_library",
        "//app/service/main/archive/server/http:go_default_library",
        "//app/service/main/archive/service:go_default_library",
        "//library/ecode/tip:go_default_library",
        "//library/log:go_default_library",
        "//library/net/trace:go_default_library",
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
