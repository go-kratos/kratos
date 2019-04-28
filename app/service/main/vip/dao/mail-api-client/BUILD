package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["mail_api_client.go"],
    importpath = "go-common/app/service/main/vip/dao/mail-api-client",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/net/http/blademaster:go_default_library",
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
    srcs = ["mail_api_client_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//library/net/http/blademaster:go_default_library",
        "//library/net/netutil/breaker:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
