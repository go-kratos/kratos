package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "appeal.go",
        "archive.go",
        "blocked.go",
        "credit.go",
        "credit_state.go",
        "jury.go",
        "labour.go",
        "manager.go",
        "param.go",
        "rpc.go",
    ],
    importpath = "go-common/app/interface/main/credit/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/archive/model/archive:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
