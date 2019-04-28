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
        "commandinfo.go",
        "conn.go",
        "doc.go",
        "errors.go",
        "log.go",
        "mock.go",
        "pool.go",
        "pubsub.go",
        "redis.go",
        "reply.go",
        "scan.go",
        "script.go",
        "trace.go",
    ],
    importpath = "go-common/library/cache/redis",
    tags = ["automanaged"],
    deps = [
        "//library/container/pool:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/stat:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "commandinfo_test.go",
        "conn_test.go",
        "pool_test.go",
        "pubsub_test.go",
        "redis_test.go",
        "reply_test.go",
        "scan_test.go",
        "script_test.go",
        "test_test.go",
        "trace_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/container/pool:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
    ],
)
