package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "audit_log.go",
        "mcn.go",
        "msg.go",
        "recommend.go",
        "service.go",
        "statistics.go",
        "upload.go",
    ],
    importpath = "go-common/app/admin/main/mcn/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/mcn/conf:go_default_library",
        "//app/admin/main/mcn/dao/bfs:go_default_library",
        "//app/admin/main/mcn/dao/msg:go_default_library",
        "//app/admin/main/mcn/dao/up:go_default_library",
        "//app/admin/main/mcn/dao/videoup:go_default_library",
        "//app/admin/main/mcn/model:go_default_library",
        "//app/interface/main/mcn/model/datamodel:go_default_library",
        "//app/interface/main/mcn/tool/worker:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/member/api:go_default_library",
        "//app/service/main/member/model/block:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/time:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "audit_log_test.go",
        "mcn_test.go",
        "msg_test.go",
        "recommend_test.go",
        "service_test.go",
        "statistics_test.go",
        "upload_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/mcn/conf:go_default_library",
        "//app/admin/main/mcn/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
