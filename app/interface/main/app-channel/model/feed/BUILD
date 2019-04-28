package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "feed.go",
        "feed2.go",
    ],
    importpath = "go-common/app/interface/main/app-channel/model/feed",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-card/model/card:go_default_library",
        "//app/interface/main/app-card/model/card/audio:go_default_library",
        "//app/interface/main/app-card/model/card/bangumi:go_default_library",
        "//app/interface/main/app-card/model/card/live:go_default_library",
        "//app/interface/main/app-card/model/card/operate:go_default_library",
        "//app/interface/main/app-card/model/card/show:go_default_library",
        "//app/interface/main/app-channel/model:go_default_library",
        "//app/interface/main/app-channel/model/activity:go_default_library",
        "//app/interface/main/app-channel/model/card:go_default_library",
        "//app/interface/main/app-channel/model/channel:go_default_library",
        "//app/interface/main/app-channel/model/dislike:go_default_library",
        "//app/interface/main/app-channel/model/recommend:go_default_library",
        "//app/interface/main/tag/model:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//app/service/openplatform/pgc-season/api/grpc/episode/v1:go_default_library",
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
