package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["const.go"],
    importpath = "go-common/app/interface/main/app-resource/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//app/interface/main/app-resource/model/tab:go_default_library"],
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
        "//app/interface/main/app-resource/model/abtest:all-srcs",
        "//app/interface/main/app-resource/model/broadcast:all-srcs",
        "//app/interface/main/app-resource/model/domain:all-srcs",
        "//app/interface/main/app-resource/model/experiment:all-srcs",
        "//app/interface/main/app-resource/model/guide:all-srcs",
        "//app/interface/main/app-resource/model/module:all-srcs",
        "//app/interface/main/app-resource/model/notice:all-srcs",
        "//app/interface/main/app-resource/model/param:all-srcs",
        "//app/interface/main/app-resource/model/plugin:all-srcs",
        "//app/interface/main/app-resource/model/show:all-srcs",
        "//app/interface/main/app-resource/model/sidebar:all-srcs",
        "//app/interface/main/app-resource/model/splash:all-srcs",
        "//app/interface/main/app-resource/model/static:all-srcs",
        "//app/interface/main/app-resource/model/tab:all-srcs",
        "//app/interface/main/app-resource/model/version:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
