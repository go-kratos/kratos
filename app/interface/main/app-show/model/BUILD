package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["const.go"],
    importpath = "go-common/app/interface/main/app-show/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//app/service/main/archive/api:go_default_library"],
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
        "//app/interface/main/app-show/model/activity:all-srcs",
        "//app/interface/main/app-show/model/ad:all-srcs",
        "//app/interface/main/app-show/model/bangumi:all-srcs",
        "//app/interface/main/app-show/model/banner:all-srcs",
        "//app/interface/main/app-show/model/card:all-srcs",
        "//app/interface/main/app-show/model/daily:all-srcs",
        "//app/interface/main/app-show/model/dislike:all-srcs",
        "//app/interface/main/app-show/model/feed:all-srcs",
        "//app/interface/main/app-show/model/live:all-srcs",
        "//app/interface/main/app-show/model/recommend:all-srcs",
        "//app/interface/main/app-show/model/region:all-srcs",
        "//app/interface/main/app-show/model/show:all-srcs",
        "//app/interface/main/app-show/model/tag:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
