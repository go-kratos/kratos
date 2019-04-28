package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "banner.go",
        "capsule.go",
        "gaea.go",
        "payGoods.go",
        "payLive.go",
        "resource.go",
        "roomMng.go",
        "splash.go",
        "token.go",
        "upload.go",
    ],
    importpath = "go-common/app/admin/live/live-admin/service/v1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/live/live-admin/api/http/v1:go_default_library",
        "//app/admin/live/live-admin/conf:go_default_library",
        "//app/admin/live/live-admin/dao:go_default_library",
        "//app/admin/live/live-admin/model:go_default_library",
        "//app/service/live/av/api/liverpc/v0:go_default_library",
        "//app/service/live/relation/api/liverpc/v1:go_default_library",
        "//app/service/live/resource/api/grpc/v1:go_default_library",
        "//app/service/live/room/api/liverpc:go_default_library",
        "//app/service/live/room/api/liverpc/v1:go_default_library",
        "//app/service/live/xlottery/api/grpc/v1:go_default_library",
        "//app/service/video/stream-mng/api/v1:go_default_library",
        "//library/database/bfs:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = ["token_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/admin/live/live-admin/api/http/v1:go_default_library",
        "//app/admin/live/live-admin/conf:go_default_library",
        "//app/admin/live/live-admin/dao:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
