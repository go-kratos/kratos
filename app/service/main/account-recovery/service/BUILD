package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "appeal.go",
        "data_check.go",
        "send_email.go",
        "service.go",
    ],
    importpath = "go-common/app/service/main/account-recovery/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/account-recovery/conf:go_default_library",
        "//app/service/main/account-recovery/dao:go_default_library",
        "//app/service/main/account-recovery/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/microcosm-cc/bluemonday:go_default_library",
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
    srcs = [
        "appeal_test.go",
        "data_check_test.go",
        "send_email_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/account-recovery/conf:go_default_library",
        "//app/service/main/account-recovery/model:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
