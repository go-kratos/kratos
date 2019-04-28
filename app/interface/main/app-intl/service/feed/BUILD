package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "black.go",
        "cache.go",
        "index.go",
        "infoc.go",
        "operate.go",
        "player.go",
        "rank.go",
        "service.go",
    ],
    importpath = "go-common/app/interface/main/app-intl/service/feed",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-card/model:go_default_library",
        "//app/interface/main/app-card/model/card:go_default_library",
        "//app/interface/main/app-card/model/card/ai:go_default_library",
        "//app/interface/main/app-card/model/card/banner:go_default_library",
        "//app/interface/main/app-card/model/card/cm:go_default_library",
        "//app/interface/main/app-card/model/card/operate:go_default_library",
        "//app/interface/main/app-card/model/card/rank:go_default_library",
        "//app/interface/main/app-intl/conf:go_default_library",
        "//app/interface/main/app-intl/dao/account:go_default_library",
        "//app/interface/main/app-intl/dao/archive:go_default_library",
        "//app/interface/main/app-intl/dao/black:go_default_library",
        "//app/interface/main/app-intl/dao/card:go_default_library",
        "//app/interface/main/app-intl/dao/location:go_default_library",
        "//app/interface/main/app-intl/dao/rank:go_default_library",
        "//app/interface/main/app-intl/dao/recommend:go_default_library",
        "//app/interface/main/app-intl/dao/relation:go_default_library",
        "//app/interface/main/app-intl/dao/tag:go_default_library",
        "//app/interface/main/app-intl/model:go_default_library",
        "//app/interface/main/app-intl/model/feed:go_default_library",
        "//app/interface/main/tag/model:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/location/model:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/text/translate/chinese:go_default_library",
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
