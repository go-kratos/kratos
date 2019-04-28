package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "base_model.go",
        "model.go",
        "request_base.go",
    ],
    importpath = "go-common/app/interface/main/mcn/model",
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
        "//app/interface/main/mcn/model/datamodel:all-srcs",
        "//app/interface/main/mcn/model/mcnmodel:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
