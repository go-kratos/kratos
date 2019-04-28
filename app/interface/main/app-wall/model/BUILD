package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["const.go"],
    importpath = "go-common/app/interface/main/app-wall/model",
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
        "//app/interface/main/app-wall/model/mobile:all-srcs",
        "//app/interface/main/app-wall/model/operator:all-srcs",
        "//app/interface/main/app-wall/model/telecom:all-srcs",
        "//app/interface/main/app-wall/model/unicom:all-srcs",
        "//app/interface/main/app-wall/model/wall:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
