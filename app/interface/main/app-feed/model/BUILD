package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["const.go"],
    importpath = "go-common/app/interface/main/app-feed/model",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-card/model/card/live:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
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
        "//app/interface/main/app-feed/model/dislike:all-srcs",
        "//app/interface/main/app-feed/model/feed:all-srcs",
        "//app/interface/main/app-feed/model/live:all-srcs",
        "//app/interface/main/app-feed/model/region:all-srcs",
        "//app/interface/main/app-feed/model/search:all-srcs",
        "//app/interface/main/app-feed/model/tag:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
