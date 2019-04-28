package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["const.go"],
    importpath = "go-common/app/interface/main/app-channel/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-card/model/card/live:go_default_library",
        "//app/interface/main/app-channel/model/tab:go_default_library",
        "//app/service/main/archive/api:go_default_library",
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
        "//app/interface/main/app-channel/model/activity:all-srcs",
        "//app/interface/main/app-channel/model/audio:all-srcs",
        "//app/interface/main/app-channel/model/bangumi:all-srcs",
        "//app/interface/main/app-channel/model/card:all-srcs",
        "//app/interface/main/app-channel/model/channel:all-srcs",
        "//app/interface/main/app-channel/model/converge:all-srcs",
        "//app/interface/main/app-channel/model/dislike:all-srcs",
        "//app/interface/main/app-channel/model/feed:all-srcs",
        "//app/interface/main/app-channel/model/game:all-srcs",
        "//app/interface/main/app-channel/model/live:all-srcs",
        "//app/interface/main/app-channel/model/recommend:all-srcs",
        "//app/interface/main/app-channel/model/shopping:all-srcs",
        "//app/interface/main/app-channel/model/special:all-srcs",
        "//app/interface/main/app-channel/model/tab:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
