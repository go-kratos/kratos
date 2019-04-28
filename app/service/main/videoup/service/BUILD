package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "archive_test.go",
        "audit_log_test.go",
        "flow_test.go",
        "msg_test.go",
        "pgc_test.go",
        "porder_test.go",
        "service_test.go",
        "staff_test.go",
        "ugc_test.go",
        "util_test.go",
        "video_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/main/videoup/conf:go_default_library",
        "//app/service/main/videoup/model/archive:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//vendor/github.com/davecgh/go-spew/spew:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "audit_log.go",
        "databus.go",
        "dede.go",
        "extra_func.go",
        "flow.go",
        "monitor.go",
        "msg.go",
        "netsafe.go",
        "pgc.go",
        "porder.go",
        "report.go",
        "service.go",
        "staff.go",
        "tx.go",
        "ugc.go",
        "util.go",
        "video.go",
    ],
    importpath = "go-common/app/service/main/videoup/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/mcn/tool/worker:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/videoup/conf:go_default_library",
        "//app/service/main/videoup/dao/agent:go_default_library",
        "//app/service/main/videoup/dao/archive:go_default_library",
        "//app/service/main/videoup/dao/bgm:go_default_library",
        "//app/service/main/videoup/dao/databus:go_default_library",
        "//app/service/main/videoup/dao/dede:go_default_library",
        "//app/service/main/videoup/dao/manager:go_default_library",
        "//app/service/main/videoup/dao/monitor:go_default_library",
        "//app/service/main/videoup/dao/msg:go_default_library",
        "//app/service/main/videoup/dao/relation:go_default_library",
        "//app/service/main/videoup/dao/ups:go_default_library",
        "//app/service/main/videoup/model/archive:go_default_library",
        "//app/service/main/videoup/model/dede:go_default_library",
        "//app/service/main/videoup/model/message:go_default_library",
        "//app/service/main/videoup/model/prom:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/stat/prom:go_default_library",
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
