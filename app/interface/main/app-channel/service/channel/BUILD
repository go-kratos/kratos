package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["channel_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-channel/conf:go_default_library",
        "//app/interface/main/app-channel/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "audit.go",
        "cache.go",
        "channel.go",
        "index.go",
        "index2.go",
        "infoc.go",
        "manager.go",
        "service.go",
        "tab.go",
    ],
    importpath = "go-common/app/interface/main/app-channel/service/channel",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-card/model:go_default_library",
        "//app/interface/main/app-card/model/card:go_default_library",
        "//app/interface/main/app-card/model/card/audio:go_default_library",
        "//app/interface/main/app-card/model/card/bangumi:go_default_library",
        "//app/interface/main/app-card/model/card/live:go_default_library",
        "//app/interface/main/app-card/model/card/operate:go_default_library",
        "//app/interface/main/app-card/model/card/show:go_default_library",
        "//app/interface/main/app-channel/conf:go_default_library",
        "//app/interface/main/app-channel/dao/account:go_default_library",
        "//app/interface/main/app-channel/dao/activity:go_default_library",
        "//app/interface/main/app-channel/dao/archive:go_default_library",
        "//app/interface/main/app-channel/dao/article:go_default_library",
        "//app/interface/main/app-channel/dao/audio:go_default_library",
        "//app/interface/main/app-channel/dao/audit:go_default_library",
        "//app/interface/main/app-channel/dao/bangumi:go_default_library",
        "//app/interface/main/app-channel/dao/card:go_default_library",
        "//app/interface/main/app-channel/dao/converge:go_default_library",
        "//app/interface/main/app-channel/dao/game:go_default_library",
        "//app/interface/main/app-channel/dao/live:go_default_library",
        "//app/interface/main/app-channel/dao/location:go_default_library",
        "//app/interface/main/app-channel/dao/region:go_default_library",
        "//app/interface/main/app-channel/dao/relation:go_default_library",
        "//app/interface/main/app-channel/dao/shopping:go_default_library",
        "//app/interface/main/app-channel/dao/special:go_default_library",
        "//app/interface/main/app-channel/dao/tab:go_default_library",
        "//app/interface/main/app-channel/dao/tag:go_default_library",
        "//app/interface/main/app-channel/model:go_default_library",
        "//app/interface/main/app-channel/model/activity:go_default_library",
        "//app/interface/main/app-channel/model/card:go_default_library",
        "//app/interface/main/app-channel/model/channel:go_default_library",
        "//app/interface/main/app-channel/model/feed:go_default_library",
        "//app/interface/main/app-channel/model/tab:go_default_library",
        "//app/interface/main/tag/model:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/location/model:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//app/service/openplatform/pgc-season/api/grpc/episode/v1:go_default_library",
        "//app/service/openplatform/pgc-season/api/grpc/season/v1:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//vendor/github.com/dgryski/go-farm:go_default_library",
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
