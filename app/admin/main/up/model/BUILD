package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "command_arg.go",
        "errorcode.go",
        "monitor.go",
        "msg.go",
        "rpc.go",
    ],
    importpath = "go-common/app/admin/main/up/model",
    tags = ["automanaged"],
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
    srcs = [
        ":package-srcs",
        "//app/admin/main/up/model/datamodel:all-srcs",
        "//app/admin/main/up/model/signmodel:all-srcs",
        "//app/admin/main/up/model/upcrmmodel:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
