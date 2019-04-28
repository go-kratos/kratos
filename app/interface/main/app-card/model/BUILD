package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["const.go"],
    importpath = "go-common/app/interface/main/app-card/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-card/model/card/audio:go_default_library",
        "//app/interface/main/app-card/model/card/live:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/log:go_default_library",
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
        "//app/interface/main/app-card/model/bplus:all-srcs",
        "//app/interface/main/app-card/model/card:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
