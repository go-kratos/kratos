package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "conn.go",
        "request.go",
        "server.go",
    ],
    importpath = "go-common/app/service/main/broadcast/libs/websocket",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//app/service/main/broadcast/libs/bufio:go_default_library"],
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
