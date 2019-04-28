package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "model.go",
        "roomAdmin.go",
    ],
    importpath = "go-common/app/service/live/xuser/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/xuser/api/grpc:go_default_library",
        "//library/time:go_default_library",
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
        "//app/service/live/xuser/model/dhh:all-srcs",
        "//app/service/live/xuser/model/exp:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
