package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["service_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/interface/live/live-demo/api/http/v2:go_default_library",
        "//app/interface/live/live-demo/conf:go_default_library",
        "//app/interface/live/live-demo/dao:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "foo.go",
        "foo2.go",
    ],
    importpath = "go-common/app/interface/live/live-demo/service/v2",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/live-demo/api/http/v2:go_default_library",
        "//app/interface/live/live-demo/conf:go_default_library",
        "//app/interface/live/live-demo/dao:go_default_library",
        "//app/service/live/room/api/liverpc/v2:go_default_library",
        "//library/ecode:go_default_library",
        "//library/ecode/pb:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
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
