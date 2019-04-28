package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["feed_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/web-feed/conf:go_default_library",
        "//app/interface/main/web-feed/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "feed.go",
        "service.go",
    ],
    importpath = "go-common/app/interface/main/web-feed/service",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/web-feed/conf:go_default_library",
        "//app/interface/main/web-feed/dao:go_default_library",
        "//app/interface/main/web-feed/model:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/account/rpc/client:go_default_library",
        "//app/service/main/feed/model:go_default_library",
        "//app/service/main/feed/rpc/client:go_default_library",
        "//library/cache:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
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
