package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "delCache_test.go",
        "exchange_test.go",
        "getTid_test.go",
        "get_test.go",
        "pay_test.go",
        "recharge_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/live/wallet/conf:go_default_library",
        "//app/service/live/wallet/dao:go_default_library",
        "//app/service/live/wallet/model:go_default_library",
        "//app/service/live/wallet/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "http.go",
        "wallet.go",
    ],
    importpath = "go-common/app/service/live/wallet/http",
    tags = ["automanaged"],
    deps = [
        "//app/service/live/wallet/conf:go_default_library",
        "//app/service/live/wallet/model:go_default_library",
        "//app/service/live/wallet/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
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
