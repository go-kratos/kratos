package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "hbase.go",
        "memcache.go",
        "mysql.go",
        "rpc.go",
    ],
    importpath = "go-common/app/interface/main/creative/dao/weeklyhonor",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/up/util/hbaseutil:go_default_library",
        "//app/interface/main/creative/conf:go_default_library",
        "//app/interface/main/creative/model/weeklyhonor:go_default_library",
        "//app/service/main/up/api/v1:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/hbase.v2:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/tsuna/gohbase/hrpc:go_default_library",
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
        "dao_test.go",
        "hbase_test.go",
        "memcache_test.go",
        "mysql_test.go",
        "rpc_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/creative/conf:go_default_library",
        "//app/interface/main/creative/model/weeklyhonor:go_default_library",
        "//app/service/main/up/api/v1:go_default_library",
        "//library/database/hbase.v2:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
        "//vendor/github.com/golang/mock/gomock:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/github.com/tsuna/gohbase/hrpc:go_default_library",
    ],
)
