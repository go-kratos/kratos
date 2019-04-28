package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "account_test.go",
        "album_test.go",
        "app_test.go",
        "archive_test.go",
        "article_test.go",
        "bangumi_test.go",
        "channel_archive_test.go",
        "channel_test.go",
        "coin_test.go",
        "dynamic_test.go",
        "fav_test.go",
        "notice_test.go",
        "service_test.go",
        "setting_test.go",
        "tag_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/space/conf:go_default_library",
        "//app/interface/main/space/model:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "account.go",
        "album.go",
        "app.go",
        "archive.go",
        "article.go",
        "assist.go",
        "bangumi.go",
        "blacklist.go",
        "channel.go",
        "channel_archive.go",
        "coin.go",
        "dynamic.go",
        "fav.go",
        "game.go",
        "notice.go",
        "service.go",
        "setting.go",
        "shop.go",
        "tag.go",
    ],
    importpath = "go-common/app/interface/main/space/service",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/space/conf:go_default_library",
        "//app/interface/main/space/dao:go_default_library",
        "//app/interface/main/space/model:go_default_library",
        "//app/interface/main/tag/model:go_default_library",
        "//app/interface/main/tag/rpc/client:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/interface/openplatform/article/rpc/client:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/assist/model/assist:go_default_library",
        "//app/service/main/assist/rpc/client:go_default_library",
        "//app/service/main/coin/api:go_default_library",
        "//app/service/main/favorite/api/gorpc:go_default_library",
        "//app/service/main/favorite/model:go_default_library",
        "//app/service/main/filter/model/rpc:go_default_library",
        "//app/service/main/filter/rpc/client:go_default_library",
        "//app/service/main/member/api/gorpc:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//app/service/main/member/model/block:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//app/service/main/relation/rpc/client:go_default_library",
        "//app/service/main/thumbup/model:go_default_library",
        "//app/service/main/thumbup/rpc/client:go_default_library",
        "//app/service/main/up/api/v1:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/sync/errgroup.v2:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
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

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)
