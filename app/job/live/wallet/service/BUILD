package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "service_test.go",
        "wallet_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/live/wallet/conf:go_default_library",
        "//app/job/live/wallet/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "service.go",
        "wallet.go",
    ],
    importpath = "go-common/app/job/live/wallet/service",
    tags = ["automanaged"],
    deps = [
        "//app/job/live/wallet/conf:go_default_library",
        "//app/job/live/wallet/dao:go_default_library",
        "//app/job/live/wallet/model:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
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
