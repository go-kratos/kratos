package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
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
        "//app/service/main/broadcast/api/grpc/v1:all-srcs",
        "//app/service/main/broadcast/cmd:all-srcs",
        "//app/service/main/broadcast/dao:all-srcs",
        "//app/service/main/broadcast/libs:all-srcs",
        "//app/service/main/broadcast/model:all-srcs",
        "//app/service/main/broadcast/server/grpc:all-srcs",
        "//app/service/main/broadcast/server/http:all-srcs",
        "//app/service/main/broadcast/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
