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
        "dao_test.go",
        "delay_test.go",
        "desc_format_test.go",
        "flow_test.go",
        "history_test.go",
        "netsafe_test.go",
        "new_video_test.go",
        "porder_test.go",
        "redis_test.go",
        "report_test.go",
        "staff_test.go",
        "type_test.go",
        "video_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/main/videoup/conf:go_default_library",
        "//app/service/main/videoup/model/archive:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
        "//vendor/github.com/davecgh/go-spew/spew:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "dao.cache.go",
        "dao.go",
        "delay.go",
        "desc_format.go",
        "flow.go",
        "forbid.go",
        "history.go",
        "mc.cache.go",
        "netsafe.go",
        "new_video.go",
        "porder.go",
        "redis.go",
        "relation.go",
        "report.go",
        "staff.go",
        "staff_apply.go",
        "type.go",
        "video.go",
    ],
    importpath = "go-common/app/service/main/videoup/dao/archive",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/videoup/conf:go_default_library",
        "//app/service/main/videoup/model/archive:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/dgryski/go-farm:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/golang.org/x/sync/singleflight:go_default_library",
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
