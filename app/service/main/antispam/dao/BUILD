package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "cond_test.go",
        "dao_test.go",
        "keyword_test.go",
        "mysql_test.go",
        "redis_test.go",
        "regexp_test.go",
        "rule_test.go",
        "transaction_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/antispam/conf:go_default_library",
        "//app/service/main/antispam/model:go_default_library",
        "//app/service/main/antispam/util:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "cond.go",
        "dao.go",
        "keyword.go",
        "mysql.go",
        "redis.go",
        "regexp.go",
        "rule.go",
        "transaction.go",
    ],
    importpath = "go-common/app/service/main/antispam/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/antispam/conf:go_default_library",
        "//app/service/main/antispam/model:go_default_library",
        "//app/service/main/antispam/util:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
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
