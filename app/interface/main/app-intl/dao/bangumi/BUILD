package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["dao_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-intl/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "grpc.go",
    ],
    importpath = "go-common/app/interface/main/app-intl/dao/bangumi",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-intl/conf:go_default_library",
        "//app/interface/main/app-intl/model/bangumi:go_default_library",
        "//app/service/openplatform/pgc-season/api/grpc/season/v1:go_default_library",
        "//library/ecode:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/xstr:go_default_library",
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
