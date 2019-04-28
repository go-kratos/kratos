load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
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
)

go_library(
    name = "go_default_library",
    srcs = [
        "mysql.go",
        "sql.go",
    ],
    importpath = "go-common/library/database/sql",
    tags = ["automanaged"],
    deps = [
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/netutil/breaker:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/stat:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/go-sql-driver/mysql:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["mysql_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/net/netutil/breaker:go_default_library",
        "//library/time:go_default_library",
    ],
)
