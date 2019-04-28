package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "roomAdmin_test.go",
        "room_test.go",
        "topic_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/interface/live/app-ucenter/api/http/v1:go_default_library",
        "//app/interface/live/app-ucenter/conf:go_default_library",
        "//library/net/metadata:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "room.go",
        "roomAdmin.go",
        "topic.go",
    ],
    importpath = "go-common/app/interface/live/app-ucenter/service/v1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/app-ucenter/api/http/v1:go_default_library",
        "//app/interface/live/app-ucenter/conf:go_default_library",
        "//app/interface/live/app-ucenter/dao:go_default_library",
        "//app/service/live/xuser/api/grpc/v1:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
