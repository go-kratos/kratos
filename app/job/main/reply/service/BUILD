package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "index_test.go",
        "reply_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/reply/conf:go_default_library",
        "//vendor/github.com/go-sql-driver/mysql:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "admin.go",
        "business.go",
        "fold.go",
        "index.go",
        "notify.go",
        "reply.go",
        "service.go",
        "spam.go",
        "upper.go",
    ],
    importpath = "go-common/app/job/main/reply/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/interface/openplatform/article/rpc/client:go_default_library",
        "//app/job/main/reply/conf:go_default_library",
        "//app/job/main/reply/dao/message:go_default_library",
        "//app/job/main/reply/dao/notice:go_default_library",
        "//app/job/main/reply/dao/reply:go_default_library",
        "//app/job/main/reply/dao/search:go_default_library",
        "//app/job/main/reply/dao/spam:go_default_library",
        "//app/job/main/reply/dao/stat:go_default_library",
        "//app/job/main/reply/model/reply:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/assist/model/assist:go_default_library",
        "//app/service/main/assist/rpc/client:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//app/service/openplatform/pgc-season/api/grpc/episode/v1:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/time:go_default_library",
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
