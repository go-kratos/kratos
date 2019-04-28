package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "AppRelationFilterStrategy.go",
        "AppRelationSortStrategy.go",
        "Tools.go",
        "rpcWraper.go",
    ],
    importpath = "go-common/app/interface/live/app-interface/service/v1/relation",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/app-interface/api/http/v1:go_default_library",
        "//app/interface/live/app-interface/conf:go_default_library",
        "//app/interface/live/app-interface/dao:go_default_library",
        "//app/service/live/av/api/liverpc/v1:go_default_library",
        "//app/service/live/fans_medal/api/liverpc/v1:go_default_library",
        "//app/service/live/live_data/api/liverpc/v1:go_default_library",
        "//app/service/live/relation/api/liverpc/v1:go_default_library",
        "//app/service/live/room/api/liverpc/v1:go_default_library",
        "//app/service/live/room/api/liverpc/v2:go_default_library",
        "//app/service/live/room_ex/api/liverpc/v1:go_default_library",
        "//app/service/live/userext/api/liverpc/v1:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/rpc/liverpc/context:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
