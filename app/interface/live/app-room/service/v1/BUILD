load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "gift_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/interface/live/app-room/api/http/v1:go_default_library",
        "//app/interface/live/app-room/conf:go_default_library",
        "//app/interface/live/app-room/model:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "banner.go",
        "gift.go",
        "roomNotice.go",
    ],
    importpath = "go-common/app/interface/live/app-room/service/v1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/app-room/api/http/v1:go_default_library",
        "//app/interface/live/app-room/conf:go_default_library",
        "//app/interface/live/app-room/dao:go_default_library",
        "//app/interface/live/app-room/model:go_default_library",
        "//app/service/live/gift/api/grpc/v1:go_default_library",
        "//app/service/live/resource/api/grpc/v1:go_default_library",
        "//app/service/live/xuserex/api/grpc/v1:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/sync/errgroup:go_default_library",
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
        "//app/interface/live/app-room/service/v1/dm:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
