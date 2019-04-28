package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "feed.go",
        "player.go",
    ],
    importpath = "go-common/app/interface/main/app-feed/model/feed",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-card/model:go_default_library",
        "//app/interface/main/app-card/model/card/ai:go_default_library",
        "//app/interface/main/app-card/model/card/audio:go_default_library",
        "//app/interface/main/app-card/model/card/bangumi:go_default_library",
        "//app/interface/main/app-card/model/card/banner:go_default_library",
        "//app/interface/main/app-card/model/card/cm:go_default_library",
        "//app/interface/main/app-card/model/card/live:go_default_library",
        "//app/interface/main/app-card/model/card/operate:go_default_library",
        "//app/interface/main/app-card/model/card/rank:go_default_library",
        "//app/interface/main/app-card/model/card/show:go_default_library",
        "//app/interface/main/app-feed/model:go_default_library",
        "//app/interface/main/app-feed/model/dislike:go_default_library",
        "//app/interface/main/app-feed/model/live:go_default_library",
        "//app/interface/main/tag/model:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/feed/model:go_default_library",
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
