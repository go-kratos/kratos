package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["const.go"],
    importpath = "go-common/app/interface/main/app-interface/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-interface/model/live:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/conf/env:go_default_library",
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
        "//app/interface/main/app-interface/model/audio:all-srcs",
        "//app/interface/main/app-interface/model/bangumi:all-srcs",
        "//app/interface/main/app-interface/model/banner:all-srcs",
        "//app/interface/main/app-interface/model/bplus:all-srcs",
        "//app/interface/main/app-interface/model/coin:all-srcs",
        "//app/interface/main/app-interface/model/community:all-srcs",
        "//app/interface/main/app-interface/model/elec:all-srcs",
        "//app/interface/main/app-interface/model/favorite:all-srcs",
        "//app/interface/main/app-interface/model/history:all-srcs",
        "//app/interface/main/app-interface/model/live:all-srcs",
        "//app/interface/main/app-interface/model/pay:all-srcs",
        "//app/interface/main/app-interface/model/region:all-srcs",
        "//app/interface/main/app-interface/model/relation:all-srcs",
        "//app/interface/main/app-interface/model/search:all-srcs",
        "//app/interface/main/app-interface/model/shop:all-srcs",
        "//app/interface/main/app-interface/model/sidebar:all-srcs",
        "//app/interface/main/app-interface/model/sp:all-srcs",
        "//app/interface/main/app-interface/model/space:all-srcs",
        "//app/interface/main/app-interface/model/tag:all-srcs",
        "//app/interface/main/app-interface/model/topic:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
