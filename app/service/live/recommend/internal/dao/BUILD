package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "online_filter.go",
        "room_feature.go",
    ],
    importpath = "go-common/app/service/live/recommend/internal/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/recommend/internal/conf:go_default_library",
        "//app/service/live/recommend/recconst:go_default_library",
        "//app/service/live/relation/api/liverpc:go_default_library",
        "//app/service/live/relation/api/liverpc/v1:go_default_library",
        "//app/service/live/room/api/liverpc:go_default_library",
        "//app/service/live/room/api/liverpc/v1:go_default_library",
        "//app/service/live/room/api/liverpc/v2:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/log:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "dao_test.go",
        "room_feature_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/live/recommend/internal/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
