package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "card.go",
        "double.go",
        "ipad.go",
        "single.go",
    ],
    importpath = "go-common/app/interface/main/app-card/model/card",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-card/model:go_default_library",
        "//app/interface/main/app-card/model/bplus:go_default_library",
        "//app/interface/main/app-card/model/card/ai:go_default_library",
        "//app/interface/main/app-card/model/card/audio:go_default_library",
        "//app/interface/main/app-card/model/card/bangumi:go_default_library",
        "//app/interface/main/app-card/model/card/banner:go_default_library",
        "//app/interface/main/app-card/model/card/cm:go_default_library",
        "//app/interface/main/app-card/model/card/live:go_default_library",
        "//app/interface/main/app-card/model/card/operate:go_default_library",
        "//app/interface/main/app-card/model/card/show:go_default_library",
        "//app/interface/main/tag/model:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//app/service/openplatform/pgc-season/api/grpc/episode/v1:go_default_library",
        "//app/service/openplatform/pgc-season/api/grpc/season/v1:go_default_library",
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
        "//app/interface/main/app-card/model/card/ai:all-srcs",
        "//app/interface/main/app-card/model/card/audio:all-srcs",
        "//app/interface/main/app-card/model/card/bangumi:all-srcs",
        "//app/interface/main/app-card/model/card/banner:all-srcs",
        "//app/interface/main/app-card/model/card/cm:all-srcs",
        "//app/interface/main/app-card/model/card/live:all-srcs",
        "//app/interface/main/app-card/model/card/operate:all-srcs",
        "//app/interface/main/app-card/model/card/rank:all-srcs",
        "//app/interface/main/app-card/model/card/show:all-srcs",
        "//app/interface/main/app-card/model/card/topic:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
