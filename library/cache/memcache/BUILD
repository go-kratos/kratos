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
    srcs = [
        ":package-srcs",
        "//library/cache/memcache/test:all-srcs",
    ],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
        "conn.go",
        "errors.go",
        "memcache.go",
        "mock.go",
        "pool.go",
        "trace.go",
        "util.go",
    ],
    importpath = "go-common/library/cache/memcache",
    tags = ["automanaged"],
    deps = [
        "//library/container/pool:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/stat:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "client_test.go",
        "conn_test.go",
        "memcache_test.go",
        "pool_test.go",
        "util_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/cache/memcache/test:go_default_library",
        "//library/container/pool:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
    ],
)
