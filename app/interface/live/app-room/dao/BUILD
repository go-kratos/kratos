package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "dao_test.go",
        "livewallet_test.go",
        "paycenter_test.go",
        "userext_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/interface/live/app-room/conf:go_default_library",
        "//app/interface/live/app-room/model:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "api.go",
        "dao.go",
        "livewallet.go",
        "paycenter.go",
        "userext.go",
    ],
    importpath = "go-common/app/interface/live/app-room/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/app-room/conf:go_default_library",
        "//app/interface/live/app-room/model:go_default_library",
        "//app/service/live/userext/api/liverpc:go_default_library",
        "//app/service/live/userext/api/liverpc/v0:go_default_library",
        "//app/service/live/xuserex/api/grpc/v1:go_default_library",
        "//library/cache:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/interface/live/app-room/dao/av:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
