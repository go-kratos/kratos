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
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/member/conf:go_default_library",
        "//app/admin/main/member/model/block:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "audit_log.go",
        "block.go",
        "msg.go",
        "notify.go",
        "rpc.go",
        "service.go",
    ],
    importpath = "go-common/app/admin/main/member/service/block",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/member/conf:go_default_library",
        "//app/admin/main/member/dao/block:go_default_library",
        "//app/admin/main/member/model/block:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/figure/model:go_default_library",
        "//app/service/main/figure/rpc/client:go_default_library",
        "//app/service/main/spy/model:go_default_library",
        "//app/service/main/spy/rpc/client:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
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
