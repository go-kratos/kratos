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
        "discovery.go",
        "node_proc.go",
        "sql.go",
        "tidb.go",
    ],
    importpath = "go-common/library/database/tidb",
    tags = ["automanaged"],
    deps = [
        "//library/conf/env:go_default_library",
        "//library/log:go_default_library",
        "//library/naming:go_default_library",
        "//library/naming/discovery:go_default_library",
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
    srcs = [
        "sql_test.go",
        "tidb_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/net/netutil/breaker:go_default_library",
        "//library/time:go_default_library",
    ],
)
