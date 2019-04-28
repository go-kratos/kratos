load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

package(default_visibility = ["//visibility:public"])

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
        "archive_test.go",
        "dao_test.go",
        "es_test.go",
        "log_test.go",
        "mng_test.go",
        "mng_v2_test.go",
        "query_extra_test.go",
        "query_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/search/conf:go_default_library",
        "//app/admin/main/search/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/olivere/elastic.v5:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "dao.go",
        "es.go",
        "log.go",
        "mng.go",
        "mng_v2.go",
        "query.go",
        "query_extra.go",
    ],
    importpath = "go-common/app/admin/main/search/dao",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/search/conf:go_default_library",
        "//app/admin/main/search/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/gopkg.in/olivere/elastic.v5:go_default_library",
    ],
)
