package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "allowance_test.go",
        "code_test.go",
        "coupon_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/coupon/api:go_default_library",
        "//app/service/main/coupon/conf:go_default_library",
        "//app/service/main/coupon/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "allowance.go",
        "branch.go",
        "code.go",
        "coupon.go",
        "coupon_cartoon.go",
        "service.go",
    ],
    importpath = "go-common/app/service/main/coupon/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/coupon/api:go_default_library",
        "//app/service/main/coupon/conf:go_default_library",
        "//app/service/main/coupon/dao:go_default_library",
        "//app/service/main/coupon/model:go_default_library",
        "//app/service/main/vipinfo/api:go_default_library",
        "//app/service/main/vipinfo/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/time:go_default_library",
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
