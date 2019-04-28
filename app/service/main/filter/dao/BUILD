package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "ai_test.go",
        "area_test.go",
        "dao_test.go",
        "filter_test.go",
        "hbase_test.go",
        "http_test.go",
        "key_test.go",
        "memcache_test.go",
        "rpc_test.go",
        "white_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/filter/conf:go_default_library",
        "//app/service/main/filter/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "ai.go",
        "area.go",
        "dao.go",
        "filter.go",
        "hbase.go",
        "http.go",
        "key.go",
        "memcache.go",
        "rpc.go",
        "white.go",
    ],
    importpath = "go-common/app/service/main/filter/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/antispam/model:go_default_library",
        "//app/service/main/antispam/rpc/client:go_default_library",
        "//app/service/main/filter/conf:go_default_library",
        "//app/service/main/filter/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/hbase.v2:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
