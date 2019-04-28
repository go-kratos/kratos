package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "bag_test.go",
        "databus_test.go",
        "discount_test.go",
        "gift_test.go",
        "infoc_test.go",
        "list_test.go",
        "package_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/live/gift/api/grpc/v1:go_default_library",
        "//app/service/live/gift/conf:go_default_library",
        "//app/service/live/gift/dao:go_default_library",
        "//app/service/live/gift/model:go_default_library",
        "//app/service/live/resource/sdk:go_default_library",
        "//vendor/github.com/satori/go.uuid:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "bag.go",
        "databus.go",
        "discount.go",
        "gift.go",
        "infoc.go",
        "list.go",
        "package.go",
    ],
    importpath = "go-common/app/service/live/gift/service/v1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/activity/api/liverpc/v1:go_default_library",
        "//app/service/live/fans_medal/api/liverpc/v2:go_default_library",
        "//app/service/live/gift/api/grpc/v1:go_default_library",
        "//app/service/live/gift/conf:go_default_library",
        "//app/service/live/gift/dao:go_default_library",
        "//app/service/live/gift/model:go_default_library",
        "//app/service/live/live_user/api/liverpc/v1:go_default_library",
        "//app/service/live/resource/sdk:go_default_library",
        "//app/service/live/room/api/liverpc/v1:go_default_library",
        "//app/service/live/user/api/liverpc/v3:go_default_library",
        "//app/service/live/xuser/api/grpc/v1:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/rpc/liverpc/context:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//vendor/github.com/satori/go.uuid:go_default_library",
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
