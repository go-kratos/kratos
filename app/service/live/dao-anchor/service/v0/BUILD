package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["createData_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/live/dao-anchor/api/grpc/v0:go_default_library",
        "//app/service/live/dao-anchor/conf:go_default_library",
        "//app/service/live/dao-anchor/dao:go_default_library",
        "//library/net/metadata:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["createData.go"],
    importpath = "go-common/app/service/live/dao-anchor/service/v0",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/dao-anchor/api/grpc/v0:go_default_library",
        "//app/service/live/dao-anchor/api/grpc/v1:go_default_library",
        "//app/service/live/dao-anchor/conf:go_default_library",
        "//app/service/live/dao-anchor/dao:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/log:go_default_library",
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
