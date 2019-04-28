package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "history.go",
        "rpc.go",
        "toview.go",
    ],
    importpath = "go-common/app/interface/main/history/model",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/history/model:go_default_library",
        "//library/ecode:go_default_library",
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
